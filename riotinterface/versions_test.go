package riotinterface

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersionsArray(t *testing.T) {
	assert := assert.New(t)
	versions, err := GetVersionsArray()
	assert.Equal("lolpatch_3.7", versions[len(versions)-1])
	assert.Equal(nil, err)
}

func TestGetLastVersion(t *testing.T) {
	assert := assert.New(t)
	versions, _ := GetVersionsArray()

	lastVersion, err := GetLastVersion()
	assert.Equal(versions[0], lastVersion)
	assert.Equal(nil, err)
}

func TestGetLastVersionFromGameVersion(t *testing.T) {
	assert := assert.New(t)

	lastVersion, err := GetLastVersionFromGameVersion("11.7.366.7612")
	assert.Equal("11.7.1", lastVersion)
	assert.Equal(nil, err)

	lastVersion, err = GetLastVersionFromGameVersion("9.5.11111111.7612")
	assert.Equal("9.5.1", lastVersion)
	assert.Equal(nil, err)
}
