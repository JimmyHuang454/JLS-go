package jls

import (
	"context"
	"log"
	"net"
	"net/http"
	"testing"

	"github.com/JimmyHuang454/JLS-go/tls"
	"github.com/stretchr/testify/assert"
)

func TestFallback(t *testing.T) {
	serverName := "uif03.top"

	// server
	cert, err := tls.X509KeyPair(certPem, keyPem)
	assert.Nil(t, err)
	serverConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   serverName,
		UseJLS:       true, JLSPWD: []byte("abc"), JLSIV: []byte("abc")}
	port := "2000"
	listener, err := tls.Listen("tcp", ":"+port, serverConfig)
	assert.Nil(t, err)

	go func() {
		inClient, err := listener.Accept()
		defer inClient.Close()
		assert.NotNil(t, err)
		// buf := make([]byte, 200)
		// inClient.Read(buf)
	}()

	// client
	config := &tls.Config{
		ServerName: serverName,
	}
	tcpAddress := "127.0.0.1:" + port
	tcp, err := net.Dial("tcp", tcpAddress)
	assert.Nil(t, err)
	tlsDial := tls.Client(tcp, config)
	err = tlsDial.Handshake()
	assert.Nil(t, err)

	client := &http.Client{
		Transport: &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return tlsDial, nil
			},
		},
	}
	request, _ := http.NewRequest("GET", "https://"+serverName, nil)
	response, err := client.Do(request)
	assert.Nil(t, err)
	defer response.Body.Close()
	buf := make([]byte, 200)
	response.Body.Read(buf)
	log.Println(string(buf))
}
