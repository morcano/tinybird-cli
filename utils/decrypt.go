package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

// BlockDecrypt performs ECB decryption
func BlockDecrypt(block cipher.Block, src []byte) []byte {
	dst := make([]byte, len(src))
	bs := block.BlockSize()

	if len(src)%bs != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	for i := 0; i < len(src); i += bs {
		block.Decrypt(dst[i:], src[i:])
	}

	return dst
}

// DecryptAES256ECB function decodes and decrypts AES-256-ECB encrypted content
func DecryptAES256ECB(encodedMessage, key string) (string, error) {
	ciphertext, _ := base64.StdEncoding.DecodeString(encodedMessage)

	if len(key) < 32 {
		return "", errors.New("key length is less than 32")
	}

	block, err := aes.NewCipher([]byte(key[:32]))
	if err != nil {
		return "", err
	}

	decrypted := BlockDecrypt(block, ciphertext)
	decrypted, err = RemovePKCS7Padding(decrypted)
	if err != nil {
		return "", fmt.Errorf("failed to remove PKCS7 padding: %v", err)
	}

	return string(decrypted), nil
}

// RemovePKCS7Padding removes padding from data
func RemovePKCS7Padding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("input data length cannot be 0")
	}
	paddingLength := int(data[length-1])
	if paddingLength > length {
		return nil, errors.New("error: padding is larger than block")
	}
	return data[:(length - paddingLength)], nil
}
