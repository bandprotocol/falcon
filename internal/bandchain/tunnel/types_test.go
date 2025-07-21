package tunnel_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bandprotocol/falcon/internal/bandchain/tunnel"
)

func TestIsTssRouteType(t *testing.T) {
	assert.True(t, tunnel.IsTssRouteType("/band.tunnel.v1beta1.TSSRoute"))
	assert.True(t, tunnel.IsTssRouteType("band.tunnel.v1beta1.TSSRoute"))
	assert.False(t, tunnel.IsTssRouteType("/band.tunnel.v1beta1.IBCRoute"))
}
