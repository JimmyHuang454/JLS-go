package tls

import (
	"errors"
	"fmt"
	"io"
	"net"

	r "github.com/JimmyHuang454/JLS-go/jls"
)

func JLSHandler(c *Conn, tlsError error) error {
	if !c.config.UseJLS {
		return tlsError
	}

	if c.isClient {
		if tlsError == nil && !c.IsValidJLS {
			// it is a valid TLS Client but Not JLS,
			defer c.Close()
			return errors.New("not JLS")
			// so we must TODO: act like a normal http request at here
		}
	} else if tlsError != nil && !c.IsValidJLS && c.quic == nil {
		// It is not JLS. Forward at here.
		// TODO: if we using sing-box, we need to use its forward method, since it may take over traffic by Tun.
		defer c.conn.Close()
		if c.config.ServerName != "" {
			server, forwardError := net.Dial("tcp", c.config.ServerName+":443")
			fmt.Println(c.config.ServerName + ":443 forwarding...")
			if forwardError == nil {
				defer server.Close()
				server.Write(c.ClientHelloRecord)
				server.Write(c.ForwardClientHello)
				c.ClientHelloRecord = nil // improve memory.
				c.ForwardClientHello = nil
				go io.Copy(server, c.conn)
				io.Copy(c.conn, server) // block until forward finish.
			}
		}
	}
	return tlsError
}

func BuildJLSClientHello(c *Conn, hello *clientHelloMsg) {
	if !c.config.UseJLS || c.IsBuildedFakeRandom {
		return
	}
	zeroArray := BuildZeroArray()
	hello.random = zeroArray
	withoutBinder, _ := hello.marshalWithoutBinders()
	hello.random, _ = BuildFakeRandom(c.config, withoutBinder)
	copy(hello.raw[6:], hello.random)
	c.IsBuildedFakeRandom = true
}

func BuildJLSServerHello(c *Conn, hello *serverHelloMsg) {
	if !c.config.UseJLS {
		return
	}

	hello.random = BuildZeroArray()
	hello.marshal()

	hello.random, _ = BuildFakeRandom(c.config, hello.raw)
	copy(hello.raw[6:], hello.random)
}

func CheckJLSServerHello(c *Conn, serverHello *serverHelloMsg) {
	c.IsValidJLS = false
	if !c.config.UseJLS {
		return
	}
	serverHello.marshal() // init
	zeroArray := BuildZeroArray()
	raw := make([]byte, len(serverHello.raw))
	copy(raw, serverHello.raw)
	copy(raw[6:], zeroArray)

	c.IsValidJLS, _ = CheckFakeRandom(c.config, raw, serverHello.random)
	c.config.InsecureSkipVerify = c.IsValidJLS
}

// return false means need to forward.
func CheckJLSClientHello(c *Conn, clientHello *clientHelloMsg) (bool, error) {
	c.IsValidJLS = false
	if !c.config.UseJLS {
		return true, errors.New("disable JLS.") // == TLS.
	}
	zeroArray := BuildZeroArray()
	withoutBinder, err := clientHello.marshalWithoutBinders()
	c.ForwardClientHello = clientHello.raw
	if err != nil {
		return false, errors.New("failed to get clientHello raw bytes.")
	}
	raw := make([]byte, len(withoutBinder))
	copy(raw, withoutBinder)
	copy(raw[6:], zeroArray)

	c.IsValidJLS, err = CheckFakeRandom(c.config, raw, clientHello.random)
	if err != nil {
		return false, errors.New("failed to check fakeRandom.")
	}
	if !c.IsValidJLS || c.vers != VersionTLS13 {
		return false, errors.New("wrong fakeRandom.")
	}
	if len(clientHello.keyShares) == 0 {
		fmt.Println("JLS missing keyShare can be not safty.")
	}
	return true, nil // valid JLS.
}

func BuildZeroArray() []byte {
	const byteLen = 32
	zeroArray := make([]byte, byteLen)
	for i := 0; i < byteLen; i++ {
		zeroArray[i] = 0
	}
	return zeroArray
}

func BuildFakeRandom(config *Config, AuthData []byte) ([]byte, error) {
	iv := append(config.JLSIV, AuthData...)
	pwd := append(config.JLSPWD, AuthData...)
	fakeRandom := r.NewFakeRandom(pwd, iv)

	err := fakeRandom.Build()
	return fakeRandom.Random, err
}

func CheckFakeRandom(config *Config, AuthData []byte, random []byte) (bool, error) {
	iv := append(config.JLSIV, AuthData...)
	pwd := append(config.JLSPWD, AuthData...)
	fakeRandom := r.NewFakeRandom(pwd, iv)

	IsValid, err := fakeRandom.Check(random)
	return IsValid, err
}
