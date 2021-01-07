package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/betom84/kl-logger/api"
	"github.com/betom84/kl-logger/klimalogg"
	"github.com/betom84/kl-logger/repository"
	"github.com/betom84/kl-logger/transceiver"
	"github.com/sirupsen/logrus"
)

var (
	log      *string = flag.String("log", "stdout", "Logfile")
	logLevel *string = flag.String("logLevel", "info", "Log level (e.g. error, info, debug, trace)")
	apiPort  *int    = flag.Int("apiPort", 8088, "Port to serve http api requests")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage of kl-logger:")
		flag.PrintDefaults()
	}

	flag.Parse()

	ll, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logrus.WithError(err).Panic()
	}

	logrus.SetLevel(ll)

	if *log != "stdout" {
		lf, err := os.OpenFile(*log, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			logrus.WithError(err).Panic()
		}

		logrus.SetOutput(lf)
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logrus.Infof("log level is set to '%s'", logrus.GetLevel().String())
}

func main() {
	err := run()
	if err != nil {
		logrus.WithError(err).Panic()
	}

	os.Exit(0)
}

func run() error {
	t := transceiver.NewTransceiver(0x6666, 0x5555)
	err := t.Open()
	if err != nil {
		return err
	}

	defer func() {
		err = t.Close()
		if err != nil {
			logrus.WithError(err).Error("could not close transceiver")
		}

		logrus.Debug("usb transceiver disconnected")
	}()

	logrus.WithFields(logrus.Fields{
		"vendorID":  fmt.Sprintf("0x%04x", t.VendorID),
		"productID": fmt.Sprintf("0x%04x", t.ProductID),
	}).Info("usb transceiver connected")

	c := klimalogg.NewConsole(t)
	err = c.Initialise(&repository.Default)
	if err != nil {
		return err
	}

	defer func() {
		c.Close()
		logrus.Debug("klimalogg console closed")
	}()

	c.StartCommunication()
	if err != nil {
		return err
	}

	logrus.Info("klimalogg console ready")

	server := api.NewServer(&repository.Default)
	go func() {
		logrus.Infof("start api http server on :%d", *apiPort)

		err := http.ListenAndServe(fmt.Sprintf(":%d", *apiPort), server)
		if err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Error("error running api http server")
		}
	}()

	waitForSignal()

	return nil
}

func waitForSignal() bool {
	var sigChan = make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGINT)

	logrus.Info("waiting for signal to exit")

	sig := <-sigChan
	logrus.Infof("exiting due to %+v\n", sig)

	return true
}
