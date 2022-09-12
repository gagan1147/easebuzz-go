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

// Client structs holds the Payment Client
type Client struct {
	Key  string
	Salt string
}

// Reponse structs holds the reponse got from easebuzz payment initiate api
type Response struct {
	Status    int    `json:"status"`
	Data      string `json:"data"`
	ErrorDesc string `json:"error_desc"`
}

// PaymentParams structs is request parameter required for initiating the payment request
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

// TransactionWebhookReponse structs is webhook response which will get api call from easebuzz server to our server for every transaction status update
type TransactionWebhookReponse struct {
	Txnid          string `json:"txnid"`
	Firstname      string `json:"firstname"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Key            string `json:"key"`
	Mode           string `json:"mode"`
	Status         string `json:"status"`
	CardCategory   string `json:"cardCategory"`
	Addedon        string `json:"addedon"`
	PaymentSource  string `json:"payment_source"`
	PgType         string `json:"pg_type"`
	BankRefNum     string `json:"bank_ref_num"`
	Bankcode       string `json:"bankcode"`
	Error          string `json:"error"`
	NameOnCard     string `json:"name_on_card"`
	Cardnum        string `json:"cardnum"`
	CardType       string `json:"card_type"`
	Easepayid      string `json:"easepayid"`
	Amount         string `json:"amount"`
	NetAmountDebit string `json:"net_amount_debit"`
	Productinfo    string `json:"productinfo"`
	Udf1           string `json:"udf1"`
	Hash           string `json:"hash"`
	Surl           string `json:"surl"`
	Furl           string `json:"furl"`
	ErrorMsg       string `json:"error_Message"`
	MerchantLogo   string `json:"merchant_logo"`
}

// InitiatePayment function helps in generating access_key which is used for payment purpose
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
