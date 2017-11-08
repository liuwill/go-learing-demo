package main

import (
	"strings"
	"net/http"
	"os"
	_ "errors"
	_ "io"
	_ "bufio"
	"io/ioutil"
	_ "regexp"
	"flag"
	"github.com/tidwall/gjson"
	"./logger"
)

var name = flag.String("name", "movie.ts", "save file name")
var mode = flag.String("mode", "simple", "input m3u8 url")
var info = flag.String("info", "", "url for simple mode")

type DownResult struct {
	content []byte
	pos int
}

var out = make(chan []byte)
var downChan = make(chan DownResult)
var batchChan = make(chan int)
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

func getFileName() string {
	outputFile := "movie.ts"
	if len(*name) > 0 {
		outputFile = *name
	}

	return outputFile
}

func saveFile(bytesData []byte){
	outputFile := getFileName()

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

func fetchM3u8File(m3u8Url string) []byte {
	go fetchMeta(m3u8Url)
	m3u8Content := <- out

	return m3u8Content
}

func decodeVideoUrl(m3u8Content string) []string{
	contentList := strings.Split(string(m3u8Content), "\r\n")
	videoUrls := []string{}

	for _, contentLine := range contentList {
		if strings.Contains(contentLine, "http") {
			videoUrls = append(videoUrls, strings.Trim(contentLine, " "))
		}
	}
	return videoUrls
}

func saveVideoSegment(videoUrls []string, progress int, segmentCount int) {
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

func Exist(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil || os.IsExist(err)
}

func main (){
	flag.Parse()

	targetFile := getFileName()
	if Exist(targetFile) {
		logger.Error("File [ " + targetFile + " ] is exist")
		return
	}

	var m3u8Url = ""
	if *mode == "simple" {
		m3u8Url = loadVideoUrlSimple()
	} else {
		go loadVideoUrlFromRemote()
		m3u8Url = <- urlChan
	}

	println(m3u8Url)

	m3u8Content := fetchM3u8File(m3u8Url)
	videoUrls := decodeVideoUrl(string(m3u8Content))

	ctrlCount := 0
	finishMark := false

	for i, vl := range videoUrls {
		finishMark = false
		go downLoadFile(i, vl)

		ctrlCount += 1
		if ctrlCount > 50 {
			finishMark = true

			saveVideoSegment(videoUrls, i, ctrlCount)
			ctrlCount = 0
		}
	}

	if !finishMark {
		saveVideoSegment(videoUrls, len(videoUrls), ctrlCount)
	}
}
