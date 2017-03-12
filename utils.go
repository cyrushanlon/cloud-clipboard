package client

import (
	"errors"
	"io"
	"net"
	"runtime"
	"strings"

	"fmt"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"

	"encoding/base64"
	"encoding/hex"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
)

//StringArrayContains returns if the slice contains a target string
func StringArrayContains(list []string, target string) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == target {
			return true
		}
	}
	return false
}

//StringArrayRemove removes target string from slice
func StringArrayRemove(list *[]string, target string) {
	for i := 0; i < len(*list); i++ {
		if (*list)[i] == target {
			l := append((*list)[:i], (*list)[i+1:]...)
			list = &l
			return
		}
	}
}

//TrimBuffer takes a slice and trims it to the content
func TrimBuffer(buffer []byte) []byte {
	var bufferSlice []byte

	for k, v := range buffer {
		if v == 0 {
			bufferSlice = buffer[:k]
			break
		}
	}

	return bufferSlice
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

//HashString returns hashed salted string as bytes
func HashString(str string) string {
	salt := []byte("fixedSaltStringThatNooneShouldBeReading")

	return hex.EncodeToString(pbkdf2.Key([]byte(str), salt, 4096, sha256.Size, sha256.New))
}

func getKey() []byte {
	key := []byte(Conf.Password + Conf.Username)

	if len(key) > 32 {
		key = key[0:32]
	} else if len(key) == 0 {

	} else {
		for i := len(key) - 1; i < 31; i++ {
			key = append(key, 0)
		}
	}
	return key
}

func Encrypt(text []byte) ([]byte, error) {
	key := getKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func Decrypt(text []byte) ([]byte, error) {
	key := getKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getLogOutput(v ...interface{}) string {
	_, file, line, _ := runtime.Caller(2)

	split := strings.Split(file, "/")
	fileOut := split[len(split)-1]

	return fmt.Sprint("[", fileOut, " ", line, "] ", runtime.NumGoroutine(), fmt.Sprint(v...))
}

//LogInfo outputs info with function and line number of where it was called
func LogInfo(v ...interface{}) {
	log.Info(getLogOutput(v))
}

//LogErr outputs error with function and line number of where it was called
func LogErr(v ...interface{}) {
	log.Error(getLogOutput(v))
}

//LogWarn outputs warning with function and line number of where it was called
func LogWarn(v ...interface{}) {
	log.Warn(getLogOutput(v))
}
