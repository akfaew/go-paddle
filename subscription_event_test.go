package paddle

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCancellationEffectiveDate(t *testing.T) {
	e := SubscriptionCancelled{
		CancellationEffectiveDate: "2021-05-12",
	}
	val := e.GetCancellationEffectiveDate()
	assert.Equal(t, 2021, val.Year())
	assert.Equal(t, time.Month(5), val.Month())
	assert.Equal(t, 12, val.Day())
}
