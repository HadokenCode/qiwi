// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package qiwi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"time"
)

// Cards for payment-history endpoints
type Cards struct {
	client *Client
}

// NewCards returns new Cards obj
func NewCards(c *Client) *Cards {
	return &Cards{client: c}
}

// CardsDetectResponse api response
type CardsDetectResponse struct {
	Code struct {
		Value json.Number `json:"value"`
		Name  string      `json:"_name"`
	} `json:"code"`
	Message json.Number `json:"message"`
}

// Detect detect card PS
func (c *Cards) Detect(cardNumber string) (id int64, err error) {
	body, err := c.client.makePostRequest(EndpointCardsDetect, url.Values{"cardNumber": {cardNumber}})
	if err != nil {
		return
	}
	defer body.Close()

	bts, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}

	buf := bytes.NewReader(bts)

	log.Printf("%s", bts)

	dec := json.NewDecoder(buf)

	var r CardsDetectResponse
	err = dec.Decode(&r)
	if err != nil {
		return
	}

	if r.Code.Value.String() != "0" {
		return 0, fmt.Errorf("%s", r.Message.String())
	}

	return r.Message.Int64()
}

// PaymentRequest request of payment
type PaymentRequest struct {
	ID  string `json:"id"`
	Sum struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	} `json:"sum"`
	PaymentMethod struct {
		Type      string `json:"type"`
		AccountID string `json:"accountId"`
	} `json:"paymentMethod"`
	Fields struct {
		Account string `json:"account"`
	} `json:"fields"`
}

// PaymentResponse foemat of payment response
type PaymentResponse struct {
	ID     string `json:"id"`
	Terms  string `json:"terms"`
	Fields struct {
		Account string `json:"account"`
	} `json:"fields"`
	Sum struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	} `json:"sum"`
	Source      string `json:"source"`
	Transaction struct {
		ID    string `json:"id"`
		State struct {
			Code string `json:"code"`
		} `json:"state"`
	} `json:"transaction"`
}

// Payment make mayment
func (c *Cards) Payment(psID int64, amount float64, cardNumber string) (res PaymentResponse, err error) {
	req := PaymentRequest{
		ID: strconv.Itoa(int(time.Now().Unix()) * 1000),
	}
	// constants
	req.PaymentMethod.Type = "Account"
	req.PaymentMethod.AccountID = "643"

	req.Sum.Amount = amount
	req.Sum.Currency = "643"
	req.Fields.Account = cardNumber

	endpoint := fmt.Sprintf(EndpointCardsPayment, psID)

	body, err := c.client.makePostRequest(endpoint, req)
	if err != nil {
		return
	}
	defer body.Close()

	bts, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}

	log.Printf("%s\n", bts)

	dec := json.NewDecoder(bytes.NewReader(bts))
	err = dec.Decode(&res)

	return
}
