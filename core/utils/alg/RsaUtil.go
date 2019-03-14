package alg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// 配置文件使用rsa加密，网络传输使用密文传输，具体执行加密解密的是
// goos客户端和goos的react中端

var privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAntcofxD+O7wJg722up6ggThTClKJpnl9+nRMptCP+GYmqNpG
a6e9hEmxO+BL5//+FF/JDIvJgm8KoFjiosZV+6p2maDWp82gIuF8DZiqLRRRLb1H
+JOqelchRXK3sii7usaS+OmHdoJUVDH/Ddp78ESjchDg6ollb8ijJ1JMD0lmvcDF
LYHkfLEtm35V/yn2B8LxoWyQ+V0qZ25J3U/3NKVGxXgvqPGqc5nClcPDJB6nLAJV
AcVVPh7bgS15dLtUWicimFcF+ibPYR9GHDn8Ts4/WldQiMeMPXJwWusyCO2Qh4px
18fmVYVbNEYDxx8Neizi15W73cAa+cAuKUEbRwIDAQABAoIBAFRbJemdp7cnnNH3
TfT8S3d05z0DKpFb0ljRrOemud8MuRlRmraPxelcjaCj9QwH+vLtD8P29RStTjJs
IiiaNo3KMORT88v4O0RrBcveuPnp4VbuQCu2mIIShdvxGbenRFPkI6fHtiZs/sYc
Kz07PDkU6syoRBqqz1E3d9ZpMXWoCmGDlajlO+cYBk/Ifv7LnhmKZ/8LwraLPuld
qQ927YdB2qPLCNncCq2UKNHhfNgs/ZkQG0kfBD5fAJzK200Vi6x7gR4k5aFcrRWt
cIwuLoQ/7w5wu7oUB2DsY2G+265jY0QlZibqPRTUxcHn6qWtELY+ChUjd7DLbcTx
+0HcfgkCgYEA0qIJY38Aa4gjQbN0e9QVOqEHDkAxayPhwTq42bGzahLuFU50Stcq
bAg/NHmFN/mZYnl3jjfr7AyJf8djBoBGIt1QsE4pX68SaIhgcV1T2WStPWZ/0Leq
xhEHIlTcniCFmzvzwT3gKawbAB4RGmOFTBgKEAkKn4dMILp9mLVk6xsCgYEAwQ1d
0LLRFHrekaqDWwA/LDU1DCs1WoZY2jpZOrvCld5PI5hwnzmGbN7iUrG3OxeARpof
8jKUIudnXa8KlDN+WFSh6wdt3rxVdK3njjRC9kl40AHr24Bph8HQTDS3jBC4Q/h6
Am0b2Vqp6IfmQBpwnVNYrwWqxN3xiPzsvxHgB0UCgYADnnEW2onBTzd723922S/8
L+QVJJk0xAsO7NlcNCdm/ShGCXEKECRUctfTKPtW+NzfykJ5mRBen+CE54IWDIIn
+zF1tgIT+MgSL7WofPgB1i77zRUJGv6+JHDA8EBSHzDsd906RrvhO2nDWMDmeStD
IW6a0+zwzLxMG6goxYUUXQKBgDhQ2+M7pr5gsXiw0yzCv6r3wofQOvozYswWZV06
1KK/fTqXB5OLRmmQA1m4Oglk66is4VDX7FraQk8T9vQQqXS/C5TyT9y/9/XXnUrg
eAA8op+bT+Byb1aI9WiloD2dywMZAw6eIZegWRxaOJLOo7dhTuePsadIT2N01ONQ
JZg5AoGBAJAjMSj7wjhqegbYSVGVmNsVXITcvr+IRJQyw205AQTubcXSmW2q9yHo
t1NHs9tJoRy24cSSM3Zn4n9izeZSc1TKvz0Qt3M10BL77P4CHEJ8MA/DvYzcxD3j
P3uTnZMD3mM+Gzc/zLYG2avffVZiKuGZnxdOS3ZLiAGTS5P7QDKQ
-----END RSA PRIVATE KEY-----`)

var publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAntcofxD+O7wJg722up6g
gThTClKJpnl9+nRMptCP+GYmqNpGa6e9hEmxO+BL5//+FF/JDIvJgm8KoFjiosZV
+6p2maDWp82gIuF8DZiqLRRRLb1H+JOqelchRXK3sii7usaS+OmHdoJUVDH/Ddp7
8ESjchDg6ollb8ijJ1JMD0lmvcDFLYHkfLEtm35V/yn2B8LxoWyQ+V0qZ25J3U/3
NKVGxXgvqPGqc5nClcPDJB6nLAJVAcVVPh7bgS15dLtUWicimFcF+ibPYR9GHDn8
Ts4/WldQiMeMPXJwWusyCO2Qh4px18fmVYVbNEYDxx8Neizi15W73cAa+cAuKUEb
RwIDAQAB
-----END PUBLIC KEY-----`)

// 加密
func RsaEncrypt(data string) (string, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)

	if strBytes, e := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(data)); e != nil {
		return "", e
	} else {
		return base64.StdEncoding.EncodeToString(strBytes), nil
	}
}

// 解密
func RsaDecrypt(data string) (string, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	if encryptedDecodeBytes, err:=base64.StdEncoding.DecodeString(data); err != nil {
		return "", err
	} else {
		if strBytes, e := rsa.DecryptPKCS1v15(rand.Reader, priv, encryptedDecodeBytes); e != nil {
			return "", e
		} else {
			return string(strBytes), nil
		}
	}
}