package database

import (
	"agricultural_meta/crypto"
	"agricultural_meta/utils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDb(t *testing.T) {
	config, err := utils.NewConfig("../config.yaml")
	assert.Nil(t, err)
	db := InitNodeDb(config.Mysql.Dsn)
	assert.NotNil(t, db)
}

func TestAddNode(t *testing.T) {
	config, err := utils.NewConfig("../config.yaml")
	assert.Nil(t, err)
	db := InitNodeDb(config.Mysql.Dsn)
	assert.NotNil(t, db)
	key := crypto.GenerateKeyPair(1)
	assert.Nil(t, err)

	leaderId, err := db.AddNode(1, "127.0.0.1:8081", key)
	fmt.Println(leaderId)
	assert.Nil(t, err)
	for i := 2; i <= 5; i++ {
		key := crypto.GenerateKeyPair(i)
		leaderId, err = db.AddNode(i, "127.0.0.1:8081", key)
		assert.Equal(t, leaderId, 1)
		assert.Nil(t, err)
	}

}

func TestDeleteNode(t *testing.T) {
	config, err := utils.NewConfig("../config.yaml")
	assert.Nil(t, err)
	db := InitNodeDb(config.Mysql.Dsn)
	assert.NotNil(t, db)
	assert.Nil(t, db.DeleteNode(1))
}

func TestSearchNode(t *testing.T) {
	config, err := utils.NewConfig("../config.yaml")
	assert.Nil(t, err)
	db := InitNodeDb(config.Mysql.Dsn)
	assert.NotNil(t, db)
	node, err := db.SearchNodeById(1)
	assert.Nil(t, err)
	fmt.Println(node)
}

func TestCompareByteInDb(t *testing.T) {
	config, err := utils.NewConfig("../config.yaml")
	assert.Nil(t, err)
	db := InitNodeDb(config.Mysql.Dsn)
	assert.NotNil(t, db)
	key := crypto.GenerateKeyPair(1)
	assert.Nil(t, err)

	_, err = db.AddNode(1, "127.0.0.1:8081", key)
	assert.Nil(t, err)
	node, err := db.SearchNodeById(1)
	assert.Nil(t, err)
	keyBytes := key.ToSlice()
	assert.True(t, utils.CompareBytes(keyBytes, node.PubKey))
}
