package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"golang.org/x/net/publicsuffix"
)

// ... pack /// unpack ...
type option func(*http.Client)

func NewClient(options ...option) *http.Client {
	jar, _ := cookiejar.New(&cookiejar.Options{
		// whole url mattern
		PublicSuffixList: publicsuffix.List,
	})
	client := &http.Client{
		Timeout: 100 * time.Second,
		Jar:     jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("too many redirects")
			}
			fmt.Printf("→ redirecting to %s\n", req.URL)
			return nil
		},
	}
	for _, opt := range options {
		opt(client)
	}
	return client
}

// INFO: go dosent' allow name funton inside funtion ohh it's so f imp
func withTimeout(d time.Duration) option {
	return func(c *http.Client) {
		c.Timeout = d
	}
}

func call() {
	// now we cna only set the apramete whihc we want and we cna ingor  ah we want no direct arugme all
	client := NewClient(withTimeout(time.Second))
	res, err := client.Head("https://www.shekharx.in/")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	fmt.Println(res.Status)
}
