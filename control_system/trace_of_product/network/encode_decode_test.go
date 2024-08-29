package network

import (
	"agricultural_meta/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitData(t *testing.T) {
	hash := types.RandomHash()
	nonce := uint64(randomId())
	cData := CommitData(hash, nonce)
	hash1, nonce1 := CommitDataSplit(cData)
	assert.Equal(t, hash, hash1)
	assert.Equal(t, nonce, nonce1)
}
