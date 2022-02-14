package riotinterface

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPuuidFromSummonerName(t *testing.T) {
	assert := assert.New(t)
	puuid, err := GetPuuidFromSummonerName("AngleAbuser")
	assert.Equal("gyFa4zzThNzWkwO2VMRHuBqqHs1VwSJg_evd7QKSvwsWj0ipsax0u2dmn4b-mRp66hnjEUYZu0gm2Q", puuid)
	assert.Equal(nil, err)
}
