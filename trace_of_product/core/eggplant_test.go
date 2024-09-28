package core

import (
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func randomMetaData() MetaData {
	return MetaData{
		EggplantId:      rand.Intn(math.MaxInt),
		ProductHeight:   rand.Intn(math.MaxInt),
		ProductHash:     types.RandomHash(),
		TransportHeight: rand.Intn(math.MaxInt),
		TransportHash:   types.RandomHash(),
		ProcessHeight:   rand.Intn(math.MaxInt),
		ProcessHash:     types.RandomHash(),
		StorageHeight:   rand.Intn(math.MaxInt),
		StorageHash:     types.RandomHash(),
		SellHeight:      rand.Intn(math.MaxInt),
		SellHash:        types.RandomHash(),
	}
}

func randomEggpalnt() *Eggplant {
	egg := &Eggplant{}
	egg.MetaData = randomMetaData()
	egg.FirstSeen = time.Now().Unix()
	return egg

}

func randomEggplantWithSignature() *Eggplant {
	egg := randomEggpalnt()
	priKey := crypto.GeneratePrivateKey()
	if err := egg.Sign(priKey); err != nil {
		logrus.Error(err)
		return nil
	}
	return egg
}

func TestEggplantSignAndVerify(t *testing.T) {
	test_data := ""
	for i := 0; i < 20; i++ {
		egg := randomEggpalnt()
		test_data += fmt.Sprintf("农产品 %d\n%+v\n", i+1, egg)
		priKey := crypto.GeneratePrivateKey()
		assert.Nil(t, egg.Sign(priKey))
		assert.Nil(t, egg.Verify())
	}
	fmt.Println(test_data)

}

func TestEggplantHash(t *testing.T) {
	test_data := ""
	for i := 0; i < 20; i++ {
		egg := randomEggpalnt()
		test_data += fmt.Sprintf("农产品 %d\n%+v\n", i+1, egg)
		egg.SetHash(EggplantHasher{})
		assert.False(t, egg.Hash.IsZero())
		fmt.Printf("Eggplant %d hash %v\n", i+1, egg.Hash)
	}
	fmt.Println(test_data)

}

func TestEggplantEncodeDecode(t *testing.T) {
	test_data := ""
	test_res := ""
	for i := 0; i < 20; i++ {
		egg := randomEggpalnt()
		egg.SetHash(EggplantHasher{})
		priKey := crypto.GeneratePrivateKey()
		assert.Nil(t, egg.Sign(priKey))
		buf := &bytes.Buffer{}
		test_data += fmt.Sprintf("农产品 %d (编码前)\n%+v\n", i+1, egg)
		test_res += fmt.Sprintf("农产品 %d (编码前)\n%+v\n", i+1, egg)
		assert.Nil(t, NewGobEggplantEncoder(buf).Encode(egg))
		test_res += fmt.Sprintf("农产品 %d (编码结果)\n%+v\n", i+1, buf.Bytes())
		newEgg := &Eggplant{}
		assert.Nil(t, NewGobEggplantDecoder(buf).Decode(newEgg))
		test_res += fmt.Sprintf("农产品 %d (解码结果)\n%+v\n", i+1, egg)
		assert.Equal(t, *egg, *newEgg)
	}
	fmt.Println(test_data)
	fmt.Println(test_res)

}
