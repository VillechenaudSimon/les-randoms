package riotinterface

import (
	"bytes"
	"errors"
	"io/ioutil"
	"les-randoms/utils"
	"net/http"
	"os"
)

func isResponseStatusNotOK(status string) bool {
	return status[0:3] != "200"
}

func requestRIOTAPI(url string) ([]byte, error) {
	utils.LogClassic("REQUEST RIOT API : " + url)

	request, error := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if error != nil {
		return nil, error
	}
	request.Header.Set("Accept-Charset", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("X-Riot-Token", os.Getenv("X_RIOT_TOKEN"))

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()

	if isResponseStatusNotOK(response.Status) {
		return nil, errors.New("RESPONSE RIOT API : " + response.Status)
	}

	//fmt.Println("response Headers:", response.Header)

	return ioutil.ReadAll(response.Body)
}
