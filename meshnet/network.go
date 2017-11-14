package meshnet

import "net"

type PeerNode struct {
	Host string
	Port int
}

type PeerConnection struct {
	Node   PeerNode
	Socket *net.TCPConn
}

type Network struct{}

func (net *Network) Join(peers []PeerNode) {
	//
}
