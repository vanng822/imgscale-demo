package main

import (
	"github.com/go-martini/martini"
	"net/http"
	"io/ioutil"
)

func getimgHandler(res http.ResponseWriter, req *http.Request, params martini.Params) int {
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
}

func indexHandler(res http.ResponseWriter, req *http.Request) string {
	return `<html>
				<head>
				</head>
				<body>
					<div style="text-align: center;">
					<p>
						Demo server for <a href="https://github.com/vanng822/imgscale" target="_blank">https://github.com/vanng822/imgscale</a>
					</p>
					<p>
						<a href="/img/0x360/http://images4.fanpop.com/image/photos/16100000/Cute-Kitten-kittens-16123796-1280-800.jpg" target="_blank">
							<img src="/img/100x0/http://images4.fanpop.com/image/photos/16100000/Cute-Kitten-kittens-16123796-1280-800.jpg" />
						</a>
					</p>
					</div>
				</body>
			</html>`
}
