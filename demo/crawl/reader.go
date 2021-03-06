package main

import (
	"os"
	"fmt"
	"io"
	"bufio"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"errors"
)

type ReaderNode struct {
	word  string
	count int
}

func readConfig() ([]string, error){
	inputFile, err := os.Open("urls.txt")
	if err != nil {
		return nil, errors.New("read fail")
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	urls :=[]string{}
	index := 0
	for {
		inputString, readerError := inputReader.ReadString('\n')
		fmt.Printf("The input was: %s", inputString)

		if readerError == io.EOF {
			return urls, nil
		}

		readLine := strings.Trim(inputString, "\n")
		if len(readLine) > 5 {
			urls = append(urls, readLine)
			index++
		}
	}
}

func parseReader() {
	out := make(chan []byte)
	urls, err := readConfig()
	if err != nil {
		return
	}
	println(len(urls))

	for _, url := range urls {
		println(url)

		go func(url string) {
			resp, err := http.Get(url)
			if err != nil {
				out <- []byte{}
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			out <- body
		}(url)
	}

	word := make(chan string)
	themap := make(map[string]int)
	finishMark := make(chan bool)
	var myExp = regexp.MustCompile(`[a-zA-Z]+`)

	go func() {
		for _, _ = range urls {
			body := <-out

			// bodystring := string(body)
			dicts := myExp.FindAll(body, -1)
			for _, a := range dicts {
				key := string(a)
				word <- key
			}
			// println(bodystring)
		}

		finishMark <- true

		// func containString(list []string, word string){

		// }
	}()

	go func(){
		for {
			key := <-word

			if _, ok := themap[key]; ok {
				themap[key]++
			} else {
				themap[key] = 1
			}
			// println(key)
			// println(themap[key])
		}
	}()

	for i:=0 ; i < 2; i++{
		_ =  <-finishMark

		keys := make([]string, len(themap))
		vals := make([]int, len(themap))
		nodeList := make([]ReaderNode, len(themap))

		i := 0
		for k, v := range themap {
			keys[i] = k
			vals[i] = v

			nodeList[i] = ReaderNode{word: k, count: v}
			// append(nodeList, &)
			i++
		}

		// 排序1 冒泡
		for index, _ := range vals {
			for iindex := 0; iindex < len(vals)-1-index; iindex++ {
				if vals[iindex] > vals[iindex+1] {
					tempkey := keys[iindex]
					tempval := vals[iindex]

					keys[iindex] = keys[iindex+1]
					vals[iindex] = vals[iindex+1]

					keys[iindex+1] = tempkey
					vals[iindex+1] = tempval
				}
			}
		}

		for index, count := range vals {
			println(count, keys[index])
		}

		// 排序2 sort模块
		sort.Slice(nodeList, func(i, j int) bool {
			return nodeList[i].count < nodeList[j].count
		})

		for _, nodeItem := range nodeList {
			println(nodeItem.count, nodeItem.word)
		}
	}
}
