package singlenode

import (
	"os"
	"sync"

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
