package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
)

var cipherKey = `/=fbaL|)vOJPA6Xu.i_uvl*ieE=lx:g5`

func GCMCipher() cipher.AEAD {
	block, err := aes.NewCipher([]byte(cipherKey))
	if err != nil {
		log.Fatalln("Error creating cipher block:", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalln("Error creating gcm block:", err)
	}

	return gcm
}

func Encrypt(text string) string {
	gcm := GCMCipher()
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalln(err)
	}

	cipher := gcm.Seal(nonce, nonce, []byte(text), nil)
	return fmt.Sprintf("%s", cipher)
}

func Decrypt(cipherText []byte) string {
	gcm := GCMCipher()
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		log.Fatalln("GCM nonce err")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	text, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return string(text)
}
