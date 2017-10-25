package main

import (
		"fmt"
		"net/http"
		"io/ioutil"
		"flag"
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

func main() {
		flag.Parse()
		fmt.Printf(" 🔭 server run listen: %s \n", *addr)
		http.HandleFunc("/", handler)
		http.ListenAndServe(*addr, nil)
}
