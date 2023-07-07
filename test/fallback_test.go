package jls

import (
	"log"
	"net/http"
	"testing"

	"github.com/JimmyHuang454/JLS-go/tls"
	"github.com/stretchr/testify/assert"
)

func TestFallback(t *testing.T) {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	assert.Nil(t, err)

	serverConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "baidu.com",
		UseJLS:       true, JLSPWD: []byte("abc"), JLSIV: []byte("abc")}
	port := "2000"
	listener, err := tls.Listen("tcp", ":"+port, serverConfig)
	assert.Nil(t, err)

	go func() {
		inClient, err := listener.Accept()
		defer inClient.Close()
		assert.NotNil(t, err)
	}()

	client := &http.Client{}
	request, _ := http.NewRequest("GET", "https://baidu.com", nil)
	response, err := client.Do(request)
	defer response.Body.Close()
	assert.Nil(t, err)
	assert.NotEmpty(t, response)
	log.Println(response)
}
