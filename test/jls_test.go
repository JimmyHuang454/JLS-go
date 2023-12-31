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

	conn, err := tls.Dial("tcp", "apple.com:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("GET http://apple.com HTTP/1.1\r\nHost: apple.com\r\n\r\n"))
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

	serverConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	port := "2200"
	listener, err := tls.Listen("tcp", ":"+port, serverConfig)
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

func TestRightJLS(t *testing.T) {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	assert.Nil(t, err)

	serverConfig := &tls.Config{Certificates: []tls.Certificate{cert},
		UseJLS: true, JLSPWD: []byte("abc"), JLSIV: []byte("abc")}
	port := "2001"
	listener, err := tls.Listen("tcp", ":"+port, serverConfig)
	assert.Nil(t, err)

	go func() {
		inClient, err := listener.Accept()
		assert.Nil(t, err)
		buf := make([]byte, 200)
		n, err := inClient.Read(buf)
		assert.Equal(t, n, 1)
		inClient.Close()
	}()

	conn, err := tls.Dial("tcp", "127.0.0.1:"+port,
		&tls.Config{InsecureSkipVerify: false,
			ServerName: "abc.com",
			UseJLS:     true, JLSPWD: []byte("abc"), JLSIV: []byte("abc")})
	assert.Nil(t, err)
	conn.Write([]byte{1})
	err = conn.Close()
	assert.Nil(t, err)
}

func TestWrongJLS(t *testing.T) {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	assert.Nil(t, err)

	serverConfig := &tls.Config{Certificates: []tls.Certificate{cert},
		UseJLS: true, JLSPWD: []byte("abc"), JLSIV: []byte("abc")}
	port := "2002"
	listener, err := tls.Listen("tcp", ":"+port, serverConfig)
	assert.Nil(t, err)

	go func() {
		inClient, err := listener.Accept()
		assert.NotNil(t, err)
		inClient.Close()
	}()

	conn, err := tls.Dial("tcp", "127.0.0.1:"+port,
		&tls.Config{InsecureSkipVerify: false,
			ServerName: "abc.com",
			UseJLS:     true, JLSPWD: []byte("abc"), JLSIV: []byte("abcd")})
	assert.NotNil(t, err)
	assert.Nil(t, conn)
}

func TestProvideChannal(t *testing.T) {
	serverName := "apple.com"
	cert, err := tls.X509KeyPair(certPem, keyPem)
	serverConfig := &tls.Config{Certificates: []tls.Certificate{cert},
		ServerName: serverName,
		UseJLS:     true, JLSPWD: []byte("abc"), JLSIV: []byte("abc")}

	port := "2003"
	address := "127.0.0.1:" + port

	listener, err := net.Listen("tcp", address)
	assert.Nil(t, err)

	// ok JLS
	go func() {
		inClient, err := listener.Accept()
		assert.Nil(t, err)
		assert.NotNil(t, inClient)

		safeServer := tls.Server(inClient, serverConfig)
		assert.NotNil(t, safeServer)
		err = safeServer.Handshake()
		assert.Nil(t, err)
		inClient.Close()
	}()
	clientConfig := &tls.Config{InsecureSkipVerify: false,
		ServerName: serverName,
		UseJLS:     true, JLSPWD: []byte("abc"), JLSIV: []byte("abc")}

	c, err := net.Dial("tcp", address)
	safeClient := tls.Client(c, clientConfig)
	err = safeClient.Handshake()
	assert.Nil(t, err)
}

func TestWrongProvideChannal(t *testing.T) {
	serverName := "apple.com"
	cert, err := tls.X509KeyPair(certPem, keyPem)
	serverConfig := &tls.Config{Certificates: []tls.Certificate{cert},
		ServerName: serverName,
		UseJLS:     true, JLSPWD: []byte("1"), JLSIV: []byte("2")}

	port := "2004"
	address := "127.0.0.1:" + port

	listener, err := net.Listen("tcp", address)
	assert.Nil(t, err)

	clientConfig := &tls.Config{InsecureSkipVerify: false,
		ServerName: serverName,
		UseJLS:     true, JLSPWD: []byte("3"), JLSIV: []byte("4")}

	// wrong JLS
	go func() {
		inClient, err := listener.Accept()
		assert.Nil(t, err)
		assert.NotNil(t, inClient)

		safeServer := tls.Server(inClient, serverConfig)
		err = safeServer.Handshake()
		assert.NotNil(t, err)
		inClient.Close()
	}()
	c, err := net.Dial("tcp", address)
	errorClient := tls.Client(c, clientConfig)
	err = errorClient.Handshake()
	assert.NotNil(t, err)
	log.Println(err)
}
