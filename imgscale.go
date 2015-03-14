package main

import (
	"flag"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/vanng822/imgscale/imgscale"
	provider_http "github.com/vanng822/imgscale/provider/http"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(2)
	var (
		configPath string
		host       string
		port       int
	)

	flag.StringVar(&host, "h", "127.0.0.1", "Host to listen on")
	flag.IntVar(&port, "p", 9999, "Port number to listen on")
	flag.StringVar(&configPath, "c", "./config/formats.json", "Path to configurations")
	flag.Parse()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2)
	app := martini.Classic()
	handler := imgscale.Configure(configPath)
	defer handler.Cleanup()
	handler.SetImageProvider(provider_http.New("http://imgscale.isgoodness.com/getimg/"))
	app.Use(handler.ServeHTTP)
	app.Get("/", indexHandler)
	app.Get("/getimg/(?P<url>.+)", getimgHandler)
	log.Printf("listening to address %s:%d", host, port)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), app)
	
	for {
		sig := <-sigc
		switch sig {
		case syscall.SIGUSR2:
			log.Println("Reloading config")
			handler.Reload()
		default:
			log.Printf("Got signal: %s", sig)
			return
		}
	}
	
}
