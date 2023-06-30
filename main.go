package main

import (
	"log"

	"github.com/jls-go/tls"
)

func main() {
	log.SetFlags(log.Lshortfile)

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
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 2000)
	for true {
		n, err = conn.Read(buf)
		if err != nil {
			log.Println(n, err)
			return
		}

		println(string(buf[:n]))
		conn.Close()
	}
}
