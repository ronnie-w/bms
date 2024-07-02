package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"embed"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"time"
)

// bms_windows.exe
//
//go:embed bms_unix
var Binaries embed.FS

func getMac() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}

	return as, nil
}

func encrypt() {
	var (
		f           *os.File
		err         error
		collectHash [][64]byte
		hashArr     []byte

		logger = func(strerr string, err error) {
			f.Write([]byte(time.Now().String() + strerr + err.Error()))
		}
	)
	defer f.Close()

	f, err = os.OpenFile("crash.dump", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		os.Exit(1)
	}
	mac, err := getMac()
	if err != nil || len(mac) == 0 {
		logger("Unsupported device!", err)
		os.Exit(2)
	}

	for _, m := range mac {
		collectHash = append(collectHash, sha512.Sum512([]byte(m)))
	}

	for _, ch := range collectHash {
		hashArr = append(hashArr, ch[:]...)
	}

	key := []byte("Hs&x+FS_LmT@Y^F#Hc-5hbJ'k_{$Lwp9")
	c, err := aes.NewCipher(key)
	if err != nil {
		logger("Fatal Error! Quiting (xe1)...", err)
		os.Exit(3)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		logger("Fatal Error! Quiting (xe2)...", err)
		os.Exit(3)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logger("Fatal Error! Quiting (xe3)...", err)
		os.Exit(3)
	}

	f, err = os.OpenFile("licence", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(1)
	}
	f.Write(gcm.Seal(nonce, nonce, hashArr, nil))

	os.Exit(0)
}

func dumpExe() {
	fmt.Println("Dumping exe...")
	if runtime.GOOS == "windows" {
		data, _ := Binaries.ReadFile("bms_windows.exe")
		f, _ := os.Create("bms_windows.exe")
		defer f.Close()

		f.Write(data)
	} else if runtime.GOOS == "linux" {
		data, _ := Binaries.ReadFile("bms_unix")
		f, _ := os.Create("bms_unix")
		defer f.Close()

		f.Write(data)
	}
}

func main() {
	fmt.Println("xoxo")
	dumpExe()
	encrypt()
}
