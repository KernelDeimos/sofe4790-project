package main

import "github.com/KernelDeimos/sofe4790/singlenode"
import "github.com/sirupsen/logrus"

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	singlenode.RunApplication()
}
