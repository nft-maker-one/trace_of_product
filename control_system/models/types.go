package models

import (
	"encoding/hex"
	"fmt"
)

type User struct {
	UserName string `gorm:"user_name" json:"user_name"`
	Password string `gorm:"password" json:"password"`
}

type ConsortiumNode struct {
	Id         int    `gorm:"id"`
	Addr       string `gorm:"addr"`
	PubKey     []byte `gorm:"pubkey"`
	CreateTime int64  `gorm:"time"`
}

type Hash [32]uint8

type ReqMetaData struct {
	EggplantId      int    `json:"eggplant_id"`
	ProductHeight   int    `json:"product_height"`
	ProductHash     string `json:"product_hash"`
	TransportHeight int    `json:"transport_height"`
	TransportHash   string `json:"transport_hash"`
	ProcessHeight   int    `json:"process_height"`
	ProcessHash     string `json:"process_hash"`
	StorageHeight   int    `json:"storage_height"`
	StorageHash     string `json:"storage_hash"`
	SellHeight      int    `json:"sell_height"`
	SellHash        string `json:"sell_hash"`
}
type MetaData struct {
	EggplantId      int  `json:"eggplant_id"`
	ProductHeight   int  `json:"product_height"`
	ProductHash     Hash `json:"product_hash"`
	TransportHeight int  `json:"transport_height"`
	TransportHash   Hash `json:"transport_hash"`
	ProcessHeight   int  `json:"process_height"`
	ProcessHash     Hash `json:"process_hash"`
	StorageHeight   int  `json:"storage_height"`
	StorageHash     Hash `json:"storage_hash"`
	SellHeight      int  `json:"sell_height"`
	SellHash        Hash `json:"sell_hash"`
}

type ClientRequest struct {
	Header   string `json:"header"`
	RespAddr string `json:"resp_addr"`
	Content  []byte `content:"content"`
}

type Eggplant struct {
	MetaData
	NodeId     int
	PublickKey []byte //the Validator of this Eggplant
	Signature  []byte //the Signature of the Validator
	Hash       Hash   //the digest for eggplant's metadata
	FirstSeen  int64  //the creation time of the eggplant
}

func (r ReqMetaData) ToMetaData() MetaData {
	pHash, err := hex.DecodeString(r.ProductHash)
	if err != nil {
		return MetaData{}
	}
	tHash, err := hex.DecodeString(r.TransportHash)
	if err != nil {
		return MetaData{}
	}
	pcHash, err := hex.DecodeString(r.ProcessHash)
	if err != nil {
		return MetaData{}
	}
	sHash, err := hex.DecodeString(r.StorageHash)
	if err != nil {
		return MetaData{}
	}
	seHash, err := hex.DecodeString(r.SellHash)
	if err != nil {
		return MetaData{}
	}
	res := MetaData{}
	res.EggplantId = r.EggplantId
	res.ProductHash, _ = BytesToHash(pHash)
	res.ProductHeight = r.ProductHeight
	res.TransportHash, _ = BytesToHash(tHash)
	res.TransportHeight = r.TransportHeight
	res.ProcessHash, _ = BytesToHash(pcHash)
	res.ProductHeight = r.ProcessHeight
	res.StorageHash, _ = BytesToHash(sHash)
	res.StorageHeight = r.StorageHeight
	res.SellHash, _ = BytesToHash(seHash)
	res.SellHeight = r.SellHeight
	return res

}

func BytesToHash(data []byte) (Hash, error) {
	if len(data) != 32 {
		return Hash{}, fmt.Errorf("expected length of 32,but got %v", len(data))
	}
	hash := Hash{}
	for i := 0; i < len(hash); i++ {
		hash[i] = data[i]
	}
	return hash, nil
}
