package common

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
)

//读取响应
func ParseResponseBody(resp *http.Response) (e error, content string) {
	var body []byte
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("ParseResponseBody gzip.NewReader error:", err)
			e = err
			return
		}
		defer reader.Close()
		body, e = ioutil.ReadAll(reader)
		if e != nil {
			fmt.Println("ParseResponseBody gzip ioutil.ReadAll error:", e)
			return
		}
	default:
		body, e = ioutil.ReadAll(resp.Body)
		if e != nil {
			fmt.Println("ParseResponseBody ioutil.ReadAll error:", e)
			return
		}

	}
	content = string(body)
	return
}
