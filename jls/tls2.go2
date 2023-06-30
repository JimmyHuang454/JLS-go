package jls

import (
	"log"
	"testing"

	"github.com/jls-go/tls"
	"github.com/stretchr/testify/assert"
)

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
	assert.Nil(t, err)
	err = conn.Close()
	assert.Nil(t, err)
}
