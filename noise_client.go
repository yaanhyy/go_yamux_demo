package main

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-noise"
	"github.com/libp2p/go-libp2p-core/crypto"
	"net"
	"context"
	"fmt"
)

func NewNoiseTransport(typ, bits int) *noise.Transport {
	priv, pub, err := crypto.GenerateKeyPair(typ, bits)

	if err != nil {
		print(err)
	}
	//localID, err := peer.IDFromPrivateKey(priv)
	id, err := peer.IDFromPublicKey(pub)
	if err != nil {

		print(err)
	} else {
		fmt.Printf("id:%s", id)
	}
	transport,err := noise.New(priv)
	return  transport

}

func newClientConn()(net.Conn) {

	client, clientErr := net.Dial("tcp", "localhost:8981")

	if clientErr != nil {
		print("Failed to accept:", clientErr)
	}



	return client
}

type privKey struct {

}

func (priv privKey)Read(p []byte) (n int, err error) {
	src  := make([]byte, 32)
	for i := 0; i<32; i++ {
		src[i] = byte(i)
	}
	len := copy(p, src)
	return len, nil
}




func main() {
	client := newClientConn()
	initTransport := NewNoiseTransport(crypto.Ed25519, 2048)
	//_, pub, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, privKey{})
	//id, err := peer.IDFromPublicKey(pub)
	//secureConn, err := initTransport.SecureOutbound(context.TODO(), client, id)
	id, _ := peer.Decode("12D3KooWA4Xop1JaT3MHxwYMkCepYsv4iPVopMXwCz5iHYdBfeSB")
	secureConn, err := initTransport.SecureOutbound(context.TODO(), client, id)
	if err == nil {
		before := []byte("hello world")
		secureConn.Write(before)
	} else {
		fmt.Printf("%+v", err)
	}
}
