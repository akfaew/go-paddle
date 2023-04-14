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
	"time"
)

// https://developer.paddle.com/webhook-reference/product-fulfillment/fulfillment-webhook
type FulfillmentWebhook struct {
	EventTime   string `json:"event_time"`
	Quantity    string `json:"quantity"`
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
	UserID             string `json:"user_id"`
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

func (s *SubscriptionCancelled) GetCancellationEffectiveDate() time.Time {
	var year, month, day int
	_, err := fmt.Sscanf(s.CancellationEffectiveDate, "%d-%d-%d", &year, &month, &day)
	if err != nil {
		return time.Time{}
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
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

type SubscriptionUpdated struct {
	AlertID               string `json:"alert_id"`
	AlertName             string `json:"alert_name"`
	CancelURL             string `json:"cancel_url"`
	CheckoutID            string `json:"checkout_id"`
	Currency              string `json:"currency"`
	CustomDate            string `json:"custom_date"`
	EventTime             string `json:"event_time"`
	MarketingConsent      string `json:"marketing_consent"`
	NewPrice              string `json:"new_price"`
	NewQuantity           string `json:"new_quantity"`
	NewUnitPrice          string `json:"new_unit_price"`
	NewBillDate           string `json:"new_bill_date"`
	OldNextBillDate       string `json:"old_next_bill_date"`
	OldPrice              string `json:"old_price"`
	OldQuantity           string `json:"old_quantity"`
	OldStatus             string `json:"old_status"`
	OldSubscriptionPlanID string `json:"old_subscription_plan_id"`
	OldUnitPrice          string `json:"old_unit_price"`
	Status                string `json:"status"`
	SubscriptionID        string `json:"subscription_id"`
	SubscriptionPlanID    string `json:"subscription_plan_id"`
	UpdateURL             string `json:"update_url"`
	UserID                string `json:"user_id"`
	PausedAt              string `json:"paused_at"`
	PausedFrom            string `json:"paused_from"`
	PausedReason          string `json:"paused_reason"`
}

type SubscriptionPaymentFailed struct {
	AlertID               string `json:"alert_id"`
	AlertName             string `json:"alert_name"`
	Amount                string `json:"amount"`
	AttemptNumber         string `json:"attempt_number"`
	CancelURL             string `json:"cancel_url"`
	CheckoutID            string `json:"checkout_id"`
	Currency              string `json:"currency"`
	CustomData            string `json:"custom_data"`
	Email                 string `json:"email"`
	EventTime             string `json:"event_time"`
	Instalments           string `json:"instalments"`
	MarketingConsent      string `json:"marketing_consent"`
	NextRetryDate         string `json:"next_retry_date"`
	OrderID               string `json:"order_id"`
	UserID                string `json:"user_id"`
	Quantity              string `json:"quantity"`
	Status                string `json:"status"`
	SubscriptionID        string `json:"subscription_id"`
	SubscriptionPaymentID string `json:"subscription_payment_id"`
	SubscriptionPlanID    string `json:"subscription_plan_id"`
	UnitPrice             string `json:"unit_price"`
	UpdateURL             string `json:"update_url"`
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
	case "subscription_updated":
		var ret SubscriptionUpdated
		if j, err := json.Marshal(payload); err != nil {
			return nil, err
		} else {
			if err := json.Unmarshal(j, &ret); err != nil {
				return nil, err
			}
		}
		return &ret, nil
	case "subscription_payment_failed":
		var ret SubscriptionPaymentFailed
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
