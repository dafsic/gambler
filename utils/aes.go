package utils

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"fmt"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesECBEncrypt(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, LineNo())
	}
	bs := block.BlockSize()
	src = PKCS5Padding(src, bs)
	if len(src)%bs != 0 {
		return nil, fmt.Errorf("Need a multiple of the blocksize,%s", LineNo())
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

func AesECBDecrypt(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, LineNo())
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return nil, fmt.Errorf("crypto/cipher: input not full blocks,%s", LineNo())
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	length := len(out)
	unpadding := int(out[length-1])
	if (length-unpadding) < 0 || (length-unpadding) > length {
		return nil, fmt.Errorf("crypto/cipher: unpadding error,%s", LineNo())
	}
	out = PKCS5UnPadding(out)
	return out, nil
}

func AesEncryptoBase64(data, key string) (string, error) {
	cdata, err := AesECBEncrypt([]byte(data), []byte(key))
	if err != nil {
		return "", fmt.Errorf("%w%s", err, LineNo())
	}

	return base64.StdEncoding.EncodeToString(cdata), nil

}

func AesDecryptooBase64(data, key string) (string, error) {
	cdata, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("%w%s", err, LineNo())
	}

	ptext, err := AesECBDecrypt(cdata, []byte(key))
	if err != nil {
		return "", fmt.Errorf("%w%s", err, LineNo())
	}

	return string(ptext), nil
}
