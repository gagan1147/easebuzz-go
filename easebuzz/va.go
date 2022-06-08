package easebuzz

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/CloudStuffTech/go-utils/request"
	"github.com/CloudStuffTech/go-utils/security"
)

const easebuzzVaURL = "https://wire.easebuzz.in/api/v1/insta-collect/virtual_accounts/"

type VA struct {
	ID          string    `json:"id"`
	UPI         string    `json:"upi_qrcode_url"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	VAN         string    `json:"virtual_account_number"`
	VaIFSC      string    `json:"virtual_ifsc_number"`
	VaUPI       string    `json:"virtual_upi_handle"`
	Active      string    `json:"is_active"`
	PhoneNum    []string  `json:"phone_numbers"`
	CreatedAt   time.Time `json:"created_at"`
}

type EaseBuzzVA struct {
	VA VA `json:"virtual_account"`
}

type EaseBuzzVAResp struct {
	Success bool       `json:"success"`
	Data    EaseBuzzVA `json:"data"`
}

func GetEasebuzzVA(key string, salt string, vaId string) *VA {
	easeBuzzVAResp := &EaseBuzzVAResp{}
	authorizationHash := security.Sha512Hash(fmt.Sprintf("%s|%s|%s", key, vaId, salt))
	httpClient := request.NewClient(&request.ClientOptions{Timeout: 30})
	var opts = &request.RequestOptions{
		Method:        "GET",
		URL:           easebuzzVaURL + vaId,
		Body:          "",
		Retries:       1,
		RetryInterval: 1,
		Query:         map[string]string{"key": key},
		Headers:       map[string]string{"Authorization": authorizationHash},
	}
	resp, err := httpClient.Request(opts)
	if err != nil {
		return &easeBuzzVAResp.Data.VA
	}
	json.Unmarshal([]byte(resp.Body), &easeBuzzVAResp)
	return &easeBuzzVAResp.Data.VA
}

func CreateEaseBuzzVA(salt string, body map[string]string) string {
	authorizationHash := security.Sha512Hash(fmt.Sprintf("%s|%s|%s", body["key"], body["label"], salt))
	jsonBody, _ := json.Marshal(body)
	httpClient := request.NewClient(&request.ClientOptions{Timeout: 30})
	var opts = &request.RequestOptions{
		Method:        "POST",
		URL:           easebuzzVaURL,
		Body:          string(jsonBody),
		Retries:       1,
		RetryInterval: 1,
		Headers:       map[string]string{"Authorization": authorizationHash, "Content-Type": "application/json; charset=UTF-8"},
	}
	resp, err := httpClient.Request(opts)
	if err != nil {
		return ""
	}
	easeBuzzVAResp := EaseBuzzVAResp{}
	json.Unmarshal([]byte(resp.Body), &easeBuzzVAResp)
	if easeBuzzVAResp.Success == false {
		return ""
	}
	return easeBuzzVAResp.Data.VA.ID
}
