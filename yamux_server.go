package main
// 多路复用
import (
	"fmt"
	"github.com/hashicorp/yamux"
	"net"
	"time"
)

func Recv(session *yamux.Session,stream net.Conn, id int){
	for {
		buf := make([]byte, 15)
		n, err := stream.Read(buf)
		if err == nil{
			fmt.Println("ID:", id, ", len:", n, time.Now().Unix(), string(buf))
			stream.Write([]byte("server ok" ))
			stream1, _ := session.Open()
			stream1.Write([]byte("pong"))
		}else{
			fmt.Println(time.Now().Unix(), err)
			return
		}
	}
}
func main()  {
	// 建立底层复用连接
	tcpaddr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:8980");
	tcplisten, _ := net.ListenTCP("tcp", tcpaddr);
	conn, _ := tcplisten.Accept()
	session, _ := yamux.Server(conn, nil)

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
