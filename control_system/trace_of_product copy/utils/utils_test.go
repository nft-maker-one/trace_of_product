package utils

import (
	"agricultural_meta/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	assert.Nil(t, LogMsg([]string{"name"}, []string{"张三"}))
	assert.NotNil(t, LogMsg([]string{"name"}, []string{"alice", "bob"}))
	assert.Nil(t, LogWarn([]string{"name"}, []string{"张三"}))
	assert.NotNil(t, LogWarn([]string{"name"}, []string{"alice", "bob"}))
	assert.Nil(t, LogDebug([]string{"name"}, []string{"张三"}))
	assert.NotNil(t, LogDebug([]string{"name"}, []string{"alice", "bob"}))
	assert.Nil(t, LogError([]string{"name"}, []string{"张三"}))
	assert.NotNil(t, LogError([]string{"name"}, []string{"alice", "bob"}))

}

func TestYaml(t *testing.T) {
	config, err := NewConfig("../config.yaml")
	assert.Nil(t, err)
	LogMsg([]string{"dsn"}, []string{config.Mysql.Dsn})
}

func TestCompareBytes(t *testing.T) {
	a := types.RandomBytes(20)
	b := make([]byte, 20)
	copy(b, a)
	assert.True(t, CompareBytes(a, b))
	b[1] = 13
	assert.False(t, CompareBytes(a, b))
}
