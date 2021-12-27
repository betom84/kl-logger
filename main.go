package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/betom84/kl-logger/api"
	"github.com/betom84/kl-logger/klimalogg"
	"github.com/betom84/kl-logger/metrics"
	"github.com/betom84/kl-logger/repository"
	"github.com/sirupsen/logrus"

	_ "net/http/pprof"
)

var (
	log      *string = flag.String("log", "stdout", "Logfile")
	logLevel *string = flag.String("logLevel", "info", "Log level (e.g. error, info, debug, trace)")
	usbTrace *bool   = flag.Bool("usbTrace", false, "Trace usb control messages")
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
	t := klimalogg.NewTransceiver(0x6666, 0x5555)
	err := t.Open()
	if err != nil {
		return err
	}

	if *usbTrace {
		trace, err := os.OpenFile(fmt.Sprintf("%s_transceiver.trace", time.Now().Format("20060102_15040507")), os.O_CREATE|os.O_RDWR, 0666)
		if err == nil {
			t.StartTracing(trace)
			defer trace.Close()
		}
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

	c, err := klimalogg.NewConsole(t)
	if err != nil {
		return err
	}

	defer func() {
		c.Close()
		logrus.Debug("klimalogg console closed")
	}()

	c.AddListener(repository.Default.NewListener())
	c.AddListener(metrics.KlimaloggCurrentValuesPublisher())
	c.StartCommunication()
	if err != nil {
		return err
	}

	logrus.Info("klimalogg console ready")

	server := api.NewServer(repository.Default, c.Transceiver())
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
	var sigChan = make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGINT)

	logrus.Info("waiting for signal to exit")

	sig := <-sigChan
	logrus.Infof("exiting due to %+v\n", sig)

	return true
}
