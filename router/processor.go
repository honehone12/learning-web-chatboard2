package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const pwSalt = "jTfoU75lMD5awz93AdUmgKKoxuVtgvRi"
const runeSource = "aA1bB2cC3dD4eE5fFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvV6wW7xX8yY9zZ"

func encrypt(plainText string) (crypted string) {
	asBytes := sha256.Sum256([]byte(plainText))
	crypted = fmt.Sprintf("%x", asBytes)
	return
}

func processPassword(pw string) string {
	return encrypt(fmt.Sprint(pwSalt, pw))
}

func generate(length uint) (str string, err error) {
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
