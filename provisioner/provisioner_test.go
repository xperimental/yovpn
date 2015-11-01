package provisioner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoToken(t *testing.T) {
	_, err := NewProvisioner("")
	assert.Equal(t, ErrNoToken, err)
}
