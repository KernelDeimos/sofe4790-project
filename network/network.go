package network

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

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
