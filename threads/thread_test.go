package main

import (
	"learning-web-chatboard2/common"
	"log"
	"net/http"
	"strings"
	"testing"
)

func Test_CreateThread(t *testing.T) {
	newThre := common.Thread{
		Topic:  "I want eat meat pretty much.",
		Owner:  "TestingTaro",
		UserId: 1,
	}
	client := http.DefaultClient
	req, err := common.MakeRequestFromThread(
		&newThre,
		http.MethodPost,
		"http://localhost:8082/create",
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
	createdThre, err := common.MakeThreadFromResponse(res)
	if err != nil {
		log.Panicln(err.Error())
	}

	if strings.Compare(newThre.Topic, createdThre.Topic) != 0 ||
		strings.Compare(newThre.Owner, createdThre.Owner) != 0 {
		log.Panicln("found different fields")
	}
}
