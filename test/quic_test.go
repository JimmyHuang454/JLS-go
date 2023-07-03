package jls

import (
	"log"
	"net"
	"testing"

	"github.com/quic-go/quic-go"
	"github.com/stretchr/testify/assert"
)

func TestQuic(t *testing.T) {
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: 1234})
	// ... error handling
	tr := quic.Transport{
		Conn: udpConn,
	}
	ln, err := tr.Listen(tlsConf, quicConf)
	// ... error handling
	go func() {
		for {
			conn, err := ln.Accept()
			// ... error handling
			// handle the connection, usually in a new Go routine
		}
	}()

}
