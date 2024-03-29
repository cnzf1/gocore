/*
* @Author: cnzf1
* @Date: 2021-08-05 17:31:50
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-01 16:50:56
* @Description:
*/
package codec_test

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/cnzf1/gocore/codec"
	"github.com/stretchr/testify/assert"
)

const (
	priKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAujka1xH+mye8zIlAzFuQ6UsVLKWiNDsINBTDSPHZ3qQujakL
3LQQjYpiKMH5qE/inVJncZ/rQl1bspxSjncU4GLYu6GBcMwlbDNCpH1DZwXuZQlr
wCXghFnFVqdtnPJ28Oq55oDDjLQrQVdCGM79nwWnwpyUkkCl55JKinkWDHgZFRX4
6vRvu+pPEZm3q/LEdqLqFej2PgDiPrU5Xz/inD5JJpfWsALOQ2F2FrBx7IIK9Eat
Mfwp84p5wwqBaSi2l54CBPdbvKhL4K5cI+Mo3TKi7xtIWp+pfghVsFdZvP18I8K6
MXFALNc/GNRqZaHnu3glODVvHY0IRm/ON/nbZQIDAQABAoIBAQCjEJwDFdu3qx00
kT8vc0K6Niftd4BIciSlzkSOXFDmFyg4nX0onngcKL/5Zpmhm4oZLm4sXddYvn0s
MpxL6dRbA9M6wZqh1fEzBNPnS1S5IsV0rcIveDtYSW92iJeAJgSmwzNTtw8E50M1
LR5QsPf+xqn2zLuAMaHU3BHvnUYEVarFp53vtj0G4fXSTDNon8HEO3sIXvvKudqR
qe29ySIg+mS5qeBfF3JP63gDndSNT7/1JDhDPxWEm58S8FmwEaaU/RTOM0wxMUyV
jR/ZdYeJMir2JZLHC/5IrK67+nDTD0OuNqMWZzCkWEwEMU+SLJazMCkRX471wFoH
j8mhiaElAoGBAOc5fx4RkNS0hriXtsoWNCia0HnYQMu8vHb8xFC51lQdfA2slcX9
Y7AB6q6J3eJUi2l2fjzsUt9TANa6ab2O6g3KmWDT50Pxq0QFMOOPJVjDbn1QdGZA
A2AUp+YzPuWJ0/gf1FD0tTk5bbv3p7CB/muR4/TO6LPyDnK5mrSkZ9UXAoGBAM4t
Nr3ET3PEgHFGcn3KQzH14JU5A1Ma8sBP4RfKZyufU2I8ybN6RG4PcKHyzOB9/Gpb
607UsN4DdDQVSUdb8xaQXBZaQBGQMgqCAmjjqIQxy001jLvIHYcUxm26ZHyv0XwR
e3JkEwPqPaDa6LKlnhDGcs6G6J4B4aUR88EvmFjjAoGBAMySqyvwQKJgQh2JZRiw
wl72ceKLePCIwHnJsur1MHJlT79NZYmxYQR0/ayEn8JCKMIbKx89uyiI6GIStcEX
c27WRBNOB/uuEmfw68s5d8JrzhKjHwjkM9hLDi12Q3yUD+0kRBWIG9pQPA0k1MEu
kemcPwH2Gh4y16ObIQwXtSHrAoGAV+Hf5o2qDD+jPCV6IfI4KDCVNSYjK6Zd+OlT
mg91YJu+MC6XD0C7sGo2aWGUQNCS6kcaCvUQGuJAAv9bx+YCvQh1qDV5/8KGAgKe
wlTf/NE4xkVgIp7PL0gEuLrtoFRVJ9xP0Vek31NWR51n+NYthRsBztSkjM1igDkh
vKPr/V8CgYEAvSdqrkc1d5meeanRZjUPwcP2FGau2NHkaxcTlEYG6VWqpuGjoXEg
ANITv+uG8t+Uhkcx/xhborHEk/1xrqfpxeJQTPyGtwJHA3hR0/2NRA8RP5Cu0atZ
DOdjuNFpMQpZ4OQXYjYDU4GEgaCVBZkJEosz/kl32qQqidfyCS7IsbE=
-----END RSA PRIVATE KEY-----`
	pubKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAujka1xH+mye8zIlAzFuQ
6UsVLKWiNDsINBTDSPHZ3qQujakL3LQQjYpiKMH5qE/inVJncZ/rQl1bspxSjncU
4GLYu6GBcMwlbDNCpH1DZwXuZQlrwCXghFnFVqdtnPJ28Oq55oDDjLQrQVdCGM79
nwWnwpyUkkCl55JKinkWDHgZFRX46vRvu+pPEZm3q/LEdqLqFej2PgDiPrU5Xz/i
nD5JJpfWsALOQ2F2FrBx7IIK9EatMfwp84p5wwqBaSi2l54CBPdbvKhL4K5cI+Mo
3TKi7xtIWp+pfghVsFdZvP18I8K6MXFALNc/GNRqZaHnu3glODVvHY0IRm/ON/nb
ZQIDAQAB
-----END PUBLIC KEY-----`
	testBody = `this is the content`
)

