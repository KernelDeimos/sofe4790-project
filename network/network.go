package network

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	MessageTypeElection    = "_elect"
	MessageTypeCoordinator = "_coord"
)

type Network struct {
	Self          NetworkPeer
	ElectionState bool
	Peers         []NetworkPeer
	Timeout       time.Duration
}

func NewDefaultNetwork(host string, id, port int, leader bool) *Network {
	self := NetworkPeer{}
	self.Host = host
	self.Identifier = id
	self.Port = port
	self.IsLeader = leader

	return &Network{
		Self:          self,
		ElectionState: false,
		Peers:         []NetworkPeer{},
		Timeout:       time.Second * 10,
	}
}

func (n *Network) AddPeer(peer NetworkPeer) {
	n.Peers = append(n.Peers, peer)
}

func (n *Network) SetLeader(leaderID int) {
	if n.Self.Identifier == leaderID {
		n.Self.IsLeader = true
	} else {
		for i := range n.Peers {
			if n.Peers[i].Identifier == leaderID {
				n.Peers[i].IsLeader = true
				break
			}
		}
	}
}

func (n *Network) StartElection() {
	n.ElectionState = true

	wg := &sync.WaitGroup{}
	anyAnswers := false

	// Clear the IsLeader attribute of any nodes who are currently a leader
	for _, peer := range n.Peers {
		peer.IsLeader = false
	}

	// Send an election message to all nodes with a higher identifier
	for _, peer := range n.Peers {
		go func() {
			if peer.Identifier > n.Self.Identifier {
				_, err := peer.Send(MessageTypeElection, "", n.Timeout)
				if err == nil {
					anyAnswers = true
				}
				wg.Done()
			}
		}()
	}

	// Wait for all nodes to either answer or timeout
	wg.Wait()

	// Set self to leader if there were no answers (all timeouts)
	n.Self.IsLeader = !anyAnswers

	if n.Self.IsLeader {
		// Send a coordinator message to all nodes
		for _, peer := range n.Peers {
			peer.Send(MessageTypeCoordinator, "", n.Timeout)
		}
	}

}

func (n *Network) Attach(r *gin.Engine) {
	r.GET("/:type/:data", func(c *gin.Context) {
		logrus.Infof("Received message of type %s\n", c.Param("type"))
		switch c.Param("type") {
		case MessageTypeElection:
			c.JSON(http.StatusOK, struct {
				Status string `json:"status"`
			}{"okay"})
		case MessageTypeCoordinator:
			//
		}
	})
}

type NetworkPeer struct {
	Identifier int
	IsLeader   bool
	PeerNode
}

func NewPeer(host string, port, id int) NetworkPeer {
	return NetworkPeer{
		id, false, PeerNode{
			Host: host,
			Port: port,
		},
	}
}

type PeerNode struct {
	Host string
	Port int
}

func (node *PeerNode) generateURI(mtype string, data interface{}) (string, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	dataString := url.PathEscape(string(dataBytes))

	route := "/" + mtype + "/" + dataString
	host := node.Host + ":" + strconv.Itoa(node.Port)

	uri := host + route

	return uri, nil
}

func (node *PeerNode) Send(
	mtype string, data interface{}, timeout time.Duration,
) (string, error) {
	uri, err := node.generateURI(mtype, data)
	if err != nil {
		return "", err
	}

	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(uri)
	if err != nil {
		return "", err
	}
	replyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	replyString := string(replyBytes)

	return replyString, nil
}
