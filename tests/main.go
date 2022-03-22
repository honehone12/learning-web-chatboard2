package main

import (
	"crypto/rand"
	"fmt"
	"learning-web-chatboard2/common"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
)

const runeSource = "aA1bB2cC3dD4eE5fFgGhHiIjJkKlLm0MnNoOpPqQrRsStTuUvV6wW7xX8yY9zZ"

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

func main() {
	fmt.Println(generate(8))
}

func doFunc() {
	c := make(chan bool)
	go do(c)
	for {
		if <-c {
			os.Exit(0)
		}
	}
}

func do(c chan bool) {
	newUser := common.User{
		Name:     "TestingTaro",
		Email:    "TestingTaro@go.com",
		Password: "TaroTaroTesting0721",
	}
	client := http.DefaultClient

	req, err := common.MakeRequestFromUser(
		&newUser,
		"POST",
		"http://localhost:8081/users/post",
	)
	if err != nil {
		log.Panicln(err.Error())
	}
	res, err := client.Do(req)
	if err != nil {
		log.Panicln(err.Error())
	}
	if res.StatusCode != http.StatusOK {
		log.Panicln(res.Status)
	}
	createdUser, err := common.MakeUserFromResponse(res)
	if err != nil {
		log.Panicln(err.Error())
	}

	if strings.Compare(newUser.Name, createdUser.Name) != 0 ||
		strings.Compare(newUser.Email, createdUser.Email) != 0 {
		log.Panicln(err.Error())
	}
	c <- true
}
