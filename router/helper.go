package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	runeSource = "aA1bB2cC3dD4eE5fFgGhHiIjJkKlLm0MnNoOpPqQrRsStTuUvV6wW7xX8yY9zZ"
	macSalt    = "uPUqL7dZ"
	pwSalt     = "LV2vP8vq"
)
const (
	aes256KeySize uint = 32
	macKeySize    uint = 32
	stateSize     uint = 32
)

var helper struct {
	block  cipher.Block
	macKey []byte
}

// every time server is restarted, cookie become no longer valid
func newProcessor() (err error) {
	bKyeStr, err := generateString(aes256KeySize)
	if err != nil {
		return
	}
	bKey := []byte(bKyeStr)
	macKeyStr, err := generateString(macKeySize)
	if err != nil {
		return
	}
	helper.macKey = []byte(macKeyStr)
	helper.block, err = aes.NewCipher(bKey)
	return
}

func makeHash(plainText string) (hashed string) {
	asBytes := sha256.Sum256([]byte(plainText))
	hashed = fmt.Sprintf("%x", asBytes)
	return
}

func processPassword(pw string) string {
	// see these pkgs
	// https://pkg.go.dev/golang.org/x/crypto/bcrypt
	// https://pkg.go.dev/golang.org/x/crypto/scrypt
	return makeHash(fmt.Sprint(pwSalt, pw))
}

func generateString(length uint) (str string, err error) {
	var i uint
	maxEx := int64(len(runeSource))
	runePool := []rune(runeSource)
	for i = 0; i < length; i++ {
		bigN, err := rand.Int(rand.Reader, big.NewInt(maxEx))
		if err != nil {
			break
		}
		n := bigN.Uint64()
		str = fmt.Sprint(str, string(runePool[n]))
	}
	return
}

func encrypt(plainText string) (cipherText []byte, err error) {
	cipherText = make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	n, err := io.ReadFull(rand.Reader, iv)
	if err != nil {
		err = fmt.Errorf("%s: returned %d", err.Error(), n)
		return
	}

	encryptStream := cipher.NewCTR(helper.block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))
	return
}

func decrypt(cipherText []byte) (plainText string, err error) {
	decryptText := make([]byte, len(cipherText[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(helper.block, cipherText[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptText, cipherText[aes.BlockSize:])
	plainText = string(decryptText)
	return
}

func makeMAC(value []byte) []byte {
	hash := hmac.New(sha256.New, helper.macKey)
	hash.Write(value)
	return hash.Sum([]byte(macSalt))
}

func verifyMAC(mac []byte, value []byte) bool {
	hash := hmac.New(sha256.New, helper.macKey)
	hash.Write(value)
	hashedVal := hash.Sum([]byte(macSalt))
	return hmac.Equal(mac, hashedVal)
}

func encode(value []byte) string {
	return base64.URLEncoding.EncodeToString(value)
}

func decode(encoded string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(encoded)
}

func storeSessionCookie(ctx *gin.Context, value string) (err error) {
	encrypted, err := encrypt(value)
	if err != nil {
		return
	}
	// add mac value first
	bytesVal := makeMAC(encrypted)
	// separated '|'
	bytesVal = append(bytesVal, []byte("|")...)
	// add encrypted value
	bytesVal = append(bytesVal, encrypted...)

	valToStore := encode(bytesVal)

	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(
		shortTimeSession,
		valToStore,
		0,
		"/",
		"localhost",
		true,
		true,
	)
	return
}

func pickupSessionCookie(ctx *gin.Context) (uuid string, err error) {
	rawStored, err := ctx.Cookie(shortTimeSession)
	if err != nil {
		return
	}
	bytesVal, err := decode(rawStored)
	if err != nil {
		return
	}
	splited := bytes.SplitN(bytesVal, []byte("|"), 2)
	mac := splited[0]
	encrypted := splited[1]
	if !verifyMAC(mac, encrypted) {
		err = errors.New("invalid cookie")
		return
	}
	uuid, err = decrypt(encrypted)
	return
}

// this should be middleware!!
func generateState(ctx *gin.Context) (state string, err error) {
	state, err = generateString(stateSize)
	sess, err := GetSessionPtr(ctx)
	if err != nil {
		return
	}
	sess.State = state

	return
}
