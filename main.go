package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/agnivade/mdns"
	"github.com/syntaqx/serve"
)

var (
	host   string
	port   string
	dir    string
	iface  string
	fName  string
	prefix string
)

func main() {
	flag.StringVar(&host, "host", "", "host address to bind to")
	flag.StringVar(&port, "port", "8080", "listening port")
	flag.StringVar(&dir, "dir", "", "directory path to serve")
	flag.StringVar(&iface, "iface", "wlp4s0", "interface to publish service info")
	flag.StringVar(&fName, "friendlyName", "Gompa1", "A friendly name to identify this service")
	flag.StringVar(&prefix, "prefix", "so-1", "Playlist prefix")
	flag.Parse()

	logger := log.New(os.Stdout, "[serve] ", log.LstdFlags|log.Lshortfile)

	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}

	fs := serve.NewFileServer(serve.Options{
		Directory: dir,
	})

	fs.Use(
		Logger(logger),
		Recover(),
		CORS(),
	)

	addr := net.JoinHostPort(host, port)
	server := &http.Server{
		Addr:         addr,
		Handler:      fs,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	iPort, err := strconv.Atoi(port)
	if err != nil {
		logger.Fatal(err)
	}

	mServer, err := mdns.Publish(iface, iPort, "stream_publisher._tcp", fName+"_"+prefix)
	if err != nil {
		logger.Fatal(err)
	}

	// listen for signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	logger.Printf("Started server at %s\n", port)

	// Block until one of the signals above is received
	<-signalCh
	logger.Println("Quit signal received, initializing shutdown...")
	logger.Println("Stopping HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = server.Shutdown(ctx)
	if err != nil {
		logger.Println(err)
	}
	cancel()

	logger.Println("Stopping mDNS service")
	err = mServer.Shutdown()
	if err != nil {
		logger.Println(err)
	}
}
