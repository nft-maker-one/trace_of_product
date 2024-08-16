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
	key, err := crypto.ReadPubKey(1)
	assert.Nil(t, err)
	assert.Nil(t, db.AddNode(1, "127.0.0.1:8081", *key))
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
