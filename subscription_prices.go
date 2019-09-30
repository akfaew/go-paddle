package paddle

import (
	"context"
	"errors"
)

var ErrCountryDoesNotExist = errors.New("Country does not exist")

type Price struct {
	Gross float64 `json:"gross"`
	Net   float64 `json:"net"`
	Tax   float64 `json:"tax"`
}

// https://developer.paddle.com/api-reference/checkout-api/prices/getprices
type SubscriptionPricesProduct struct {
	Currency     string `json:"currency"`
	ListPrice    Price  `json:"list_price"`
	Price        Price  `json:"price"`
	ProductID    int    `json:"product_id"`
	ProductTitle string `json:"product_title"`
	Subscription struct {
		Frequency int    `json:"frequency"`
		Interval  string `json:"interval"`
		ListPrice Price  `json:"list_price"`
		Price     Price  `json:"price"`
		TrialDays int    `json:"trial_days"`
	} `json:"subscription"`
	VendorSetPricesIncludedTax bool `json:"vendor_set_prices_included_tax"`
}

type SubscriptionPrices struct {
	CustomerCountry string                      `json:"customer_country"`
	Products        []SubscriptionPricesProduct `json:"products"`
}

type SubscriptionPricesResponse struct {
	Success  bool               `json:"success"`
	Response SubscriptionPrices `json:"response"`
}

type SubscriptionPricesOptions struct {
	ProductIDs      string `url:"product_ids,omitempty"`
	CustomerCountry string `url:"customer_country,omitempty"`
	CustomerIP      string `url:"customer_ip,omitempty"`
	Coupons         string `url:"coupons,omitempty"`
}

func (s *SubscriptionService) Prices(ctx context.Context, options SubscriptionPricesOptions) (*SubscriptionPricesResponse, error) {
	u, err := addOptions("prices", options)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	prices := new(SubscriptionPricesResponse)
	_, err = s.client.Do(ctx, req, prices)
	return prices, err
}
