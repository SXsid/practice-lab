package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

type Location struct {
	Path       string
	LineNumber string
}

func (l Location) String() string {
	return fmt.Sprintf("path:%s\nlineNumber:%s\n", l.Path, l.LineNumber)
}

func parser(stacTrace []byte) string {
	var res []Location
	data := bytes.Split(stacTrace, []byte("\t"))
	for _, slice := range data {
		slice = bytes.TrimSpace(slice)
		if len(slice) == 0 {
			continue
		}
		parts := bytes.SplitN(slice, []byte(":"), 2)
		if len(parts) != 2 {
			continue
		}
		LineNumber := bytes.Split(parts[1], []byte(" "))
		if len(LineNumber) == 0 {
			continue
		}
		res = append(res, Location{
			string(parts[0]),
			string(LineNumber[0]),
		})

	}
	var links []string
	for _, lt := range res[1:] {
		v := url.Values{}
		v.Add("path", lt.Path)
		v.Add("linenumber", lt.LineNumber)
		queryString := v.Encode()
		link := fmt.Sprintf("<a href=\"/debug?%s\">%s:%s</a>", queryString, lt.Path, lt.LineNumber)
		links = append(links, link)
	}
	return strings.Join(links, "<br/>")
}
