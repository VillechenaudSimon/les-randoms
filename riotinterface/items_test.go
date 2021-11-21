package riotinterface

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetItemsMap(t *testing.T) {
	assert := assert.New(t)
	items := GetItemsMap()
	assert.Equal("Boots", items["1001"].Name)
}
