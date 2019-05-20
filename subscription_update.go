package paddle

import (
	"context"
)

type SubscriptionUpdateOptions struct {
	VendorID       int    `url:"vendor_id"`
	VendorAuthCode string `url:"vendor_auth_code"`

	// https://paddle.com/docs/subscription-update-api/
	SubscriptionID  int    `url:"subscription_id,omitempty"` // required
	Quantity        int    `url:"quantity,omitempty"`        // required
	RecurringPrice  string `url:"recurring_price,omitempty"`
	Currency        string `url:"currency,omitempty"`
	BillImmediately bool   `url:"bill_immediately,omitempty"`
	PlanID          int    `url:"plan_id,omitempty"`
	Prorate         bool   `url:"prorate,omitempty"`
	KeepModifiers   bool   `url:"keep_modifiers,omitempty"`
}

type SubscriptionUpdate struct {
	SubscriptionID int     `json:"subscription_id"`
	PlanID         int     `json:"plan_id"`
	UserID         int     `json:"user_id"`
	NextPayment    Payment `json:"next_payment"`
}

type SubscriptionUpdateResponse struct {
	Success  bool               `json:"success"`
	Response SubscriptionUpdate `json:"response"`
}

func (s *SubscriptionService) Update(ctx context.Context, options *SubscriptionUpdateOptions) (*SubscriptionUpdateResponse, error) {
	options.VendorID = s.client.conf.VendorID
	options.VendorAuthCode = s.client.conf.APIKey
	u, err := addOptions("subscription/users/update", options)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	update := new(SubscriptionUpdateResponse)
	_, err = s.client.Do(ctx, req, update)

	return update, err
}
