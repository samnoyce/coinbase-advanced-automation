package coinbase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	httpMethod      = http.MethodPost
	createOrderPath = "/api/v3/brokerage/orders"
)

type CreateOrderRequest struct {
	ClientOrderId      string             `json:"client_order_id"`
	ProductId          string             `json:"product_id"`
	Side               string             `json:"side"`
	OrderConfiguration OrderConfiguration `json:"order_configuration"`
}

type CreateOrderResponse struct {
	Success            bool               `json:"success"`
	SuccessResponse    SuccessResponse    `json:"success_response,omitempty"`
	ErrorResponse      ErrorResponse      `json:"error_response,omitempty"`
	OrderConfiguration OrderConfiguration `json:"order_configuration"`
}

func (c *Client) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	signedToken, err := c.BuildJWT(httpMethod, createOrderPath)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(ctx, httpMethod, createOrderPath, signedToken, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("coinbase: unexpected status code %d in response", resp.StatusCode)
	}

	var createOrder CreateOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&createOrder); err != nil {
		return nil, err
	}

	return &createOrder, nil
}
