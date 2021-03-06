package main

import (
	"learning-web-chatboard2/common"
	"log"
	"net/http"
	"strings"
	"testing"
)

func Test_CreateUser(t *testing.T) {
	newUser := common.User{
		Name:     "TestingTaro",
		Email:    "TestingTaro@go.com",
		Password: "TaroTaroTesting0721",
	}
	client := http.DefaultClient
	req, err := common.MakeRequestFromUser(
		&newUser,
		http.MethodPost,
		"http://localhost:8081/signup-account",
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
		log.Panicln("found different fields")
	}
}
