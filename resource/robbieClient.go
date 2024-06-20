package resource

import (
	"bytes"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/wish/kcd/registry"
	"net/http"
	"time"
)

const RobbieEndpoint = "https://robbie-dev.i.wish.com/signoff"

type robbieSignOffRequest struct {
	KcdName      string
	KcdNameSpace string
	KcdLables    map[string]string
	KcdTag       string
	KcdImageRepo string
	Versions     registry.Versions
	Digest       registry.Digest
}

type signOffReview struct {
	Result bool   `json:"result"`
	Uuid   string `json:"uuid"`
}

func signOffPost(signoffReq *robbieSignOffRequest, endpoint string) (*signOffReview, error) {
	requestBody, err := json.Marshal(signoffReq)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	review := &signOffReview{}
	derr := json.NewDecoder(res.Body).Decode(review)
	if derr != nil {
		return nil, derr
	}
	return review, nil

}

func signOffRetryalbe(signoffReq *robbieSignOffRequest, endpoint string) (bool, error) {
	// max attempts to Robbie Sign off is set as 3, and init sleep duration is 3 seconds
	attempts := 3
	sleep := 3
	var result *signOffReview
	var err error
	for i := 0; i < attempts; i++ {
		glog.V(2).Infof("Querying with Robbie to get sign-off review for the %v attempt", i)
		result, err = signOffPost(signoffReq, endpoint)
		if err == nil {
			glog.V(2).Infof("Successfully get with Robbie sign-off review as %v, with UUId as %v", result.Result, result.Uuid)
			return result.Result, nil
		} else {
			glog.V(2).Infof("Querying with Robbie to get sign-off review fails, sleep for %v second to retry", sleep)
			time.Sleep(time.Duration(sleep) * time.Second)
			sleep *= 2
		}
	}
	return false, err
}
