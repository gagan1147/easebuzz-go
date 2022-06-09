package easebuzz

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/CloudStuffTech/go-utils/request"
	"github.com/CloudStuffTech/go-utils/security"
)

type EbzAccessKeyResp struct {
	Status    int    `json:"status"`
	Data      string `json:"data"`
	ErrorDesc string `json:"error_desc"`
}

const EasebuzzURL = "https://pay.easebuzz.in/payment/initiateLink"

func GetEbzAccessKey(key string, salt string, body map[string]string) EbzAccessKeyResp {
	ebzAccessKeyResp := EbzAccessKeyResp{}
	easeBuzzPreHash := key + "|"
	form := url.Values{}
	form.Add("key", key)
	var formParams = []string{"txnid", "amount", "productinfo", "firstname", "email", "phone", "udf1", "udf2", "udf3", "udf4", "udf5", "udf6", "udf7", "udf8", "udf9", "udf10", "furl", "surl"}
	for _, formParam := range formParams {
		if val, ok := body[formParam]; ok {
			if !(formParam == "phone" || formParam == "furl" || formParam == "surl") {
				easeBuzzPreHash += fmt.Sprintf("%s|", val)
			}
			form.Add(formParam, val)
		} else {
			easeBuzzPreHash += "|"
		}
	}
	easeBuzzPreHash += salt
	form.Add("hash", security.Sha512Hash(easeBuzzPreHash))
	httpClient := request.NewClient(&request.ClientOptions{Timeout: 30})
	var opts = &request.RequestOptions{
		Method:        "POST",
		URL:           EasebuzzURL,
		Body:          form.Encode(),
		Retries:       1,
		RetryInterval: 1,
		Headers:       map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	}
	resp, err := httpClient.Request(opts)
	if err != nil {
		return ebzAccessKeyResp
	}
	json.Unmarshal([]byte(resp.Body), &ebzAccessKeyResp)
	return ebzAccessKeyResp
}
