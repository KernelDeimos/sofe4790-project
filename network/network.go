package network

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	MessageTypeElection    = "_elect"
	MessageTypeCoordinator = "_coord"
	MessageTypeEvent       = "_event"
)

type Network struct {
	Self         NetworkPeer
	ElectionWait *sync.WaitGroup
	Peers        []NetworkPeer
	Timeout      time.Duration
	Messages     chan interface{}
}

func NewDefaultNetwork(host string, id, port int, leader bool) *Network {
	self := NetworkPeer{}
	self.Host = host
	self.Identifier = id
	self.Port = port
	self.IsLeader = leader

	return &Network{
		Self:         self,
		ElectionWait: nil,
		Peers:        []NetworkPeer{},
		Timeout:      time.Second * 10,
		Messages:     make(chan interface{}, 2000),
	}
}

func (n *Network) StartService() {
	for {
		data := <-n.Messages
		logrus.Info("Got a message! ", data)

		for i := range n.Peers {
			if n.Peers[i].IsLeader {
				_, err := n.Peers[i].Send(MessageTypeEvent, data, n.Timeout)
				if err != nil {
					logrus.Error(err)
					logrus.Error("Connection to leader failed; starting election.")
					n.Messages <- data
					n.StartElection()
					n.ElectionWait.Wait()
					n.ElectionWait = nil
				}
			}
		}
	}
}

func (n *Network) AddPeer(peer NetworkPeer) {
	n.Peers = append(n.Peers, peer)
}

func (n *Network) SetLeader(leaderID int) {
	// Clear the IsLeader attribute of any nodes who are currently a leader
	for i := range n.Peers {
		n.Peers[i].IsLeader = false
	}
	n.Self.IsLeader = false

	if n.Self.Identifier == leaderID {
		n.Self.IsLeader = true
	}
	for i := range n.Peers {
		if n.Peers[i].Identifier == leaderID {
			n.Peers[i].IsLeader = true
			break
		}
	}
}

func (n *Network) SendToLeader(data interface{}) bool {
	if n.Self.IsLeader {
		return false
	}
	n.Messages <- data
	return true
}

func (n *Network) StartElection() {
	n.ElectionWait = &sync.WaitGroup{}
	n.ElectionWait.Add(1)

	wg := &sync.WaitGroup{}
	anyAnswers := false

	// Send an election message to all nodes with a higher identifier
	for _, peer := range n.Peers {
		if peer.Identifier > n.Self.Identifier {
			wg.Add(1)
			go func() {
				_, err := peer.Send(MessageTypeElection, "", n.Timeout)
				if err == nil {
					anyAnswers = true
				}
				wg.Done()
			}()
		}
	}

	// Wait for all nodes to either answer or timeout
	wg.Wait()

	if !anyAnswers {
		logrus.Info("Sending coordinator message to other nodes.")
		// Set self to leader if there were no answers (all timeouts)
		//n.SetLeader(n.Self.Identifier)
		//n.ElectionWait.Done()
		// Send a coordinator message to all nodes
		for _, peer := range n.Peers {
			peer.Send(MessageTypeCoordinator, n.Self.Identifier, n.Timeout)
		}
	}

}

func (n *Network) Attach(r *gin.Engine, onEvent func(data interface{})) {
	r.GET("/:type/:data", func(c *gin.Context) {
		logrus.Infof("Received message of type %s\n", c.Param("type"))
		switch c.Param("type") {
		case MessageTypeElection:
			c.JSON(http.StatusOK, struct {
				Status string `json:"status"`
			}{"okay"})
		case MessageTypeCoordinator:
			jsonString := c.Param("data")
			var leaderID int
			err := json.Unmarshal([]byte(jsonString), &leaderID)
			if err != nil {
				logrus.Fatal("Failed to get coordinator ID: ", err)
			}
			n.SetLeader(leaderID)
			spew.Dump(n.Peers)
			if n.ElectionWait != nil {
				n.ElectionWait.Done()
			}
		case MessageTypeEvent:
			jsonString := c.Param("data")
			onEvent(jsonString)
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
	host := "http://" + node.Host + ":" + strconv.Itoa(node.Port)

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
