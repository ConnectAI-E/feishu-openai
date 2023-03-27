package openai

import (
	"fmt"
	"net/http"
	"time"
)

//https://api.openai.com/dashboard/billing/credit_grants
type Billing struct {
	Object         string  `json:"object"`
	TotalGranted   float64 `json:"total_granted"`
	TotalUsed      float64 `json:"total_used"`
	TotalAvailable float64 `json:"total_available"`
	Grants         struct {
		Object string `json:"object"`
		Data   []struct {
			Object      string  `json:"object"`
			ID          string  `json:"id"`
			GrantAmount float64 `json:"grant_amount"`
			UsedAmount  float64 `json:"used_amount"`
			EffectiveAt float64 `json:"effective_at"`
			ExpiresAt   float64 `json:"expires_at"`
		} `json:"data"`
	} `json:"grants"`
}

type BalanceResponse struct {
	TotalGranted   float64   `json:"total_granted"`
	TotalUsed      float64   `json:"total_used"`
	TotalAvailable float64   `json:"total_available"`
	EffectiveAt    time.Time `json:"effective_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

func (gpt *ChatGPT) GetBalance() (*BalanceResponse, error) {
	var data Billing
	err := gpt.sendRequestWithBodyType(
		gpt.ApiUrl+"/dashboard/billing/credit_grants",
		http.MethodGet,
		nilBody,
		nil,
		&data,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing data: %v", err)
	}

	balance := &BalanceResponse{
		TotalGranted:   data.TotalGranted,
		TotalUsed:      data.TotalUsed,
		TotalAvailable: data.TotalAvailable,
		ExpiresAt:      time.Now(),
		EffectiveAt:    time.Now(),
	}

	if len(data.Grants.Data) > 0 {
		balance.EffectiveAt = time.Unix(int64(data.Grants.Data[0].EffectiveAt), 0)
		balance.ExpiresAt = time.Unix(int64(data.Grants.Data[0].ExpiresAt), 0)
	}

	return balance, nil
}
