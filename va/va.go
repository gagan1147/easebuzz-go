package va

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/CloudStuffTech/go-utils/request"
	"github.com/CloudStuffTech/go-utils/security"
)

const URL = "https://wire.easebuzz.in/api/v1/insta-collect/virtual_accounts/"

type Client struct {
	Key  string
	Salt string
}

type VA struct {
	ID               string            `json:"id"`
	AR               []string          `json:"authorized_remitters"`
	UPI              string            `json:"upi_qrcode_url"`
	Label            string            `json:"label"`
	Description      string            `json:"description"`
	VAN              string            `json:"virtual_account_number"`
	VaIFSC           string            `json:"virtual_ifsc_number"`
	VaUPI            string            `json:"virtual_upi_handle"`
	Active           string            `json:"is_active"`
	PhoneNum         []string          `json:"phone_numbers"`
	CreatedAt        time.Time         `json:"created_at"`
	ADA              string            `json:"auto_deactivate_at"`
	ServiceCharge    float32           `json:"service_charge"`
	GST              float32           `json:"gst_amount"`
	ServiceChargeGST float32           `json:"service_charge_with_gst"`
	BA               float32           `json:"balance_amount"`
	AccountType      int               `json:"account_type"`
	KF               bool              `json:"kyc_flow"`
	CB               string            `json:"created_by"`
	NS               map[string]string `json:"notification_settings"`
	UpiIMG           string            `json:"upi_qrcode_remote_file_location"`
	UpiPDF           string            `json:"upi_qrcode_scanner_remote_file_location"`
	BankName         string
}

type EaseBuzzVA struct {
	VA VA `json:"virtual_account"`
}

type Response struct {
	Success bool       `json:"success"`
	Data    EaseBuzzVA `json:"data"`
}

type VaParams struct {
	Label       string
	Description string
}

func (c Client) GetVA(vaID string) (*VA, error) {
	response := &Response{}
	response.Data.VA.BankName = "Yes Bank Ltd"
	authorizationHash := security.Sha512Hash(fmt.Sprintf("%s|%s|%s", c.Key, vaID, c.Salt))
	httpClient := request.NewClient(&request.ClientOptions{Timeout: 30})
	var opts = &request.RequestOptions{
		Method:        "GET",
		URL:           URL + vaID,
		Body:          "",
		Retries:       1,
		RetryInterval: 1,
		Query:         map[string]string{"key": c.Key},
		Headers:       map[string]string{"Authorization": authorizationHash},
	}
	resp, err := httpClient.Request(opts)
	if err != nil {
		return &response.Data.VA, err
	}
	json.Unmarshal([]byte(resp.Body), &response)
	return &response.Data.VA, err
}

func (c Client) CreateVA(params *VaParams) (*VA, error) {
	authorizationHash := security.Sha512Hash(fmt.Sprintf("%s|%s|%s", c.Key, params.Label, c.Salt))
	response := Response{}
	jsonBody, _ := json.Marshal(params)
	httpClient := request.NewClient(&request.ClientOptions{Timeout: 30})
	var opts = &request.RequestOptions{
		Method:        "POST",
		URL:           URL,
		Body:          string(jsonBody),
		Retries:       1,
		RetryInterval: 1,
		Headers:       map[string]string{"Authorization": authorizationHash, "Content-Type": "application/json; charset=UTF-8"},
	}
	resp, err := httpClient.Request(opts)
	if err != nil {
		return &response.Data.VA, err
	}
	json.Unmarshal([]byte(resp.Body), &response)
	if !response.Success {
		return &response.Data.VA, err
	}
	return &response.Data.VA, err
}
