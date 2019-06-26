package paddle

import (
	"context"
)

type ProductPayLink struct {
	URL string `json:"url"`
}

type ProductPayLinkResponse struct {
	Success  bool           `json:"success"`
	Response ProductPayLink `json:"response"`
}

type ProductGeneratePayLinkOptions struct {
	VendorID       int    `url:"vendor_id"`
	VendorAuthCode string `url:"vendor_auth_code"`

	// https://paddle.com/docs/api-custom-checkout/
	ProductID               int      `url:"product_id,omitempty"`
	Title                   string   `url:"title,omitempty"`
	WebhookURL              string   `url:"webhook_url,omitempty"`
	Prices                  []string `url:"prices,brackets,omitempty"`
	Locale                  string   `url:"locale,omitempty"`
	RecurringPrices         []string `url:"recurring_prices,omitempty"`
	TrialDays               int      `url:"trial_days,omitempty"`
	CustomMessage           string   `url:"custom_message,omitempty"`
	CouponCode              string   `url:"coupon_code,omitempty"`
	ImageURL                string   `url:"image_url,omitempty"`
	ReturnURL               string   `url:"return_url,omitempty"`
	QuantityVariable        int      `url:"quantity_variable,omitempty"`
	Quantity                int      `url:"quantity,omitempty"`
	Expires                 string   `url:"expires,omitempty"`
	Affiliates              string   `url:"affiliates,omitempty"`
	RecurringAffiliateLimit string   `url:"recurring_affiliate_limit,omitempty"`
	MarketingConsent        int      `url:"marketing_consent,omitempty"`
	CustomerEmail           string   `url:"customer_email,omitempty"`
	CustomerCountry         string   `url:"customer_country,omitempty"`
	CustomerPostcode        string   `url:"customer_postcode,omitempty"`
	VatCode                 string   `url:"vat_code,omitempty"`
	Passthrough             string   `url:"passthrough,omitempty"`
}

func (s *ProductService) GeneratePayLink(ctx context.Context, options *ProductGeneratePayLinkOptions) (*ProductPayLinkResponse, error) {
	options.VendorID = s.client.conf.VendorID
	options.VendorAuthCode = s.client.conf.APIKey
	options.ProductID = s.client.conf.ProductID
	u, err := addOptions("product/generate_pay_link", options)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	paylink := new(ProductPayLinkResponse)
	_, err = s.client.Do(ctx, req, paylink)

	return paylink, err
}

func (s *ProductService) GeneratePayLinkCustom(ctx context.Context, options *ProductGeneratePayLinkOptions) (*ProductPayLinkResponse, error) {
	options.VendorID = s.client.conf.VendorID
	options.VendorAuthCode = s.client.conf.APIKey
	u, err := addOptions("product/generate_pay_link", options)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	paylink := new(ProductPayLinkResponse)
	_, err = s.client.Do(ctx, req, paylink)

	return paylink, err
}
