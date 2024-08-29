package core

import (
	"agricultural_meta/types"
	"encoding/gob"
	"io"
	"math/big"
)

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
	return gob.NewEncoder(g.w).Encode(egg)
}

type GobEggplantDecoder struct {
	r io.Reader
}

func NewGobEggplantDecoder(r io.Reader) *GobEggplantDecoder {
	return &GobEggplantDecoder{r: r}
}

func (g *GobEggplantDecoder) Decode(egg *Eggplant) error {

	if err := gob.NewDecoder(g.r).Decode(egg); err != nil {
		return err
	}
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
	return gob.NewEncoder(enc.w).Encode(b)
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
	if err := gob.NewDecoder(dec.r).Decode(b); err != nil {
		return err
	}
	return nil
}

func init() {
	gob.Register(types.Hash{})
	gob.Register(&big.Int{})
	gob.Register(MetaData{})
	gob.Register(Header{})
}
