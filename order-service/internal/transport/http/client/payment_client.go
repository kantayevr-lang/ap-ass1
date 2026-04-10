package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/internal/domain"
	"time"
)

type PaymentServiceClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewPaymentServiceClient(baseURL string) *PaymentServiceClient {
	return &PaymentServiceClient{
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *PaymentServiceClient) Pay(ctx context.Context, orderID string, amount int64) (string, error) {
	url := fmt.Sprintf("%s/payments", c.baseURL)

	payload := map[string]interface{}{
		"order_id": orderID,
		"amount":   amount,
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", domain.ErrPaymentServiceDown, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("payment declined with status: %d", resp.StatusCode)
	}

	var result struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Status, nil
}
