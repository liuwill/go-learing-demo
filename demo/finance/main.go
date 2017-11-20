package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var targetUrl = flag.String("url", "", "jd finance url")
var targetPin = flag.String("pin", "", "pin id")

func ExampleScrape() {
	println(*targetPin)
	println(*targetUrl)
	resp, err := http.Post(*targetUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader("pin="+(*targetPin)))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		os.Exit(1)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Find the review items
	fields := []string{"author", "action", "meta", "deal", "duration", "price"}
	colList := []string{"liuwill"}
	doc.Find("table tbody").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		s.Find("tr").Each(func(ii int, trs *goquery.Selection) {
			trs.Find("td").Each(func(iii int, tds *goquery.Selection) {
				targetStr := strings.TrimSpace(tds.Text())
				targetStr = strings.Replace(targetStr, "\n", "", -1)
				targetStr = strings.Replace(targetStr, "\t", "", -1)
				colList = append(colList, targetStr)
				fmt.Println(strings.TrimSpace(tds.Text()))
			})
		})
	})

	resultList := []string{}
	for index, value := range fields {
		resultList = append(resultList, value+":"+colList[index])
	}

	jsonStr, err := json.Marshal(map[string]interface{}{
		"author": "liuwill",
		"doc":    strings.Join(resultList, "\n"),
	})

	println(string(jsonStr))
	req, err := http.NewRequest("POST", "http://192.168.2.239:7077/tmp/wechat", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	respF, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", respF.Status)
	fmt.Println("response Headers:", respF.Header)
	body, _ := ioutil.ReadAll(respF.Body)
	fmt.Println("response Body:", string(body))
}

func main() {
	flag.Parse()
	ExampleScrape()
}
