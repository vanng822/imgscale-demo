package main

import (
	"flag"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/vanng822/gopid"
	"github.com/vanng822/imgscale/imgscale"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"io/ioutil"
)

func main() {
	runtime.GOMAXPROCS(2)
	var (
		configPath string
		host       string
		port       int
		pidFile    string
		force      bool
	)

	flag.StringVar(&host, "h", "127.0.0.1", "Host to listen on")
	flag.IntVar(&port, "p", 8080, "Port number to listen on")
	flag.StringVar(&configPath, "c", "./config/formats.json", "Path to configurations")
	flag.StringVar(&pidFile, "pid", "imgscale.pid", "Pid file")
	flag.BoolVar(&force, "f", false, "Force and remove pid file")
	flag.Parse()

	if pidFile != "" {
		gopid.CheckPid(pidFile, force)
		gopid.CreatePid(pidFile)
		defer gopid.CleanPid(pidFile)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2)
	app := martini.Classic()
	handler := imgscale.Configure(configPath)
	defer handler.Cleanup()
	handler.SetImageProvider(imgscale.NewImageProviderHTTP("http://imgscale.isgoodness.com/getimg/"))
	app.Use(handler.ServeHTTP)
	app.Get("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(`
			<html>
				<head>
				</head>
				<body>
				<div>
				<p>
					Demo server for https://github.com/vanng822/imgscale
				</p>
				<p>
				<a href="/img/0x360/http://images4.fanpop.com/image/photos/16100000/Cute-Kitten-kittens-16123796-1280-800.jpg"><img src="/img/100x0/http://images4.fanpop.com/image/photos/16100000/Cute-Kitten-kittens-16123796-1280-800.jpg" /></a>
				</p>
				</div>
				</body>
			</html>`))
	})

	app.Get("/getimg/(?P<url>.+)", func(res http.ResponseWriter, req *http.Request, params martini.Params) int {
		url := params["url"]
		resp, err := http.Get(url)
		if err != nil {
			return http.StatusNotFound
		}
		defer resp.Body.Close()

		imgData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return http.StatusNotFound
		}
		res.Write(imgData)
		
		return http.StatusOK
	})
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
