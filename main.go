package main

import (
	"context"
	"log"

	"coinbase-advanced-automation/internal/coinbase"
	"coinbase-advanced-automation/internal/secret"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type EventRequest struct {
	Orders []coinbase.Order `json:"orders"`
}

func handler(ctx context.Context, event EventRequest) error {
	coinbaseSecret, err := secret.GetCoinbaseSecret(ctx)
	if err != nil {
		return err
	}

	client, err := coinbase.NewClient(coinbaseSecret)
	if err != nil {
		return err
	}

	for _, order := range event.Orders {
		req := coinbase.CreateOrderRequest{
			ClientOrderId: uuid.New().String(),
			ProductId:     order.ProductId,
			Side:          order.Side,
			OrderConfiguration: coinbase.OrderConfiguration{
				MarketMarketIoc: coinbase.MarketMarketIoc{
					QuoteSize: order.QuoteSize,
				},
			},
		}

		resp, err := client.CreateOrder(ctx, &req)
		if err != nil {
			return err
		}

		if resp.Success {
			log.Printf(
				"success=%v side=%s client_order_id=%s product_id=%s quote_size=%s",
				resp.Success,
				resp.SuccessResponse.Side,
				resp.SuccessResponse.ClientOrderId,
				resp.SuccessResponse.ProductId,
				resp.OrderConfiguration.MarketMarketIoc.QuoteSize,
			)
		} else {
			log.Printf(
				"success=%v side=%s product_id=%s quote_size=%s error=%s",
				resp.Success,
				order.Side,
				order.ProductId,
				order.QuoteSize,
				resp.ErrorResponse.ErrorDetails,
			)
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
