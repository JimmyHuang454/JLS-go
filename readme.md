# JLS-go
本仓库完整实现 JLS 第 3 版本。

## 用法
跟`crypto/tls`标准库一样，但多了一些选项。

### Client
```go
import (
	"github.com/JimmyHuang454/JLS-go/tls"
)

conn, err := tls.Dial("tcp", "127.0.0.1:443",
    &tls.Config{InsecureSkipVerify: false,
        ServerName: "abc.com", // 伪装域名
        UseJLS: true, JLSPWD: []byte("密码"), JLSIV: []byte("随机数")})
if err != nil{
    return;
}
defer conn.Close()
n, _ := conn.Read(buffer)
```

### Server
```go
import (
	"github.com/JimmyHuang454/JLS-go/tls"
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

// 证书可以自己随便自签一个
cert, err := tls.X509KeyPair(certPem, keyPem)
assert.Nil(t, err)
cfg := &tls.Config{Certificates: []tls.Certificate{cert},
    ServerName: "abc.com", // 伪装站
    UseJLS: true, JLSPWD: []byte("密码"), JLSIV: []byte("随机数")}

listener, err := tls.Listen("tcp", ":443", cfg)
assert.Nil(t, err)

for true {
    inClient, err := listener.Accept()
    if err != nil {
        log.Println(err)
        return
    }
    buf := make([]byte, 200)
    _, err = inClient.Read(buf)
    inClient.Close()
}
```
