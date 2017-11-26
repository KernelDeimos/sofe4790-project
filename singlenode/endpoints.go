package singlenode

import (
	"os"
	"sync"

	"github.com/KernelDeimos/sofe4790/configobjects"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/sirupsen/logrus"
)

type StreamEndpoint struct {
	Filename string
	stream   *os.File
}

func NewStreamEndpoint(filename string) *StreamEndpoint {
	return &StreamEndpoint{
		filename, nil,
	}
}

func (end *StreamEndpoint) Start(name string) (*sync.WaitGroup, error) {
	var err error
	end.stream, err = os.OpenFile(
		end.Filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644,
	)
	return nil, err
}

func (end *StreamEndpoint) Catch(source string, data map[string]string) {
	payload := data["payload"]
	logrus.Debug("Writing data to stream: ", payload)
	_, err := end.stream.WriteString(payload)
	if err != nil {
		logrus.Error(err)
	}

	err = end.stream.Sync()
	if err != nil {
		logrus.Error(err)
	}
}

type TweetEndpoint struct {
	client *twitter.Client
}

func NewTweetEndpoint(conf configobjects.OAuth) *TweetEndpoint {
	config, token := conf.GetOAuth1()
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	return &TweetEndpoint{
		client: client,
	}
}

func (end *TweetEndpoint) Start(name string) (*sync.WaitGroup, error) {
	return nil, nil
}

func (end *TweetEndpoint) Catch(source string, data map[string]string) {
	payload := data["payload"]

	logrus.Debug("Writing data to twitter: ", payload)

	_, _, err := end.client.Statuses.Update(payload, nil)
	if err != nil {
		logrus.Error(err)
	}
}
