package main

import (
	"context"
	"github.com/hashicorp/yamux"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-secio"
	"net"
"time"
	ci "github.com/libp2p/go-libp2p-core/crypto"
)

func newClientConn()(net.Conn) {

	client, clientErr := net.Dial("tcp", "localhost:8981")

	if clientErr != nil {
		print("Failed to accept:", clientErr)
	}



	return client
}

func NewTransport(typ, bits int) *secio.Transport {
	priv, pub, err := ci.GenerateKeyPair(typ, bits)
	if err != nil {
		print(err)
	}
	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		print(err)
	}
	return &secio.Transport{
		PrivateKey: priv,
		LocalID:    id,
	}
}


func main()  {

	client := newClientConn()
	clientTpt := NewTransport(ci.Secp256k1, 2048)
	clientConn, _ :=clientTpt.SecureOutbound(context.TODO(), client, "")

	session, _ := yamux.Client(clientConn, nil)

	// 建立应用流通道1
	stream, _ := session.Open()
	stream.Write([]byte("ping" ))
	stream.Write([]byte("pnng" ))
	time.Sleep(1 * time.Second)

	// 建立应用流通道2
	stream1, _ := session.Open()
	stream1.Write([]byte("pong"))
	time.Sleep(1 * time.Second)

	// 清理退出
	time.Sleep(5 * time.Second)
	stream.Close()
	stream1.Close()
	session.Close()

}
