package config

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type generateKubeConfigResponse struct {
	Config string `json:"config"`
}

func GetKubeConfig(url string, token string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	sEnc := b64.StdEncoding.EncodeToString([]byte(token))
	req.Header.Set("Authorization", "Basic "+sEnc)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("Got response " + strconv.Itoa(resp.StatusCode) + " when expecting 200")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response generateKubeConfigResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		panic(err)
	}
	return []byte(response.Config), nil
}
