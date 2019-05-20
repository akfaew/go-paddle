package paddle

import (
	"context"
)

type SubscriptionUser struct {
	SubscriptionID   int     `json:"subscription_id"`
	PlanID           int     `json:"plan_id"`
	UserID           int     `json:"user_id"`
	UserEmail        string  `json:"user_email"`
	MarketingConsent bool    `json:"marketing_consent"`
	State            string  `json:"state"`
	SignupDate       string  `json:"signup_date"`
	LastPayment      Payment `json:"last_payment"`
	NextPayment      Payment `json:"next_payment"`
}

type SubscriptionUsersResponse struct {
	Success  bool               `json:"success"`
	Response []SubscriptionUser `json:"response"`
}

type SubscriptionUsersOptions struct {
	VendorID       int    `url:"vendor_id"`
	VendorAuthCode string `url:"vendor_auth_code"`

	SubscriptionID string `url:"subscription_id,omitempty"`
	Plan           string `url:"plan,omitempty"`
}

func (s *SubscriptionService) Users(ctx context.Context, options *SubscriptionUsersOptions) (*SubscriptionUsersResponse, error) {
	options.VendorID = s.client.conf.VendorID
	options.VendorAuthCode = s.client.conf.APIKey
	u, err := addOptions("subscription/users", options)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	users := new(SubscriptionUsersResponse)
	_, err = s.client.Do(ctx, req, users)

	return users, err
}
