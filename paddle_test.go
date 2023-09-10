package paddle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	c := Conf{}
	err := c.Init("nonexistant")
	require.Equal(t, err.Error(), "open nonexistant: no such file or directory")

	err = c.Init("testdata/invalid.pub")
	require.Equal(t, err.Error(), "failed to parse PEM block containing the public key")
}
