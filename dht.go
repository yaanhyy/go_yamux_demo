package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-conn-security-multistream"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"

	//"github.com/whyrusleeping/go-logging"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-secio"
	"github.com/libp2p/go-tcp-transport"
	"net"

	ci "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/metrics"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	yamux "github.com/libp2p/go-libp2p-yamux"
	msmux "github.com/libp2p/go-stream-muxer-multistream"
	"github.com/libp2p/go-libp2p-testing/net"
	"github.com/libp2p/go-libp2p-core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
	log2 "github.com/ipfs/go-log/v2"
	logging "github.com/ipfs/go-log"

)
var logger     = logging.Logger("dht")
var testPrefix = dht.ProtocolPrefix("/test")

type blankValidator struct{}

func (blankValidator) Validate(_ string, _ []byte) error        { return nil }
func (blankValidator) Select(_ string, _ [][]byte) (int, error) { return 0, nil }


func setupDHT(ctx context.Context,  client bool, options ...dht.Option) *dht.IpfsDHT {
	baseOpts := []dht.Option{
		testPrefix,
		dht.NamespacedValidator("v", blankValidator{}),
		dht.DisableAutoRefresh(),
	}

	if client {
		baseOpts = append(baseOpts, dht.Mode(dht.ModeClient))
	} else {
		baseOpts = append(baseOpts, dht.Mode(dht.ModeServer))
	}

	d, err := dht.New(
		ctx,
		bhost.New(GenSwarm(ctx, false)),
		append(baseOpts, options...)...,
	)
	if err != nil {
		fmt.Printf("dht new err:%+v", err)
	}
	return d
}

// GenUpgrader creates a new connection upgrader for use with this swarm.
func GenUpgrader(n *swarm.Swarm) *tptu.Upgrader {
	id := n.LocalPeer()
	pk := n.Peerstore().PrivKey(id)
	secMuxer := new(csms.SSMuxer)
	secMuxer.AddTransport(secio.ID, &secio.Transport{
		LocalID:    id,
		PrivateKey: pk,
	})

	stMuxer := msmux.NewBlankTransport()
	stMuxer.AddTransport("/yamux/1.0.0", yamux.DefaultTransport)

	return &tptu.Upgrader{
		Secure:  secMuxer,
		Muxer:   stMuxer,
		Filters: n.Filters,
	}

}

// GenSwarm generates a new test swarm.
func GenSwarm(ctx context.Context, dial_only bool) *swarm.Swarm {
	priv, pub, err := ci.GenerateKeyPair(ci.Ed25519, 2048)
	if err != nil {
		print(err)
	}
	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		print(err)
	}

	var p tnet.PeerNetParams
	p.PrivKey = priv
	p.PubKey = pub
	p.ID = id
	p.Addr, err = ma.NewMultiaddr("/ip4/127.0.0.1/tcp/5679") //tnet.ZeroLocalTCPAddress
	if err != nil {
		print(err)
	}
	ps := pstoremem.NewPeerstore()
	ps.AddPubKey(p.ID, p.PubKey)
	ps.AddPrivKey(p.ID, p.PrivKey)
	s := swarm.NewSwarm(ctx, p.ID, ps, metrics.NewBandwidthCounter())

	tcpTransport := tcp.NewTCPTransport(GenUpgrader(s))


	if err := s.AddTransport(tcpTransport); err != nil {
		print(err)
	}

	if !dial_only {
		if err := s.Listen(p.Addr); err != nil {
			print(err)
		}

		s.Peerstore().AddAddrs(p.ID, s.ListenAddresses(), peerstore.PermanentAddrTTL)
	}

	return s
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


func main() {
	log2.SetAllLoggers(log2.LevelDebug)
	logger.Debugf("refreshing DHTs routing tables...")
	//server := newServerConn()
	//serverTpt := NewTransport(ci.Ed25519, 2048)
	//serverTpt.SecureInbound(context.TODO(), server)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupDHT(ctx ,  false)
	dht:=setupDHT(ctx ,  false)
	dht.PeerID()
	//dht.Process()
	select {

	}
}