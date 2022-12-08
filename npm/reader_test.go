package npm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseManifestV2(t *testing.T) {
	// This file should exist
	data, err := ReadManifest("test/package-lock-v2.json")
	assert.Nil(t, err)
	assert.NotNil(t, data)
	// The file should be parsable
	lock, err := ParseManifestV2(data)
	assert.Nil(t, err)
	assert.NotNil(t, lock)
	// test lockfile content
	assert.Equal(t, lock.Name, "e-commerce")
	assert.Equal(t, lock.LockfileVersion, 2)
	// test root package content
	assert.NotNil(t, lock.RootPackage)
	assert.Equal(t, lock.RootPackage.Name, "e-commerce")
	assert.Equal(t, lock.RootPackage.License, "ISC")
	assert.Equal(t, lock.RootPackage.Dependencies["bcryptjs"], "^2.4.3")
	assert.Equal(t, lock.RootPackage.DevDependencies["babel-preset-env"], "^1.7.0")
	// test packages
	assert.NotNil(t, lock.Packages)
	assert.Nil(t, lock.Packages["node_modules/ansi-regex"].Dependencies)
	assert.NotNil(t, lock.Packages["node_modules/ansi-regex"].Engines)
	assert.Equal(t, lock.Packages["node_modules/ansi-regex"].Version, "2.1.1")
	assert.Nil(t, lock.Packages["node_modules/babel-code-frame"].Engines)
	assert.True(t, lock.Packages["node_modules/babel-code-frame"].Dev)
	assert.NotNil(t, lock.Packages["node_modules/babel-code-frame"].Dependencies)
	assert.Equal(t, lock.Packages["node_modules/babel-code-frame"].Dependencies["chalk"], "^1.1.3")
	assert.Equal(t, lock.Packages["node_modules/babylon"].Bin["babylon"], "bin/babylon.js")
	assert.False(t, lock.Packages["node_modules/call-bind"].Dev)
	assert.True(t, lock.Packages["node_modules/core-js"].HasInstallScript)

}
