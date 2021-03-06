package main

import (
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/Aptomi/aptomi/cmd/aptomictl/root"
	"github.com/sirupsen/logrus"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	defer func() {
		if err := recover(); err != nil {
			logrus.Info(string(debug.Stack()))
			logrus.Fatalf("%s", err) // this will terminate the client
		}
	}()

	if err := root.Command.Execute(); err != nil {
		logrus.Fatalf("%s", err) // this will terminate the client
	}
}
