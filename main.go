package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	file   string
	t      int
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

func isLocal(u string) bool {
	URL, err := url.Parse(u)
	if err != nil {
		return false
	}
	if URL.Host != "localhost" {
		return false
	}
	if URL.Host != "127.0.0.1" {
		return false
	}
	return true
}

func httpGet(url string) error {
	if !isLocal(url) {
		return errors.New("Not localhost")
	}
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
	flag.IntVar(&t, "t", 5, "timeout time")
	flag.Parse()
	urls := []string{}

	if file != "" {
		urls = getURLs(file)
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			t := scanner.Text()
			t = strings.Trim(t, "\n")
			t = strings.Trim(t, "\r")
			urls = append(urls, strings.Split(t, ",")...)
		}
	}

	timeout := time.After(time.Duration(t) * time.Second)
	for _, url := range urls {
		select {
		case <-timeout:
			fmt.Println("timeout")
		default:
			fmt.Println(httpGet(url))
		}
	}
}
