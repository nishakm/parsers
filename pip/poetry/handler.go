// SPDX-License-Identifier: Apache-2.0

package poetry

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/opensbom-generator/parsers/internal/helper"
	"github.com/opensbom-generator/parsers/meta"
	"github.com/opensbom-generator/parsers/pip/worker"
	"github.com/opensbom-generator/parsers/plugin"
)

const cmdName = "poetry"
const manifestFile = "pyproject.toml"
const manifestLockFile = "poetry.lock"
const placeholderPkgName = "{PACKAGE}"

var errDependenciesNotFound = errors.New("Unable to generate SPDX file: no modules or vendors found. Please install them before running spdx-sbom-generator, e.g.: `poetry install` or `poetry update`")
var errBuildlingModuleDependencies = errors.New("Error building module dependencies")
var errNoPipCommand = errors.New("Cannot find the poetry command")
var errVersionNotFound = errors.New("Python version not found")
var errFailedToConvertModules = errors.New("Failed to convert modules")

type poetry struct {
	metadata   plugin.Metadata
	rootModule *meta.Package
	command    *helper.Cmd
	basepath   string
	version    string
	pkgs       []worker.Packages
	metainfo   map[string]worker.Metadata
	allModules []meta.Package
}

// New ...
func New() *poetry {
	return &poetry{
		metadata: plugin.Metadata{
			Name:       "The Python Package Index (PyPI)",
			Slug:       "poetry",
			Manifest:   []string{manifestLockFile},
			ModulePath: []string{},
		},
	}
}

// Get Metadata ...
func (m *poetry) GetMetadata() plugin.Metadata {
	return m.metadata
}

// Is Valid ...
func (m *poetry) IsValid(path string) bool {
	for i := range m.metadata.Manifest {
		if helper.Exists(filepath.Join(path, m.metadata.Manifest[i])) {
			return true
		}
	}
	return false
}

// Has Modules Installed ...
func (m *poetry) HasModulesInstalled(path string) error {
	if err := m.buildCmd(ModulesCmd, m.basepath); err != nil {
		return err
	}
	result, err := m.command.Output()
	if err == nil && len(result) > 0 && worker.IsRequirementMeet(result) {
		return nil
	}
	return errDependenciesNotFound
}

// Get Version ...
func (m *poetry) GetVersion() (string, error) {
	if err := m.buildCmd(VersionCmd, m.basepath); err != nil {
		return "", err
	}
	version, err := m.command.Output()
	m.version = worker.GetShortPythonVersion(version)
	if err != nil {
		return "Python", errVersionNotFound
	}
	return version, err
}

// Set Root Module ...
func (m *poetry) SetRootModule(path string) error {
	m.basepath = path
	return nil
}

// Get Root Module ...
func (m *poetry) GetRootModule(path string) (*meta.Package, error) {
	if m.rootModule == nil {
		module := m.fetchRootModule()
		m.rootModule = &module
	}
	return m.rootModule, nil
}

// List Used Modules...
func (m *poetry) ListUsedModules(path string) ([]meta.Package, error) {
	if err := m.LoadModuleList(path); err != nil {
		return m.allModules, errFailedToConvertModules
	}

	decoder := worker.NewMetadataDecoder(m.GetPackageDetails)
	metainfo, err := decoder.ConvertMetadataToModules(m.pkgs, &m.allModules)
	if err != nil {
		return m.allModules, err
	}
	m.metainfo = metainfo

	return m.allModules, nil
}

// List Modules With Deps ...
func (m *poetry) ListModulesWithDeps(path string, globalSettingFile string) ([]meta.Package, error) {
	modules, err := m.ListUsedModules(path)
	if err != nil {
		return nil, err
	}
	m.GetRootModule(path)
	if err := worker.BuildDependencyGraph(&m.allModules, &m.metainfo); err != nil {
		return nil, err
	}
	return modules, err
}

func (m *poetry) buildCmd(cmd command, path string) error {
	cmdArgs := cmd.Parse()
	if cmdArgs[0] != cmdName {
		return errNoPipCommand
	}

	command := helper.NewCmd(helper.CmdOptions{
		Name:      cmdArgs[0],
		Args:      cmdArgs[1:],
		Directory: path,
	})

	m.command = command

	return command.Build()
}

func (m *poetry) GetPackageDetails(packageName string) (string, error) {
	metatdataCmd := command(strings.ReplaceAll(string(MetadataCmd), placeholderPkgName, packageName))

	m.buildCmd(metatdataCmd, m.basepath)
	result, err := m.command.Output()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (m *poetry) PushRootModuleToVenv() (bool, error) {
	if err := m.buildCmd(InstallRootModuleCmd, m.basepath); err != nil {
		return false, err
	}
	result, err := m.command.Output()
	if err == nil && len(result) > 0 {
		return true, err
	}
	return false, nil
}

func (m *poetry) markRootModue() {
	for i, pkg := range m.pkgs {
		if worker.IsRootModule(pkg, m.metadata.Slug) {
			m.pkgs[i].Root = true
			break
		}
	}
}

func (m *poetry) LoadModuleList(path string) error {
	state, err := m.PushRootModuleToVenv()
	if err != nil && !state {
		return err
	}
	m.buildCmd(ModulesCmd, m.basepath)
	result, err := m.command.Output()
	if err == nil && len(result) > 0 && worker.IsRequirementMeet(result) {
		m.pkgs = worker.LoadModules(result, m.version)
		m.markRootModue()
	}
	return err
}

func (m *poetry) fetchRootModule() meta.Package {
	for _, mod := range m.allModules {
		if mod.Root {
			return mod
		}
	}
	return meta.Package{}
}
