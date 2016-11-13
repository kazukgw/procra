package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/elazarl/goproxy"
	"github.com/sevlyar/go-daemon"
)

var (
	signal = flag.String("s", "", `send signal to the daemon
		quit - graceful shutdown
		stop - fast shutdown
		reload - reloading the configuration file`)
)

func main() {
	flag.Parse()
	daemon.AddCommand(
		daemon.StringFlag(signal, "quit"),
		syscall.SIGQUIT,
		termHandler,
	)
	daemon.AddCommand(
		daemon.StringFlag(signal, "stop"),
		syscall.SIGTERM,
		termHandler,
	)
	daemon.AddCommand(
		daemon.StringFlag(signal, "reload"),
		syscall.SIGHUP,
		reloadHandler,
	)

	ctx := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "proxy.log",
		LogFilePerm: 0644,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"http-proxy"},
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := ctx.Search()
		if err != nil {
			log.Fatalln("Unable send signal to the daemon:", err)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := ctx.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer ctx.Release()

	log.Println("daemon started")

	go worker()

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println("daemon terminated")
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func worker() {
	proxy := goproxy.NewProxyHttpServer()
	log.Fatal(http.ListenAndServe(":8080", proxy))
	for {
		select {
		case <-stop:
			break
		}
	}
	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Println("configuration reloaded")
	return nil
}
