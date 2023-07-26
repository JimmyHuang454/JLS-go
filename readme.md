# JLS-go

本仓库完整实现 JLS 第 3 版本。

## 其他实现

[vincentliu77/rustls-jls](https://github.com/vincentliu77/rustls-jls)

[vincentliu77/quinn-jls](https://github.com/vincentliu77/quinn-jls)（推荐，支持 QUIC）

## 用法

跟`crypto/tls`标准库一样，但多了一些选项。

选项中可以自定义是否开启 0-RTT（默认关闭），比如说需要 QUIC 的 0-RTT，那么可以设置 SessionTicketsDisabled 为 false。需要特别注意，0-RTT 具有安全问题，不保证前向安全，不保证不会被重放攻击和特征识别，优点就是加上 QUIC 传输层，延迟很低，而且数据是无法被解密。

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

inClient, err := listener.Accept()
if err == nil {
    buf := make([]byte, 200)
    _, err = inClient.Read(buf)
    inClient.Close()
}
```

更多用法，参考 [测试用例](https://github.com/JimmyHuang454/JLS-go/tree/master/test)

<!-- ### QUIC -->
<!-- JLS 是支持 QUIC 的，因为 JLS 不依赖 SessionID，而 QUIC 对 TLS 中的 SessionID 有要求。以前的 crypto/tls 是不支持 0-RTT，所以 quic-go 是创建了 crypto/tls 的分支 qtls，以实现 quic 的 0-RTT，最近（2023 年）crypto/tls 加入了 0-RTT，但目前 quic-go 还在使用自己的 qtls。 -->
