package fwreq

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/cxjava/GolangNotes/common"
)

func ForWardRequest() {
	fmt.Println("ForWardRequest!")
	// body := doForWardRequest2("113.57.187.29", "GET", "https://kyfw.12306.cn/otn/leftTicket/init", nil)
	// body := doForWardRequest("113.57.187.29", "GET", "http://kyfw.12306.cn/otn/leftTicket/init", nil)
	body := doForWardRequest2("118.194.41.18", "GET", "http://www.zonezu.com/login.do?from=/daybook.jsp", nil)
	fmt.Println(body)
}

//转发
func doForWardRequest(forwardAddress, method, requestUrl string, body io.Reader) (content string) {
	if !strings.Contains(forwardAddress, ":") {
		forwardAddress = forwardAddress + ":80"
	}

	conn, err := net.Dial("tcp", forwardAddress)
	if err != nil {
		fmt.Println("doForWardRequest DialTimeout error:", err)
		return
	}
	defer conn.Close()
	//buf_forward_conn *bufio.Reader
	buf_forward_conn := bufio.NewReader(conn)

	req, err := http.NewRequest(method, requestUrl, body)
	if err != nil {
		fmt.Println("doForWardRequest NewRequest error:", err)
		return
	}
	common.AddReqestHeader(req, method)
	var errWrite error

	errWrite = req.Write(conn)
	if errWrite != nil {
		fmt.Println("doForWardRequest Write error:", errWrite)
		return
	}

	resp, err := http.ReadResponse(buf_forward_conn, req)

	if err != nil {
		fmt.Println("doForWardRequest ReadResponse error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var err error
		err, content = common.ParseResponseBody(resp)
		if err != nil {
			fmt.Println("doForWardRequest ParseResponseBody error:", err)
			return
		}
		fmt.Println("doForWardRequest content:", content)
	} else {
		fmt.Println("StatusCode:", resp.StatusCode, resp.Header, resp.Cookies())
	}
	return
}

//转发
func doForWardRequest2(forwardAddress, method, requestUrl string, body io.Reader) (content string) {
	tcpConn, err := net.Dial("tcp", forwardAddress+":80")
	if err != nil {
		fmt.Println("net.Dial, error", err)
		return
	}
	// cf := &tls.Config{Rand: crand.Reader}
	// ssl := tls.Client(tcpConn, cf)

	reqest, err := http.NewRequest(method, requestUrl, body)
	if err != nil {
		fmt.Println("NewRequest, error", err)
		return
	}
	common.AddReqestHeader(reqest, method)

	fmt.Println(reqest.URL.Host, ",", reqest.URL.Path)

	clientConn := httputil.NewClientConn(tcpConn, nil)

	//req, err := http.NewRequest("GET", c.path.String(), nil)
	resp, err := clientConn.Do(reqest)

	if err != nil {
		fmt.Println("Client.Do:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var err error
		err, content = common.ParseResponseBody(resp)
		if err != nil {
			fmt.Println("doForWardRequest ParseResponseBody error:", err)
			return
		}
		fmt.Println("doForWardRequest content:", content)
	} else {
		fmt.Println("StatusCode:", resp.StatusCode, resp.Header, resp.Cookies())
	}
	return
}
