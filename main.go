package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	file   string
	client = &http.Client{Timeout: time.Duration(10) * time.Second}
)

func getURLs(file string) []string {
	URLs := []string{}
	bin, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	json.Unmarshal(bin, &URLs)
	return URLs
}

func httpGet(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "remlin cache hotter.")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.StringVar(&file, "file", "", "file name")
	flag.Parse()
	urls := []string{}

	if file != "" {
		urls = getURLs(file)
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			//urls = append(urls, scanner.Text())
			t := scanner.Text()
			t = strings.Trim(t, "\n")
			t = strings.Trim(t, "\r")
			urls = append(urls, strings.Split(t, ",")...)
		}
	}

	timeout := time.After(5 * time.Second)
	for _, url := range urls {
		select {
		case <-timeout:
			fmt.Println("time out")
		default:
			fmt.Println(httpGet(url))
		}
	}
}
