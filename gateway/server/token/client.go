package token

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Client struct {
	host string
}

func New(host string) *Client {
	return &Client{
		host: host,
	}
}

const validatePath = "/api/token/validate"

type requestBody struct {
	Token string `json:"token"`
}

func (c *Client) Validate(ctx context.Context, token string) (bool, error) {
	reqBody := bytes.NewBuffer(nil)
	err := json.NewEncoder(reqBody).Encode(&requestBody{Token: token})
	if err != nil {
		panic(err.Error())
	}

	res, err := http.DefaultClient.Post(c.host+validatePath, "application/json", reqBody)
	if err != nil {
		return false, err
	}

	return res.StatusCode == http.StatusOK, nil
}
