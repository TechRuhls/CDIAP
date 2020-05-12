package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var outputdir = "output"
var extension = ".json"
var commandTermFile = "newsapi_commands.list"
var fileDelimiter = "\n"
var guardianKey string
var newsAPIKey string

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func setPaths() {
	path, err := os.Getwd()
	fmt.Println(path)
	pathSep := string(os.PathSeparator)
	if err == nil {
		abPath := path + pathSep
		outputdir = abPath + outputdir + pathSep
	} else {
		check(err)
	}

}
func getAPIURLs() map[string]string {
	var APIUrls = make(map[string]string)
	APIUrls["newsAPI"] = getNewsAPIKey()
	return APIUrls
}
func createCommandURLs(urlsFromFile []string) func(string, string) map[string]string {
	APIUrls := getAPIURLs()
	var urls = make(map[string]string)
	return func(a string, b string) map[string]string {

		for i := 0; i < len(urlsFromFile); i++ {
			safeURL := strings.Split(urlsFromFile[i], "|")[1]
			urls[strings.Split(urlsFromFile[i], "|")[0]] = safeURL
		}
		for x, url := range urls {
			url = fmt.Sprintf(url+"&"+a+"=%s", APIUrls[b])
			urls[x] = url
		}
		return urls
	}
    func readFile(fileName string) string {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func getCommandTerms() []string {
	commands := strings.Split(readFile(commandTermFile), fileDelimiter)
	return commands
}
func getURLDataByCommand(url string, k string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			d1 := []byte(string(body))
			err := ioutil.WriteFile(outputdir+strconv.Itoa(int(getTimeInMilliSeconds()))+"_"+k+extension, d1, 0644)
			check(err)
			defer wg.Done()
		}
	}
}

func getTimeInMilliSeconds() int64 {
	return time.Now().UnixNano() / 1000
}
func main() {
	setPaths()
	commands := getCommands()
	wg.Add(len(commands))
	for k, url := range commands {
		go getURLDataByCommand(url, k)
	}
	wg.Wait()
}
