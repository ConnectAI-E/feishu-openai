package openai

import (
	"fmt"
	"net/http"
	"time"
)

type BillingSubScrip struct {
	HardLimitUsd float64 `json:"hard_limit_usd"`
	AccessUntil  float64 `json:"access_until"`
}
type BillingUsage struct {
	TotalUsage float64 `json:"total_usage"`
}

type BalanceResponse struct {
	TotalGranted   float64   `json:"total_granted"`
	TotalUsed      float64   `json:"total_used"`
	TotalAvailable float64   `json:"total_available"`
	EffectiveAt    time.Time `json:"effective_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

func (gpt *ChatGPT) GetBalance() (*BalanceResponse, error) {
	fmt.Println("进入")
	var data1 BillingSubScrip
	err := gpt.sendRequestWithBodyType(
		gpt.ApiUrl+"/v1/dashboard/billing/subscription",
		http.MethodGet,
		nilBody,
		nil,
		&data1,
	)
	fmt.Println("出错1", err)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing subscription: %v", err)
	}
	nowdate := time.Now()
	enddate := nowdate.Format("2006-01-02")
	startdate := nowdate.AddDate(0, 0, -100).Format("2006-01-02")
	var data2 BillingUsage
	err = gpt.sendRequestWithBodyType(
		gpt.ApiUrl+fmt.Sprintf("/v1/dashboard/billing/usage?start_date=%s&end_date=%s", startdate, enddate),
		http.MethodGet,
		nilBody,
		nil,
		&data2,
	)
	fmt.Println(data2)
	fmt.Println("出错2", err)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing subscription: %v", err)
	}

	balance := &BalanceResponse{
		TotalGranted:   data1.HardLimitUsd,
		TotalUsed:      data2.TotalUsage / 100,
		TotalAvailable: data1.HardLimitUsd - data2.TotalUsage/100,
		ExpiresAt:      time.Now(),
		EffectiveAt:    time.Now(),
	}

	if data1.AccessUntil > 0 {
		balance.EffectiveAt = time.Now()
		balance.ExpiresAt = time.Unix(int64(data1.AccessUntil), 0)
	}

	return balance, nil
}