func TestPubEncrypt(t *testing.T) {
	var rsaObj codec.Crypter
	var err error
	rsaObj, err = codec.NewRsa(pubKey, priKey)
	assert.Nil(t, err)

	enc, err := rsaObj.Encrypt([]byte(testBody))
	assert.Nil(t, err)

	dec, err := rsaObj.Decrypt(enc)
	assert.Nil(t, err)
	assert.Equal(t, testBody, string(dec))
}

func TestPrivEncrypt(t *testing.T) {
	rsaObj, err := codec.NewRsa(pubKey, priKey)
	assert.Nil(t, err)

	enc, err := rsaObj.EncryptEx([]byte(testBody))
	assert.Nil(t, err)

	dec, err := rsaObj.DecryptEx(enc)
	assert.Nil(t, err)
	assert.Equal(t, testBody, string(dec))
}

func TestPubEncryptWithBase64(t *testing.T) {
	var rsaObj codec.Crypter
	var err error
	rsaObj, err = codec.NewRsa(pubKey, priKey)
	assert.Nil(t, err)

	enc, err := rsaObj.EncryptBase64([]byte(testBody))
	assert.Nil(t, err)

	dec, err := rsaObj.DecryptBase64(enc)
	assert.Nil(t, err)
	assert.Equal(t, testBody, string(dec))
}

func TestPrivEncryptWithBase64(t *testing.T) {
	rsaObj, err := codec.NewRsa(pubKey, priKey)
	assert.Nil(t, err)

	enc, err := rsaObj.EncryptBase64Ex([]byte(testBody))
	assert.Nil(t, err)

	dec, err := rsaObj.DecryptBase64Ex(enc)
	assert.Nil(t, err)
	assert.Equal(t, testBody, string(dec))
}

func TestBadPubKey(t *testing.T) {
	_, err := codec.NewRsa("foo", "boo")
	assert.Equal(t, codec.ErrPrivKeyNotRsa, err)
}

func TestCreateRsaPubEncrypt(t *testing.T) {
	content := strings.Repeat("H", 5) + "e"

	var rsaObj codec.Crypter
	var err error
	rsaObj, err = codec.CreateRsa(1024)
	// rsaObj, err = codec.CreateRsaPkcs8(2048)
	assert.Nil(t, err)

	secretData, err := rsaObj.Encrypt([]byte(content))
	if err != nil {
		fmt.Println(err)
	}
	plainData, err := rsaObj.Decrypt(secretData)
	if err != nil {
		fmt.Print(err)
	}

	assert.Equal(t, content, string(plainData))

	fmt.Printf("内容：%v\n加密:%v\n解密:%v\n",
		content,
		hex.EncodeToString(secretData),
		string(plainData),
	)
}

func TestCreateRsaSign(t *testing.T) {
	data := strings.Repeat("H", 5) + "e"
	content := []byte(strings.Repeat(data, 3))

	var rsaObj codec.Signer
	var err error
	rsaObj, err = codec.CreateRsa(1024)
	// rsaObj, err = codec.CreateRsaPkcs8(2048)
	assert.Nil(t, err)

	//sign,_ := rsaObj.Sign(content, crypto.SHA1)
	//verify := rsaObj.Verify(content, sign, crypto.SHA1)

	sign, _ := rsaObj.Sign(content, crypto.SHA256)
	verify := rsaObj.Verify(content, sign, crypto.SHA256)

	assert.True(t, verify)

	fmt.Printf("内容：%v\n签名:%v\n验签结果:%v\n",
		string(content),
		hex.EncodeToString(sign),
		verify,
	)
}
