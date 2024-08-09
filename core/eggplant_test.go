package core

import (
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"bytes"
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
	egg.SetFirsstSeen(time.Now().Unix())
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
	egg := randomEggpalnt()
	priKey := crypto.GeneratePrivateKey()
	assert.Nil(t, egg.Sign(priKey))
	assert.Nil(t, egg.Verify())
}

func TestEggplantHash(t *testing.T) {
	egg := randomEggpalnt()
	egg.SetHash(egg.Hash(EggplantHasher{}))
	assert.False(t, egg.hash.IsZero())
	logrus.Info("eggHash:", egg.hash)
}

func TestEggplantEncodeDecode(t *testing.T) {
	egg := randomEggpalnt()
	egg.SetHash(egg.Hash(EggplantHasher{}))
	priKey := crypto.GeneratePrivateKey()
	assert.Nil(t, egg.Sign(priKey))
	buf := &bytes.Buffer{}
	assert.Nil(t, NewGobEggplantEncoder(buf).Encode(egg))
	newEgg := &Eggplant{}
	assert.Nil(t, NewGobEggplantDecoder(buf).Decode(newEgg))
	assert.Equal(t, *egg, *newEgg)

}
