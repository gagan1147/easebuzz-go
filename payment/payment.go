package payment

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/CloudStuffTech/go-utils/request"
	"github.com/CloudStuffTech/go-utils/security"
)

const EasebuzzURL = "https://pay.easebuzz.in/payment/initiateLink"

type Client struct {
	Key  string
	Salt string
}

type Response struct {
	Status    int    `json:"status"`
	Data      string `json:"data"`
	ErrorDesc string `json:"error_desc"`
}

type PaymentParams struct {
	Txnid       string
	Amount      string
	Productinfo string
	Firstname   string
	Email       string
	Phone       string
	Udf1        string
	Udf2        string
	Udf3        string
	Udf4        string
	Udf5        string
	Udf6        string
	Udf7        string
	Udf8        string
	Udf9        string
	Udf10       string
	Furl        string
	Surl        string
}

func (c Client) InitiatePayment(p *PaymentParams) (Response, error) {
	response := Response{}
	easeBuzzPreHash := c.Key + "|"
	form := url.Values{}
	form.Add("key", c.Key)
	var fixedFormParams = []string{"txnid", "amount", "productinfo", "firstname", "email", "phone", "udf1", "udf2", "udf3", "udf4", "udf5", "udf6", "udf7", "udf8", "udf9", "udf10", "furl", "surl"}
	for _, fixedFormParam := range fixedFormParams {
		value := getField(p, strings.Title(fixedFormParam))
		if value.IsValid() && !value.IsZero() {
			typeCastedValue := value.String()
			if !(fixedFormParam == "phone" || fixedFormParam == "furl" || fixedFormParam == "surl") {
				easeBuzzPreHash += fmt.Sprintf("%s|", typeCastedValue)
			}
			form.Add(fixedFormParam, typeCastedValue)
		} else {
			easeBuzzPreHash += "|"
		}
	}
	easeBuzzPreHash += c.Salt
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
		return response, err
	}
	json.Unmarshal([]byte(resp.Body), &response)
	return response, err
}

func getField(p *PaymentParams, field string) reflect.Value {
	r := reflect.ValueOf(p)
	return reflect.Indirect(r).FieldByName(field)
}
