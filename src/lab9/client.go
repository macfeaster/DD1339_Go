package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	server := []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
	}
	for {
		before := time.Now()
		// res := Get(server[0])
		// res := Read(server[0], time.Second)
		res := MultiRead(server, time.Second)
		after := time.Now()
		fmt.Println("Response:", *res)
		fmt.Println("Time:", after.Sub(before))
		fmt.Println()
		time.Sleep(500 * time.Millisecond)
	}
}

type Response struct {
	Body       string
	StatusCode int
}

// Get makes an HTTP Get request and returns an abbreviated response.
// Status code 200 means that the request was successful.
// The function returns &Response{"", 0} if the request fails
// and it blocks forever if the server doesn't respond.
func Get(url string) *Response {
	res, err := http.Get(url)
	if err != nil {
		return &Response{}
	}
	// res.Body != nil when err == nil
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
	}
	return &Response{string(body), res.StatusCode}
}

// I've found two insidious bugs in this function; both of them are unlikely
// to show up in testing. Please fix them right away and don't forget to
// write a doc comment this time.

// Bug 1: res might be subject to a data race, if it is assigned by the
// 	      go routine and the timeout case at the same time.
// 		  Solved by using the done channel for the go routine result.
// Bug 2: the go routine continuously blocks until Get() is done, and thus
//        might never quit. To solve this, the go routine has to be killed
//        after a certain amount of time. I solve this by using select.
func Read(url string, timeout time.Duration) (res *Response) {
	done := make(chan *Response) // Use a channel for res to avoid data race
	go func() {
		select {
		case done <- Get(url):
		case <-time.After(timeout):
			// We timed out, kill this go routine
			return
		}
	}()
	select {
	case res = <-done:
	case <-time.After(timeout):
		res = &Response{"Gateway timeout\n", 504}
	}
	return
}

// MultiRead makes an HTTP Get request to each url and returns
// the response of the first server to answer with status code 200.
// If none of the servers answer before timeout, the response is
// 503 - Service unavailable.
func MultiRead(urls []string, timeout time.Duration) (res *Response) {
	done := make(chan *Response)

	for _, url := range urls {
		go func(url string) {
			response := Read(url, timeout)
			if response.StatusCode == 200 {
				done <- response
			}
		}(url)
	}

	select {
	case res = <-done:
	case <-time.After(timeout):
		res = &Response{"Service unavailable\n", 503}
	}
	return
}
