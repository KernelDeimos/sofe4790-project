package singlenode

import (
	"bufio"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type StreamSource struct {
	Stream *os.File
}

func (src *StreamSource) Start(name string, c Collector) (*sync.WaitGroup, error) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	reader := bufio.NewReader(src.Stream)

	go func() {
		for {
			input, err := reader.ReadString('\n')

			logrus.Debugf("Received input: '%s'", input)

			if err != nil {
				logrus.Error(err)
				continue
			}

			data := map[string]string{}
			data["payload"] = input

			c.Collect(Event{
				Key:  name,
				Data: data,
			})
		}
	}()

	return wg, nil
}

func NewStdinSource() *StreamSource {
	return &StreamSource{os.Stdin}
}
