package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"flag"
	"log"
)

var addr = flag.String("addr", ":8000", "http service address")

func handler(w http.ResponseWriter, r *http.Request) {
  waitReqest := make(chan string)
  waitAsync := make(chan bool)
  waitAsyncFail := make(chan bool)

  go (func () {
	println("async save success")

	waitAsync <- true
  })()

  go (func () {
	println("async save fail")

	waitAsyncFail <- false
  })()

  go (func () {
	resp, _ := http.Get("http://baobab.kaiyanapp.com/api/v1/feed.bak")
	body, _ := ioutil.ReadAll(resp.Body)
	respBody := string(body)

	waitReqest <- respBody
  })()

  var status bool
  responseData := <- waitReqest

  select {
  case status = <- waitAsync :
  case status = <- waitAsyncFail :
  default:
	println("nothing")
  }

  println(status)
  // _ = <- waitAsync
  fmt.Fprint(w, responseData)
  // fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	hub := newHub()

	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/crawl", handler)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	  serveWs(hub, w, r)
	})

	server := &http.Server{Addr: *addr, Handler: nil}
	err := server.ListenAndServe()
	// http.ListenAndServe(*addr, nil)
	if err != nil {
	  log.Fatal("ListenAndServe: ", err)
	} else {
	  fmt.Printf(" 🔭 server run listen: %s", *addr)
	}
}
