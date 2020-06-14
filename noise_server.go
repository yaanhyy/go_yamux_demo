package main

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"net"
	"github.com/libp2p/go-libp2p-core/crypto"
	"context"
	"github.com/libp2p/go-libp2p-noise"
	"io"
	"github.com/libp2p/go-libp2p-core/sec"
	"fmt"
)

func newServerConn()(net.Conn) {
	lstnr, err := net.Listen("tcp", "localhost:8981")
	if err != nil {

		return nil
	}

	server, err := lstnr.Accept()

	lstnr.Close()

	if err != nil {
		print("Failed to accept:", err)
	}



	return server
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



func NewNoiseTransport(typ, bits int) *noise.Transport {
	priv, pub, err := crypto.GenerateKeyPairWithReader(typ, bits, privKey{})
	//priv, _, err := crypto.GenerateKeyPair(typ, bits)
	if err != nil {
		print(err)
	}
	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		print(err)
	} else {
		fmt.Printf("id:%s", id)
	}

	transport,err := noise.New(priv)
	return  transport

}

func ReadBuf( serverConn sec.SecureConn) {

	after := make([]byte, 100)
	_, err := io.ReadFull(serverConn, after)
	if err != nil {
		print(err.Error())
	}
	fmt.Printf("%s", after)
}

func main()  {
	serverTpt := NewNoiseTransport(crypto.Ed25519, 2048)
	server := newServerConn()

	serverConn, _ := serverTpt.SecureInbound(context.TODO(), server)
	ReadBuf(serverConn)
}