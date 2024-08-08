package main

import (
	"agricultural_meta/core"
	"agricultural_meta/crypto"
	"agricultural_meta/network"
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
)

// Server
// Transport => tcp,udp
// Block
// Tx
// Keypair

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trReomteA := network.NewLocalTransport("REMOTE_A")
	trReomteB := network.NewLocalTransport("REMOTE_B")
	trReomteC := network.NewLocalTransport("REMOTE_C")
	// trReomteD := network.NewLocalTransport("REMOTE_D")
	trLocal.Connect(trReomteA)
	trReomteA.Connect(trLocal)
	trReomteA.Connect(trReomteB)
	trReomteB.Connect(trReomteC)
	makeRemoteServers([]network.Transport{trReomteA, trReomteB, trReomteC})

	go func() {
		for {
			if err := sendTransaction(trReomteA, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(2 * time.Second)
		}
	}()
	priKey := crypto.GeneratePrivateKey()
	localServer := makeServer("LOCAL", trLocal, &priKey)
	localServer.Start()

}

func makeRemoteServers(trs []network.Transport) {
	for i := 0; i < len(trs); i++ {
		id := fmt.Sprintf("REMOTE_%c", i+'A')
		priKey := crypto.GeneratePrivateKey()
		s := makeServer(id, trs[i], &priKey)
		go s.Start()
	}
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	priKey := crypto.GeneratePrivateKey()
	// data := []byte(strconv.FormatInt(int64(rand.Intn(100000)), 10))
	tx := core.Eggplant{}
	tx.Sign(priKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobEggplantEncoder(buf)); err != nil {
		return err
	}
	msg := network.Message{}
	msg.Header = network.MessageTypeTx
	msg.Data = buf.Bytes()
	return tr.SendMessage(to, msg.Bytes())

}

func makeServer(id string, tr network.Transport, pk *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		PrivateKey: pk,
		ID:         id,
		Transports: []network.Transport{tr},
	}
	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return s
}

// encode => io.Writer  把编码后的数据写入 io.Writer
// decode => io.Reader  把解码后的数据读入 io.Reader
