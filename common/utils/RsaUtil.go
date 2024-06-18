package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

type RsaUtil struct {
}

// 秘钥生成网址: https://tool.ip138.com/rsa/

var publicKeyStr = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCYGLFk0CjXJJOtThv91VLNDymD
T24GteAxWeGAfZNZ4RFCRvYmCfZXvChDVCoBbFwX0XXfBFLByo8WkDYXx/6s/L2q
Qwjm7q8NsGPj5SAQoQDj8/ehwIw9ONPbSxF/LG/2JgigpE3ckoZhNXVhLgGeI04U
zzCwBBfNIgX4UiP97wIDAQAB
-----END PUBLIC KEY-----`

var privateKeyStr = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCYGLFk0CjXJJOtThv91VLNDymDT24GteAxWeGAfZNZ4RFCRvYm
CfZXvChDVCoBbFwX0XXfBFLByo8WkDYXx/6s/L2qQwjm7q8NsGPj5SAQoQDj8/eh
wIw9ONPbSxF/LG/2JgigpE3ckoZhNXVhLgGeI04UzzCwBBfNIgX4UiP97wIDAQAB
AoGBAJTSZJFeVPfepFlJOn5uw2w+T8JacDBEui/P4KSXOx0Q6pBNWwDxcod6ZnMq
4UcvPhVYMNudIVTZ3JSZWzR9zqT3IwZTCEPdZvC6ILErpRifUkm4z7rf72lpk07R
GttFOehclZYNlS9xSpqkNK/K28oeVzK+wCyZ2WCAB6/sU/yhAkEAx2kqgG8lGdNg
3ozuGwkMEMeXaz0uTVEK4Jz9oKKOR9ya18rbDww2VtyubIfxvT01URIAT3Ws6iuD
tbUiD8VeNwJBAMNCOaZjya5tMjFdroMsJrHLSZFKJz/eAblvn+PFEOComr8r0hLi
WEnh+Mu5g3eGClOlikRtdqAvzvh2jgyiwgkCQEOi+ySHDml9Fe1Glfibj/kdCdH4
9XyKEYtwFGLo4COlwuuQxc6L0N6TiaIMVkVevnfaCDrrahQfYFRAtOXuhu8CQGRb
YW4X080G6slcsRlSVAEFzyYRyuKUpKY+rRtQakBN6FthlnOGSoKO1mU/UEbaaexc
JRjOei4S5Hnn1VLBRKECQDytwiYSj/BPyaDS73/jxMDoJhDlfnjrlAHVHzUEJOyD
mwHuUe5cYaD6rLD9YySrMjpgO8CGTknC2s21sVV3USs=
-----END RSA PRIVATE KEY-----`

// LoadPublicKey 从字符串加载 PEM 编码的 RSA 公钥
func LoadPublicKey(pubKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not of type *rsa.PublicKey")
	}
	return publicKey, nil
}

// LoadPrivateKey 从字符串加载 PEM 编码的 RSA 私钥
func LoadPrivateKey(privKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func EncryptCustom(plaintextString string, publicKeyString string) (string, error) {
	publicKey, _ := LoadPublicKey(publicKeyString)
	plaintext := []byte(plaintextString)

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)

	if err != nil {
		fmt.Println(err.Error())
		return "null", err
	}
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return encodedCiphertext, nil
}

// Encrypt 使用公钥对数据进行加密
func Encrypt(plaintextString string) (string, error) {
	publicKey, _ := LoadPublicKey(publicKeyStr)
	plaintext := []byte(plaintextString)

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)

	if err != nil {
		fmt.Println(err.Error())
		return "null", err
	}
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return encodedCiphertext, nil
}

// Decrypt 使用私钥对数据进行解密
func Decrypt(ciphertextString string) (string, error) {

	privateKey, _ := LoadPrivateKey(privateKeyStr)
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextString)
	if err != nil {
		return "null", err
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		fmt.Println(err.Error())
		return "null", err
	}
	return string(plaintext), nil
}
