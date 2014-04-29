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

//添加头
func AddReqestHeader(request *http.Request, method string) {
	request.Header.Set("Host", "kyfw.12306.cn")
	request.Header.Set("Connection", "keep-alive")
	// request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	// request.Header.Set("X-Requested-With", "XMLHttpRequest")
	// request.Header.Set("If-Modified-Since", "0")
	// request.Header.Set("Content-Length", fmt.Sprintf("%d", request.ContentLength))
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.131 Safari/537.36")
	request.Header.Set("DNT", "1")

	if method == "POST" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

	request.Header.Set("Referer", "https://kyfw.12306.cn/otn/")
	request.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.6,en;q=0.4")
	request.Header.Set("Cookie", "JSESSIONID=9AB1BD89D055850AADE7F7D29BFE714B; _jc_save_fromStation=%u5B9C%u660C%u4E1C%2CHAN; _jc_save_toStation=%u6B66%u6C49%2CWHN; _jc_save_fromDate=2014-05-03; _jc_save_toDate=2014-04-14; _jc_save_wfdc_flag=dc; BIGipServerotn=1457062154.50210.0000")
}
