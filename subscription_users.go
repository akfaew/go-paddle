package paddle

import (
	"context"
)

// https://paddle.com/docs/api-list-users/
type SubscriptionUser struct {
	SubscriptionID   int     `json:"subscription_id"`
	PlanID           int     `json:"plan_id"`
	UserID           int     `json:"user_id"`
	UserEmail        string  `json:"user_email"`
	MarketingConsent bool    `json:"marketing_consent"`
	UpdateURL        string  `json:"update_url"`
	CancelURL        string  `json:"cancel_url"`
	State            string  `json:"state"`
	SignupDate       string  `json:"signup_date"`
	LastPayment      Payment `json:"last_payment"`
	NextPayment      Payment `json:"next_payment"`
}

type SubscriptionUsersResponse struct {
	Success  bool               `json:"success"`
	Response []SubscriptionUser `json:"response"`
}

// https://developer.paddle.com/api-reference/e33e0a714a05d-list-users
type SubscriptionUsersOptions struct {
	VendorID       int    `url:"vendor_id"`
	VendorAuthCode string `url:"vendor_auth_code"`

	SubscriptionID string `url:"subscription_id,omitempty"`
	PlanID         string `url:"plan_id,omitempty"`
	State          string `url:"state,omitempty"`
	ResultsPerPage string `url:"results_per_page,omitempty"` // max 200
	Page           string `url:"page,omitempty"`
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
