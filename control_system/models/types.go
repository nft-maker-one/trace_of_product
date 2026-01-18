package models

import (
	"encoding/hex"
	"fmt"
	"time"
)

type User struct {
	ID                 int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName           string    `gorm:"column:user_name;size:255" json:"user_name"`
	NickName           string    `gorm:"column:nick_name;size:255" json:"nick_name"`
	Email              string    `gorm:"column:email;size:255" json:"email"`
	Password           string    `gorm:"column:password;size:255" json:"password"`
	InviteCode         string    `gorm:"column:invite_code;size:100" json:"invite_code"`
	AgreeTerms         bool      `gorm:"column:agree_terms;default:false" json:"agree_terms"`
	AvatarURL          string    `gorm:"column:avatar_url;size:255;default:'/public/avatars/profile.jpg'" json:"avatar_url"`
	ProfilePublic      bool      `gorm:"column:profile_public;default:true" json:"profile_public"`
	EmailNotifications bool      `gorm:"column:email_notifications;default:true" json:"email_notifications"`
	LastLoginAt        *int64    `gorm:"column:last_login_at" json:"last_login_at"`
	CreatedAt          time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

type ConsortiumNode struct {
	Id         int64  `gorm:"primaryKey;column:id" json:"id"`
	Addr       string `gorm:"column:addr;size:40" json:"addr"`
	PubKey     []byte `gorm:"column:pub_key" json:"pub_key"`
	CreateTime int64  `gorm:"column:create_time" json:"create_time"`
	VerifyTime int    `gorm:"column:verify_time" json:"verify_time"`
}

// TableName 指定表名
func (ConsortiumNode) TableName() string {
	return "consortium_nodes"
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
	res.ProcessHeight = r.ProcessHeight
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

func (m MetaData) Verify() bool {
	return true
}
