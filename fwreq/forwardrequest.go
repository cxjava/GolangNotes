package fwreq

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/cxjava/GolangNotes/common"
)

func ForWardRequest() {
	fmt.Println("ForWardRequest!")
}

//转发
func doForWardRequest(forwardAddress, method, requestUrl string, body io.Reader) (content string) {
	if !strings.Contains(forwardAddress, ":") {
		forwardAddress = forwardAddress + ":80"
	}

	conn, err := net.DialTimeout("tcp", forwardAddress, 20)
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
