package icap

import (
	"io"
	"io/ioutil"
	"net"
	"testing"
)

// REQMOD example 2 from RFC 3507, adjusted for order of headers, etc.
func TestREQMOD2(t *testing.T) {
	request :=
		"REQMOD icap://icap-server.net/server?arg=87 ICAP/1.0\r\n" +
			"Host: icap-server.net\r\n" +
			"Encapsulated: req-hdr=0, req-body=154\r\n" +
			"\r\n" +
			"POST /origin-resource/form.pl HTTP/1.1\r\n" +
			"Host: www.origin-server.com\r\n" +
			"Accept: text/html, text/plain\r\n" +
			"Accept-Encoding: compress\r\n" +
			"Cache-Control: no-cache\r\n" +
			"\r\n" +
			"1e\r\n" +
			"I am posting this information.\r\n" +
			"0\r\n" +
			"\r\n"
	resp :=
		"ICAP/1.0 200 OK\r\n" +
			"Date: Mon, 10 Jan 2000  09:55:21 GMT\r\n" +
			"Encapsulated: req-hdr=0, req-body=231\r\n" +
			"Istag: \"W3E4R7U9-L2E4-2\"\r\n" +
			"Server: ICAP-Server-Software/1.0\r\n" +
			"\r\n" +
			"POST /origin-resource/form.pl HTTP/1.1\r\n" +
			"Accept: text/html, text/plain, image/gif\r\n" +
			"Accept-Encoding: gzip, compress\r\n" +
			"Cache-Control: no-cache\r\n" +
			"Host: www.origin-server.com\r\n" +
			"Via: 1.0 icap-server.net (ICAP Example ReqMod Service 1.1)\r\n" +
			"\r\n" +
			"2d\r\n" +
			"I am posting this information.  ICAP powered!\r\n" +
			"0\r\n" +
			"\r\n"

	p1, p2 := net.Pipe()
	c, err := newConn(p2, HandlerFunc(HandleREQMOD2))
	go c.serve()

	io.WriteString(p1, request)
	respBuffer := make([]byte, len(resp))
	_, err = io.ReadFull(p1, respBuffer)

	if err != nil {
		t.Fatalf("error while reading response: %v", err)
	}

	response := string(respBuffer)
	checkString("Response", response, resp, t)
}

func HandleREQMOD2(w ResponseWriter, req *Request) {
	w.Header().Set("Date", "Mon, 10 Jan 2000  09:55:21 GMT")
	w.Header().Set("Server", "ICAP-Server-Software/1.0")
	w.Header().Set("ISTag", "\"W3E4R7U9-L2E4-2\"")

	req.Request.Header.Set("Via", "1.0 icap-server.net (ICAP Example ReqMod Service 1.1)")
	req.Request.Header.Set("Accept", "text/html, text/plain, image/gif")
	req.Request.Header.Set("Accept-Encoding", "gzip, compress")

	body, _ := ioutil.ReadAll(req.Request.Body)
	newBody := string(body) + "  ICAP powered!"

	w.WriteHeader(200, req.Request, true)
	io.WriteString(w, newBody)
}