// Copyright 2016 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

var (
	flagRootDir     = flag.String("dir", ".", "root dir")
	flagHttpAddr    = flag.String("http", ":8080", "HTTP service address")
	flagOpenBrowser = flag.Bool("openbrowser", true, "open browser automatically")
)

func main() {
	flag.Parse()

	host, port, err := net.SplitHostPort(*flagHttpAddr)
	if err != nil {
		log.Fatal(err)
	}

	if host == "" {
		host = getLocalIP()
	}
	httpAddr := host + ":" + port
	url := "http://" + httpAddr + "/_book"

	go func() {
		log.Printf("dir: %s\n", *flagRootDir)

		if waitServer(url) && *flagOpenBrowser && startBrowser(url) {
			log.Printf("A browser window should open. If not, please visit %s", url)
		} else {
			log.Printf("Please open your web browser and visit %s", url)
		}
		log.Printf("Hit CTRL-C to stop the server\n")
	}()

	log.Fatal(http.ListenAndServe(httpAddr, http.FileServer(http.Dir(*flagRootDir))))
}

// waitServer waits some time for the http Server to start
// serving url. The return value reports whether it starts.
func waitServer(url string) bool {
	tries := 20
	for tries > 0 {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
		tries--
	}
	return false
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// startBrowser tries to open the URL in a browser, and returns
// whether it succeed.
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
