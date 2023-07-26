package jls

import (
	"log"
	"testing"

	"github.com/JimmyHuang454/JLS-go/tls"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	serverName := "github.com"
	clientConfig := &tls.Config{InsecureSkipVerify: false,
		ServerName: serverName,
		UseJLS:     true, JLSPWD: []byte("3070111071563328618171495819203123318"), JLSIV: []byte("3070111071563328618171495819203123318")}
	address := "127.0.0.1:4443"
	c, err := tls.Dial("tcp", address, clientConfig)
	assert.Nil(t, err)
	c.Write([]byte("GET / HTTP/1.1\r\nHost: github.com\r\n"))

	defer c.Close()
	buf := make([]byte, 200)
	n, err := c.Read(buf)
	log.Println(string(buf))
	log.Println(n)
}

func TestWrongClient(t *testing.T) {
	serverName := "github.com"
	clientConfig := &tls.Config{InsecureSkipVerify: false,
		ServerName: serverName,
		UseJLS:     true, JLSPWD: []byte("1"), JLSIV: []byte("2")}
	address := "127.0.0.1:4443"
	c, err := tls.Dial("tcp", address, clientConfig)
	assert.NotNil(t, err)
	c.Close()
}

func TestTLSClient(t *testing.T) {
	serverName := "github.com"
	clientConfig := &tls.Config{InsecureSkipVerify: false,
		ServerName: serverName,
		UseJLS:     false}
	address := "127.0.0.1:4443"
	c, err := tls.Dial("tcp", address, clientConfig)
	assert.NotNil(t, err)
	c.Close()
}
