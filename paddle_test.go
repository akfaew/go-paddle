package paddle

import (
	"testing"

	"github.com/akfaew/test"
)

func TestInit(t *testing.T) {
	c := Conf{}
	err := c.Init("nonexistant")
	test.EqualStr(t, err.Error(), "open nonexistant: no such file or directory")

	err = c.Init("testdata/invalid.pub")
	test.EqualStr(t, err.Error(), "failed to parse PEM block containing the public key")
}
