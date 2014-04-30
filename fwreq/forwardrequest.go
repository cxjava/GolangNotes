package fwreq

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/cxjava/GolangNotes/common"
)

var (
	VulnerableCipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	}
)

func ForWardRequest() {
	fmt.Println("ForWardRequest!")
	fmt.Println(doForWardRequest("113.57.187.29", "GET", "http://kyfw.12306.cn/otn/", nil))
	fmt.Println(doForWardRequest("113.57.187.29", "GET", "http://kyfw.12306.cn/otn/leftTicket/init", nil))
	// fmt.Println(doForWardRequest("113.57.187.29", "GET", "https://kyfw.12306.cn/otn/leftTicket/init", nil))
	fmt.Println(doForWardRequest("118.194.41.18", "GET", "http://www.zonezu.com/login.do?from=/daybook.jsp", nil))
	// fmt.Println(doForWardRequest3())
	// fmt.Println(doForWardRequest4())
}

//转发
func doForWardRequest(forwardAddress, method, requestUrl string, body io.Reader) (content string) {
	if !strings.Contains(forwardAddress, ":") {
		forwardAddress = forwardAddress + ":80"
	}

	// conn, err := net.Dial("tcp", forwardAddress)
	conn, err := tls.Dial("tcp", forwardAddress, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println("doForWardRequest DialTimeout error:", err)
		return
	}
	defer conn.Close()

	cs := conn.ConnectionState()
	fmt.Printf("State = %#v\n", cs)

	for i, cert := range cs.PeerCertificates {
		fmt.Printf("Cert[%d] = %x\n", i, cert.Signature)
	}

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
	// httputil.NewServerConn(c, r)
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
func doForWardRequest4() (html string) {
	log.Println("hijacking TLS connection")

	tlsConn, err := tls.Dial("tcp", "113.57.187.29:80", nil)
	defer tlsConn.Close()

	if err != nil {
		log.Println("error dialing TLS, falling back:", err)
		// p.doConnectRequest(w, r)
		return
	}

	cs := tlsConn.ConnectionState()
	peerCerts := cs.PeerCertificates

	fakedCert := tls.Certificate{}

	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println("GenerateKey", err)
		return ""
	}

	fakedCert.PrivateKey = rsaKey

	for _, peerCert := range peerCerts {
		fakedCert.Certificate = append(fakedCert.Certificate, peerCert.Raw)
	}

	host, _, _ := net.SplitHostPort("kyfw.12306.cn")

	config := &tls.Config{
		Certificates:             []tls.Certificate{fakedCert},
		ServerName:               host,
		PreferServerCipherSuites: true,
		CipherSuites:             VulnerableCipherSuites,
		MaxVersion:               tls.VersionTLS11,
	}

	var client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}
	request, err := http.NewRequest("GET", "https://113.57.187.29/otn/leftTicket/init", nil)
	// request, err := http.NewRequest("GET", "http://118.194.41.18/login.do?from=/daybook.jsp", nil)
	if err != nil {
		fmt.Println("http.NewRequest", err)
		html = ""
		return
	}
	request.Close = true
	// request.Header.Set("Host", "www.zonezu.com")
	// request.Header.Set("Host", "kyfw.12306.cn")
	common.AddReqestHeader(request, "GET")
	request.Close = true

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("client.Do", err)
		html = ""
		return
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		_, html = common.ParseResponseBody(response)
		fmt.Println("postUrl:", "response body:", html)
	} else {
		fmt.Println("postUrl:", "Status Code:", response.StatusCode)
		html = ""
	}
	return
}
func doForWardRequest3() (html string) {
	conn, err := tls.Dial("tcp", "218.75.201.31:443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	cs := conn.ConnectionState()
	fmt.Printf("State = %#v\n", cs)

	for i, cert := range cs.PeerCertificates {
		fmt.Printf("Cert[%d] = %x\n", i, cert.Signature)
	}
	fakedCert := tls.Certificate{}
	for _, peerCert := range cs.PeerCertificates {
		fakedCert.Certificate = append(fakedCert.Certificate, peerCert.Raw)
	}

	// client := httputil.NewClientConn(conn, nil)
	// client := &http.Client{}
	// var client = http.Client{
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// 	},
	// }

	config := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{fakedCert},
	}

	var client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}

	request, err := http.NewRequest("GET", "https://218.75.201.31/otn/leftTicket/init", nil)
	// request, err := http.NewRequest("GET", "http://118.194.41.18/login.do?from=/daybook.jsp", nil)
	if err != nil {
		fmt.Println("http.NewRequest", err)
		html = ""
		return
	}
	request.Close = true
	// request.Header.Set("Host", "www.zonezu.com")
	// request.Header.Set("Host", "kyfw.12306.cn")
	common.AddReqestHeader(request, "GET")
	request.Close = true

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("client.Do", err)
		html = ""
		return
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		_, html = common.ParseResponseBody(response)
		fmt.Println("postUrl:", "response body:", html)
	} else {
		fmt.Println("postUrl:", "Status Code:", response.StatusCode)
		html = ""
	}
	return
}
