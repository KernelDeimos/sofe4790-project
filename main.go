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

	errs := &estate.ErrorState{}

	host := args[0]
	port := strparse.ParseI(errs, args[1])
	id := strparse.ParseI(errs, args[2])

	if errs.GetError() != nil {
		logrus.Fatal(errs)
	}

	logrus.SetLevel(logrus.DebugLevel)
	singlenode.RunApplication(host, port, id)
}
