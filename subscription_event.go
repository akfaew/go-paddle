package paddle

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
)

// https://developer.paddle.com/webhook-reference/product-fulfillment/fulfillment-webhook
type FulfillmentWebhook struct {
	EventTime   string `json:"event_time"`
	Quantity    int    `json:"quantity"`
	Passthrough string `json:"passthrough"`
}

// https://paddle.com/docs/subscriptions-event-reference/#subscription_created
type SubscriptionCreated struct {
	SubscriptionID     string `json:"subscription_id"`
	Status             string `json:"status"`
	Email              string `json:"email"`
	MarketingConsent   string `json:"marketing_consent"`
	SubscriptionPlanID string `json:"subscription_plan_id"`
	NextBillDate       string `json:"next_bill_date"`
	Passthrough        string `json:"passthrough"`
	UpdateURL          string `json:"update_url"`
	CancelURL          string `json:"cancel_url"`
	Currency           string `json:"currency"`
	CheckoutID         string `json:"checkout_id"`
	Quantity           string `json:"quantity"`
	UnitPrice          string `json:"unit_price"`
	EventTime          string `json:"event_time"`
}

// https://paddle.com/docs/subscriptions-event-reference/#subscription_cancelled
type SubscriptionCancelled struct {
	SubscriptionID            string `json:"subscription_id"`
	Status                    string `json:"status"`
	Email                     string `json:"email"`
	MarketingConsent          string `json:"marketing_consent"`
	SubscriptionPlanID        string `json:"subscription_plan_id"`
	CancellationEffectiveDate string `json:"cancellation_effective_date"`
	Passthrough               string `json:"passthrough"`
	UserID                    string `json:"user_id"`
	CheckoutID                string `json:"checkout_id"`
	Quantity                  string `json:"quantity"`
	UnitPrice                 string `json:"unit_price"`
	EventTime                 string `json:"event_time"`
	Currency                  string `json:"currency"`
}

// https://paddle.com/docs/subscriptions-event-reference/#subscription_payment_succeeded
type SubscriptionPaymentSucceeded struct {
	CheckoutID         string `json:"checkout_id"`
	Currency           string `json:"currency"`
	Email              string `json:"email"`
	EventTime          string `json:"event_time"`
	MarketingConsent   string `json:"marketing_consent"`
	NextBillDate       string `json:"next_bill_date"`
	Passthrough        string `json:"passthrough"`
	Quantity           string `json:"quantity"`
	Status             string `json:"status"`
	SubscriptionID     string `json:"subscription_id"`
	SubscriptionPlanID string `json:"subscription_plan_id"`
	UnitPrice          string `json:"unit_price"`
	UserID             string `json:"user_id"`
}

func phpserialize(form url.Values) []byte {
	var keys []string
	for k := range form {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	serialized := fmt.Sprintf("a:%d:{", len(keys))
	for _, k := range keys {
		serialized += fmt.Sprintf("s:%d:\"%s\";s:%d:\"%s\";", len(k), k, len(form.Get(k)), form.Get(k))
	}
	serialized += "}"

	return []byte(serialized)
}

// https://paddle.com/docs/reference-verifying-webhooks/
func ValidatePayload(r *http.Request, pubkey *rsa.PublicKey) (interface{}, error) {
	payload := map[string]string{}

	// Get the p_signature parameter and base64 decode it.
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	p_signature := r.Form.Get("p_signature")
	signature, err := base64.StdEncoding.DecodeString(p_signature)
	if err != nil {
		return payload, err
	}

	// Get the fields sent in the request, and remove the p_signature parameter
	r.Form.Del("p_signature")
	for k := range r.Form {
		payload[k] = r.Form.Get(k) // r.Form is a map[string][]string
	}

	// ksort() and serialize the fields
	hashed := sha1.Sum(phpserialize(r.Form))

	if err = rsa.VerifyPKCS1v15(pubkey, crypto.SHA1, hashed[:], signature); err != nil {
		return nil, err
	}

	switch r.Form.Get("alert_name") {
	case "subscription_created":
		var ret SubscriptionCreated
		if j, err := json.Marshal(payload); err != nil {
			return nil, err
		} else {
			if err := json.Unmarshal(j, &ret); err != nil {
				return nil, err
			}
		}
		return &ret, nil
	case "subscription_cancelled":
		var ret SubscriptionCancelled
		if j, err := json.Marshal(payload); err != nil {
			return nil, err
		} else {
			if err := json.Unmarshal(j, &ret); err != nil {
				return nil, err
			}
		}
		return &ret, nil
	case "subscription_payment_succeeded":
		var ret SubscriptionPaymentSucceeded
		if j, err := json.Marshal(payload); err != nil {
			return nil, err
		} else {
			if err := json.Unmarshal(j, &ret); err != nil {
				return nil, err
			}
		}
		return &ret, nil
	default:
		return nil, nil
	}
}

// https://paddle.com/docs/reference-verifying-webhooks/
// FulfillmentWebhook does not have an alert_type, so handle it in a separate url
func ValidateFulfillmentWebhookPayload(r *http.Request, pubkey *rsa.PublicKey) (*FulfillmentWebhook, error) {
	payload := map[string]string{}

	// Get the p_signature parameter and base64 decode it.
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	p_signature := r.Form.Get("p_signature")
	signature, err := base64.StdEncoding.DecodeString(p_signature)
	if err != nil {
		return nil, err
	}

	// Get the fields sent in the request, and remove the p_signature parameter
	r.Form.Del("p_signature")
	for k := range r.Form {
		payload[k] = r.Form.Get(k) // r.Form is a map[string][]string
	}

	// ksort() and serialize the fields
	hashed := sha1.Sum(phpserialize(r.Form))

	if err = rsa.VerifyPKCS1v15(pubkey, crypto.SHA1, hashed[:], signature); err != nil {
		return nil, err
	}

	ret := new(FulfillmentWebhook)
	if j, err := json.Marshal(payload); err != nil {
		return nil, err
	} else {
		if err := json.Unmarshal(j, ret); err != nil {
			return nil, err
		}
	}
	return ret, nil
}
