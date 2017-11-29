package main

import (
	"os"

	"github.com/KernelDeimos/sofe4790/estate"
	"github.com/KernelDeimos/sofe4790/singlenode"
	"github.com/KernelDeimos/sofe4790/strparse"
	"github.com/sirupsen/logrus"
)

func main() {
	args := os.Args[1:]

	if len(args) < 4 {
		logrus.Fatal("Usage: ./node <host> <port> <id> <leader>")
	}

	errs := &estate.ErrorState{}

	host := args[0]
	port := strparse.ParseI(errs, args[1])
	id := strparse.ParseI(errs, args[2])
	leader := strparse.ParseI(errs, args[3])

	if errs.GetError() != nil {
		logrus.Fatal(errs)
	}

	logrus.SetLevel(logrus.DebugLevel)
	singlenode.RunApplication(host, port, id, leader)
}
