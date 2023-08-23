package jls

import (
	"log"
	"testing"

	"github.com/JimmyHuang454/JLS-go/tls"
	"github.com/stretchr/testify/assert"
)

func TestPSK(t *testing.T) {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	assert.Nil(t, err)

	serverConfig := &tls.Config{Certificates: []tls.Certificate{cert},
		UseJLS: true, JLSPWD: []byte("abc"), JLSIV: []byte("abc"), SessionTicketsDisabled: false}
	port := "2021"
	listener, err := tls.Listen("tcp", ":"+port, serverConfig)
	assert.Nil(t, err)

	go func() {
		for {
			inClient, err := listener.Accept()
			log.Println(err)
			assert.Nil(t, err)
			buf := make([]byte, 200)
			n, err := inClient.Read(buf)
			assert.Equal(t, n, 1)
			inClient.Close()
		}
	}()

	for i := 0; i < 3; i++ {
		conn, err := tls.Dial("tcp", "127.0.0.1:"+port,
			&tls.Config{InsecureSkipVerify: false,
				ServerName: "abc.com",
				UseJLS:     true, JLSPWD: []byte("abc"), JLSIV: []byte("abc"), SessionTicketsDisabled: false})
		assert.Nil(t, err)
		conn.Write([]byte{1})
		err = conn.Close()
		assert.Nil(t, err)
	}
}
