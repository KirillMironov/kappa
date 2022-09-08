package main

import (
	"github.com/KirillMironov/kappa/internal/kapplet/transport"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01|02 15:04:05.000",
	})

	var handler = transport.NewHandler("20501", logger)

	err := handler.Start()
	if err != nil {
		logger.Fatal(err)
	}
}
