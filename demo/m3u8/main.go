package main

import (
	"strings"
	"net/http"
	_ "os"
	_ "errors"
	_ "io"
	_ "bufio"
	"io/ioutil"
	_ "regexp"
	"flag"
	"github.com/tidwall/gjson"
)

var mode = flag.String("mode", "simple", "input m3u8 url")
var info = flag.String("info", "", "url for simple mode")

type DownResult struct {
	content []byte
	pos int
}

var out = make(chan []byte)
var downChan = make(chan DownResult)
var urlChan = make(chan string)

func downLoadFile(pos int, url string){
	resp, err := http.Get(url)

	if err != nil {
		println(err.Error())
		out <- []byte{}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	println("downloading", pos, len(body))
	downChan <- DownResult{
		content: body,
		pos: pos,
	}
}

func fetchMeta(url string){
	resp, err := http.Get(url)

	if err != nil {
		println(err.Error())
		out <- []byte{}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	out <- body
}

func saveFile(bytesData []byte){
	outputFile := "movie.ts"
	err := ioutil.WriteFile(outputFile, bytesData, 0644)
	if err != nil {
		println(err.Error())
	}
}

func loadVideoUrlFromRemote() {
	configBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("config not exists")
	}
	urlValue := gjson.Get(string(configBytes), "location.0")
	realUrl := urlValue.String()

	go fetchMeta(realUrl)

	content := <- out
	jsonStr := string(content)

	value := gjson.Get(jsonStr, "segs.MP4-HD-Mobile-M3u8.0.url")
	urlChan <- value.String()
}

func loadVideoUrlSimple() string {
	return *info
}

func main (){
	flag.Parse()

	var m3u8Url = ""
	if *mode == "simple" {
		m3u8Url = loadVideoUrlSimple()
	} else {
		go loadVideoUrlFromRemote()
		m3u8Url = <- urlChan
	}

	println(m3u8Url)
	go fetchMeta(m3u8Url)
	m3u8Content := <- out
	contentList := strings.Split(string(m3u8Content), "\r\n")

	videoUrls := []string{}

	for _, contentLine := range contentList {
		if strings.Contains(contentLine, "http") {
			videoUrls = append(videoUrls, strings.Trim(contentLine, " "))
		}
	}

	for i, vl := range videoUrls {
		go downLoadFile(i, vl)
	}

	urlLen := len(videoUrls)
	videoBytes := make([][]byte, urlLen)
	for _, _ = range videoUrls {
		videoContent := <- downChan
		videoBytes[videoContent.pos] = videoContent.content
	}

	finalBytes := []byte{}
	for _, bytesData := range videoBytes {
		finalBytes = append(finalBytes, bytesData...)
	}

	if len(finalBytes) <= 0 {
		println("video not found")
		return
	}
	saveFile(finalBytes)
}
