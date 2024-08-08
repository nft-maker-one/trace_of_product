package core

import (
	"agricultural_meta/crypto"
	"agricultural_meta/types"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"io"
	"math/big"
)

type BlockTemp struct {
	Header    *Header
	Eggplants [][]byte
	PubKey_x  *big.Int
	PubKey_y  *big.Int
	Sig_r     *big.Int
	Sig_s     *big.Int
	hash      types.Hash
}

type EggplantTemp struct {
	Data      MetaData
	PubKey_x  *big.Int
	PubKey_y  *big.Int
	Sig_r     *big.Int
	Sig_s     *big.Int
	Hash      types.Hash
	FirstSeen int64
}

type Encoder[T any] interface {
	Encode(T) error
}

type Decoder[T any] interface {
	Decode(T) error
}

type GobEggplantEncoder struct {
	w io.Writer
}

func NewGobEggplantEncoder(w io.Writer) *GobEggplantEncoder {
	return &GobEggplantEncoder{
		w: w,
	}
}

func (g *GobEggplantEncoder) Encode(egg *Eggplant) error {
	eggplant := EggplantTemp{}
	eggplant.Data = egg.MetaData
	eggplant.PubKey_x = egg.PublickKey.Key.X
	eggplant.PubKey_y = egg.PublickKey.Key.Y
	eggplant.FirstSeen = egg.firstSeen
	eggplant.Hash = egg.hash
	eggplant.Sig_r = egg.Signature.R
	eggplant.Sig_s = egg.Signature.S
	return gob.NewEncoder(g.w).Encode(eggplant)
}

type GobEggplantDecoder struct {
	r io.Reader
}

func NewGobEggplantDecoder(r io.Reader) *GobEggplantDecoder {
	return &GobEggplantDecoder{r: r}
}

func (g *GobEggplantDecoder) Decode(egg *Eggplant) error {
	tempEgg := EggplantTemp{}
	if err := gob.NewDecoder(g.r).Decode(tempEgg); err != nil {
		return err
	}
	if egg.PublickKey.Key == nil {
		egg.PublickKey.Key = new(ecdsa.PublicKey)
	}
	if egg.Signature == nil {
		egg.Signature = &crypto.Signature{}
	}
	egg.MetaData = tempEgg.Data
	egg.firstSeen = tempEgg.FirstSeen
	egg.PublickKey.Key.X = tempEgg.PubKey_x
	egg.PublickKey.Key.Y = tempEgg.PubKey_y
	egg.PublickKey.Key.Curve = elliptic.P256()
	egg.Signature.R = tempEgg.Sig_r
	egg.Signature.S = tempEgg.Sig_s
	return nil

}

type GobBlockEncoder struct {
	w io.Writer
}

func NewGobBlockEncoder(w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		w: w,
	}
}

func (enc *GobBlockEncoder) Encode(b *Block) error {
	blockTemp := BlockTemp{}
	blockTemp.Header = b.Header
	blockTemp.PubKey_x = b.Validator.Key.X
	blockTemp.PubKey_y = b.Validator.Key.Y
	blockTemp.Sig_r = b.Signature.R
	blockTemp.Sig_s = b.Signature.S
	blockTemp.hash = b.hash
	blockTemp.Eggplants = make([][]byte, 0)
	for _, egg := range b.Eggplants {
		buf := &bytes.Buffer{}
		NewGobEggplantEncoder(buf).Encode(egg)
		blockTemp.Eggplants = append(blockTemp.Eggplants, buf.Bytes())
	}

	err := gob.NewEncoder(enc.w).Encode(&blockTemp)
	return err
}

type GobBlockDecoder struct {
	r io.Reader
}

func NewGobBlockDecode(r io.Reader) *GobBlockDecoder {
	return &GobBlockDecoder{
		r: r,
	}
}

func (dec *GobBlockDecoder) Decode(b *Block) error {
	blockTemp := BlockTemp{}
	if err := gob.NewDecoder(dec.r).Decode(&blockTemp); err != nil {
		return err
	}
	b.hash = blockTemp.hash
	b.Header = blockTemp.Header
	for _, eggByte := range blockTemp.Eggplants {
		egg := &Eggplant{}
		err := NewGobEggplantDecoder(bytes.NewReader(eggByte)).Decode(egg)
		if err != nil {
			return err
		}
		b.Eggplants = append(b.Eggplants, egg)

	}
	if b.Validator.Key == nil {
		b.Validator.Key = new(ecdsa.PublicKey)
	}
	if b.Signature == nil {
		b.Signature = new(crypto.Signature)
	}
	b.Validator.Key.X = blockTemp.PubKey_x
	b.Validator.Key.Y = blockTemp.PubKey_y
	b.Validator.Key.Curve = elliptic.P256()
	b.Signature.R = blockTemp.Sig_r
	b.Signature.S = blockTemp.Sig_s
	return nil
}

func init() {
	gob.Register(types.Hash{})
	gob.Register(&big.Int{})
	gob.Register(MetaData{})
	gob.Register(Header{})
}
