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
