package resource

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type signOffReview struct {
	result bool   `json:"result"`
	uuid string `json:"uuid"`
}

func signOffPost(signoffReq *robbieSignOffRequest) (bool, error) {
	requestBody, err := json.Marshal(&signoffReq)
	if err != nil {
		return false, err
	}
	authAuthenticatorUrl := ""
	r, err := http.NewRequest("POST", authAuthenticatorUrl, bytes.NewBuffer(requestBody))
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	review := &signOffReview{}
	derr := json.NewDecoder(res.Body).Decode(review)
	if derr != nil {
		return false, derr
	}
	return review.result, nil

}
