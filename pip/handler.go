// SPDX-License-Identifier: Apache-2.0

package pip

import (
	"github.com/opensbom-generator/parsers/meta"
	"github.com/opensbom-generator/parsers/pip/pipenv"
	"github.com/opensbom-generator/parsers/pip/poetry"
	"github.com/opensbom-generator/parsers/pip/pyenv"

	"github.com/opensbom-generator/parsers/plugin"
)

type pip struct {
	plugin plugin.Plugin
}

// New ...
func New() *pip {
	return &pip{
		plugin: nil,
	}
}

// Get Metadata ...
func (m *pip) GetMetadata() plugin.Metadata {
	return m.plugin.GetMetadata()
}

// Is Valid ...
func (m *pip) IsValid(path string) bool {
	if p := pipenv.New(); p.IsValid(path) {
		m.plugin = p
		return true
	}

	if p := poetry.New(); p.IsValid(path) {
		m.plugin = p
		return true
	}

	if p := pyenv.New(); p.IsValid(path) {
		m.plugin = p
		return true
	}

	return false
}

// Has Modules Installed ...
func (m *pip) HasModulesInstalled(path string) error {
	return m.plugin.HasModulesInstalled(path)
}

// Get Version ...
func (m *pip) GetVersion() (string, error) {
	return m.plugin.GetVersion()
}

// Set Root Module ...
func (m *pip) SetRootModule(path string) error {
	return m.plugin.SetRootModule(path)
}

// Get Root Module ...
func (m *pip) GetRootModule(path string) (*meta.Package, error) {
	return m.plugin.GetRootModule(path)
}

// List Used Modules...
func (m *pip) ListUsedModules(path string) ([]meta.Package, error) {
	return m.plugin.ListUsedModules(path)
}

// List Modules With Deps ...
func (m *pip) ListModulesWithDeps(path string, globalSettingFile string) ([]meta.Package, error) {
	return m.plugin.ListModulesWithDeps(path, globalSettingFile)
}
