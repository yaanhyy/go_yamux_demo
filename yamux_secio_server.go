package main
// 多路复用
import (
	"context"
	"fmt"
	"github.com/hashicorp/yamux"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-secio"
	"net"
	"time"
	ci "github.com/libp2p/go-libp2p-core/crypto"
)

func Recv(session *yamux.Session, stream net.Conn, id int){
	for {
		buf := make([]byte, 50)
		n, err := stream.Read(buf)
		if err == nil{
			fmt.Println("ID:", id, ", len:", n, time.Now().Unix(), string(buf))
			stream1, _ := session.Open()
			stream1.Write([]byte("pong"))
		}else{
			fmt.Println(time.Now().Unix(), err)
			return
		}
	}
}

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
	// 建立底层复用连接
	server := newServerConn()
	serverTpt := NewTransport(ci.Ed25519, 2048)
	serverConn, _ := serverTpt.SecureInbound(context.TODO(), server)
	//tcpaddr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:8980");
	//tcplisten, _ := net.ListenTCP("tcp", tcpaddr);
	//conn, _ := tcplisten.Accept()
	session, _ := yamux.Server(serverConn, nil)

	id :=1
	for {
		// 建立多个流通路
		stream, err := session.Accept()
		if err == nil {
			fmt.Println("accept")
			id ++
			go Recv(session, stream, id)
		}else{
			fmt.Println("session over.", err)
			return
		}
	}

}

