package main

import (
	"github.com/KirillMironov/kappa/internal/kappctl/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	// Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01|02 15:04:05.000",
	})

	// Cobra
	err := cmd.NewDefaultCommand().Execute()
	if err != nil {
		logger.Fatal(err)
	}
}
