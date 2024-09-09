package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const (
	pageSize  = 100
	namespace = "library"
	urlFormat = "https://registry.hub.docker.com/v2/repositories/%s/%s/tags"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		outErr("not set image name.")
		os.Exit(1)
	}
	image := args[1]

	s := fmt.Sprintf(urlFormat, namespace, image)
	url, err := url.Parse(s)
	if err != nil {
		outErr("invalid url: %s\n%s", s, err.Error())
		os.Exit(1)
	}
	q := url.Query()
	q.Set("page_size", fmt.Sprintf("%d", pageSize))
	url.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Get(url.String())
	if err != nil {
		outErr("cannot get from api: %s", err.Error())
		os.Exit(1)
	}
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	resJSON := buf.Bytes()
	if resp.StatusCode != http.StatusOK {
		outErr("invalid status: %d : %s", resp.StatusCode, string(resJSON))
		os.Exit(1)
	}
	var res Response
	if err := json.Unmarshal(resJSON, &res); err != nil {
		outErr("parse error: %s", err.Error())
		os.Exit(1)
	}
	for _, r := range res.Results {
		fmt.Println(r.Name)
	}
}

func outErr(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
}
