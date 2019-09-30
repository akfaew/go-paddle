package paddle

import (
	"context"
	"net/http"
	"testing"

	"github.com/akfaew/test"
)

func TestPrices(t *testing.T) {
	t.Skip("Provide your own data")

	client := NewCheckoutClient(context.Background(), &http.Client{})
	res, err := client.Subscription.Prices(context.Background(), SubscriptionPricesOptions{
		CustomerCountry: "PL",
		ProductIDs:      "1234,2345",
	})
	test.NoError(t, err)
	test.Fixture(t, res)
}
