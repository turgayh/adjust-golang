package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type result struct {
	url          string
	hashResponse string
}

// GetMD5Hash returns the MD5 hash of the given string.
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// parallelGet makes parallel HTTP GET requests to the given URLs and returns the results.
func parallelGet(urls []string, concurrencyLimit int) []result {

	limitChan := make(chan struct{}, concurrencyLimit)

	resultsChan := make(chan *result)

	defer func() {
		close(limitChan)
		close(resultsChan)
	}()

	for _, url := range urls {

		go func(url string) {

			limitChan <- struct{}{}

			url = fmt.Sprintf("http://%s", url)
			res, _ := http.Get(url)
			defer res.Body.Close()

			bodyBytes, _ := io.ReadAll(res.Body)

			bodyString := string(bodyBytes)
			result := &result{url, GetMD5Hash(bodyString)}

			resultsChan <- result

			<-limitChan

		}(url)
	}

	var results []result

	for {
		result := <-resultsChan
		results = append(results, *result)

		if len(results) == len(urls) {
			break
		}
	}

	return results
}

var urls []string
var limit int

// init parses command line arguments.
func init() {
	flag.IntVar(&limit, "parallel", 10, "number of concurrency limit")
	flag.Parse()
	urls = flag.Args()

	if len(urls) == 0 {
		fmt.Println("No URLs provided.")
		os.Exit(1)
	}
}

func main() {
	results := parallelGet(urls, limit)
	for _, val := range results {
		fmt.Printf("%s %s\n", val.url, val.hashResponse)
	}
}
