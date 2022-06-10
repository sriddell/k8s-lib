package rancher

import (
	b64 "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type config struct {
	Servers []RancherServer `json:"rancherServers"`
}

type RancherServer struct {
	RancherUrl string `json:"rancherUrl"`
	Token      string `json:"token"`
}

type clusterResult struct {
	Clusters []Cluster `json:"data"`
}

type Cluster struct {
	Name    string  `json:"name"`
	Id      string  `json:"id"`
	Actions actions `json:"actions"`
}

type actions map[string]string

func GetClusters(info RancherServer) []Cluster {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	url := info.RancherUrl + "/v3/clusters"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	sEnc := b64.StdEncoding.EncodeToString([]byte(info.Token))
	req.Header.Set("Authorization", "Basic "+sEnc)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("got code " + strconv.Itoa(resp.StatusCode) + " when attempting to get clusters for " + url)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response clusterResult
	err = json.Unmarshal(b, &response)
	if err != nil {
		panic(err)
	}
	return response.Clusters
}
