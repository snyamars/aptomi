package main

import (
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/Aptomi/aptomi/cmd/aptomi/root"
	"github.com/sirupsen/logrus"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	defer func() {
		if err := recover(); err != nil {
			logrus.Info(string(debug.Stack()))
			logrus.Fatalf("%s", err) // this will terminate the server
		}
	}()

	if err := root.Command.Execute(); err != nil {
		logrus.Fatalf("%s", err) // this will terminate the server
	}
}
