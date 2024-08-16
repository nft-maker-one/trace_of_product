package crypto

import (
	"agricultural_meta/utils"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func GenerateKeyPair(id int) PublicKey {
	if _, err := os.Stat("../Keys/" + strconv.Itoa(id)); err == nil {
		utils.LogMsg([]string{"Mkdir"}, []string{"key for  " + strconv.Itoa(id) + " already exist"})
	}
	fmt.Println("Creating directory .....")
	if err := os.Mkdir("../Keys/"+strconv.Itoa(id), 0644); err != nil {
		utils.LogMsg([]string{"Mkdir"}, []string{err.Error()})
	}
	utils.LogMsg([]string{"Mkdir"}, []string{"key for  " + strconv.Itoa(id) + " created succesfully"})
	priKey := GeneratePrivateKey()
	priKeyBytes, err := x509.MarshalECPrivateKey(priKey.key)
	if err != nil {
		logrus.Errorln(strconv.Itoa(id) + " node create private key failed" + err.Error())
	}
	block := &pem.Block{
		Type:  "ecdsa private key",
		Bytes: priKeyBytes,
	}
	blockBytes := pem.EncodeToMemory(block)
	prifile, err := os.OpenFile("../Keys/"+strconv.Itoa(id)+"/"+"PRI_KEY", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer prifile.Close()
	if _, err := prifile.Write(blockBytes); err != nil {
		panic(err)
	}
	utils.LogMsg([]string{"GenerateKey"}, []string{"create private key successfully for " + strconv.Itoa(id)})
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(priKey.PublicKey().Key)
	if err != nil {
		logrus.Errorln(strconv.Itoa(id) + " node create public key failed" + err.Error())
	}
	block = &pem.Block{
		Type:  "ecdsa public key",
		Bytes: pubKeyBytes,
	}
	blockBytes = pem.EncodeToMemory(block)
	pubfile, err := os.OpenFile("../Keys/"+strconv.Itoa(id)+"/"+"PUB_KEY", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer pubfile.Close()
	if _, err := pubfile.Write(blockBytes); err != nil {
		panic(err)
	}
	utils.LogMsg([]string{"GenerateKey"}, []string{"create public key successfully for " + strconv.Itoa(id)})
	return priKey.PublicKey()
}

func ReadPriKey(id int) (*PrivateKey, error) {
	data, err := os.ReadFile("../Keys/" + strconv.Itoa(id) + "/" + "PRI_KEY")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("decode block error")
	}
	priKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key := PrivateKey{}
	key.key = priKey
	utils.LogMsg([]string{"ReadPriKey"}, []string{fmt.Sprintf("read private key for %d successfully", id)})
	return &key, nil
}

func ReadPubKey(id int) (*PublicKey, error) {
	data, err := os.ReadFile("../Keys/" + strconv.Itoa(id) + "/" + "PUB_KEY")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("decode block error")
	}
	pubKeyRow, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubKey, ok := pubKeyRow.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not ecdsa pubKey")
	}
	key := PublicKey{}
	key.Key = pubKey
	utils.LogMsg([]string{"ReadPubKey"}, []string{fmt.Sprintf("read public key for %d successfully", id)})
	return &key, nil
}
