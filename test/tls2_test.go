package jls

import (
	"log"
	"net"
	"testing"

	"github.com/JimmyHuang454/JLS-go/tls"
	"github.com/stretchr/testify/assert"
)

var certPem = []byte(`-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE3MTAyMDE5NDMwNloXDTE4MTAyMDE5NDMwNlow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d
7VNhbWvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B
5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1
NDUzgg4xMjcuMC4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2zpJEPQyz6/l
Wf86aX6PepsntZv2GYlA5UpabfT2EZICICpJ5h/iI+i341gBmLiAFQOyTDT+/wQc
6MF9+Yw1Yy0t
-----END CERTIFICATE-----`)

var keyPem = []byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIrYSSNQFaA2Hwf1duRSxKtLYX5CB04fSeQ6tF1aY/PuoAoGCCqGSM49
AwEHoUQDQgAEPR3tU2Fta9ktY+6P9G0cWO+0kETA6SFs38GecTyudlHz6xvCdz8q
EKTcWGekdmdDPsHloRNtsiCa697B2O9IFA==
-----END EC PRIVATE KEY-----`)

func HandleClient(listener net.Listener) {
	for true {
		inClient, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		buf := make([]byte, 200)
		_, err = inClient.Read(buf)
		inClient.Close()
		return
	}
}

func TestDial(t *testing.T) {
	conf := &tls.Config{
		InsecureSkipVerify: false,
	}

	conn, err := tls.Dial("tcp", "uif03.top:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("GET http://uif03.top HTTP/1.1\r\nHost: uif03.top\r\n\r\n"))
	assert.Nil(t, err)
	log.Println(n)

	buf := make([]byte, 200)
	n, err = conn.Read(buf)
	log.Println(buf)
	assert.Nil(t, err)
	err = conn.Close()
	assert.Nil(t, err)
}

func TestWithSelfSignCert(t *testing.T) {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	assert.Nil(t, err)

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	port := "2000"
	listener, err := tls.Listen("tcp", ":"+port, cfg)
	assert.Nil(t, err)

	go HandleClient(listener)

	var config = &tls.Config{InsecureSkipVerify: true, ServerName: "abc.com"}
	assert.Equal(t, config.UseJLS, false)

	conn, err := tls.Dial("tcp", "127.0.0.1:"+port,
		config)
	assert.Nil(t, err)

	err = conn.Close()
	assert.Nil(t, err)
}

func TestJLS(t *testing.T) {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	assert.Nil(t, err)

	cfg := &tls.Config{Certificates: []tls.Certificate{cert},
		UseJLS: true, JLSPWD: []byte("abc"), JLSIV: []byte("abc")}
	port := "2001"
	listener, err := tls.Listen("tcp", ":"+port, cfg)
	assert.Nil(t, err)

	go HandleClient(listener)

	conn, err := tls.Dial("tcp", "127.0.0.1:"+port,
		&tls.Config{InsecureSkipVerify: false,
			ServerName: "abc.com",
			UseJLS:     true, JLSPWD: []byte("abc"), JLSIV: []byte("abc")})
	assert.Nil(t, err)
	err = conn.Close()
	assert.Nil(t, err)

	return
	conn, err = tls.Dial("tcp", "127.0.0.1:"+port,
		&tls.Config{InsecureSkipVerify: false,
			ServerName: "abc.com",
			UseJLS:     true, JLSPWD: []byte("abc"), JLSIV: []byte("abcd")})
	assert.NotNil(t, err)
	err = conn.Close()
	assert.Nil(t, err)
}
