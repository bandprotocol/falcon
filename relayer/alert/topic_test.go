package alert_test

import (
	"testing"

	"github.com/bandprotocol/falcon/relayer/alert"
)

func TestGetFullTopic_Variants(t *testing.T) {
	t.Run("base only", func(t *testing.T) {
		topic := alert.NewTopic("ALERT")
		got := topic.GetFullTopic()
		want := "ALERT"
		if got != want {
			t.Fatalf("GetFullTopic()=%q, want %q", got, want)
		}
	})

	t.Run("all fields set", func(t *testing.T) {
		got := alert.NewTopic("ALERT").
			WithTunnelID(42).
			WithChainName("osmosis").
			WithEndpoint("https://rpc.osmosis.zone").
			GetFullTopic()

		want := "ALERT TUNNEL_ID-42 CHAIN-osmosis ENDPOINT-https://rpc.osmosis.zone"
		if got != want {
			t.Fatalf("GetFullTopic()=%q, want %q", got, want)
		}
	})

	t.Run("some fields set", func(t *testing.T) {
		got := alert.NewTopic("BASE").
			WithTunnelID(7).
			WithChainName("bandchain").
			GetFullTopic()

		want := "BASE TUNNEL_ID-7 CHAIN-bandchain"
		if got != want {
			t.Fatalf("GetFullTopic()=%q, want %q", got, want)
		}
	})
}

func TestFluentBuilder_Immutability(t *testing.T) {
	base := alert.NewTopic("A")

	withTunnel := base.WithTunnelID(1)
	if got, want := base.GetFullTopic(), "A"; got != want {
		t.Fatalf("base mutated after WithTunnelID: got %q want %q", got, want)
	}
	if got, want := withTunnel.GetFullTopic(), "A TUNNEL_ID-1"; got != want {
		t.Fatalf("withTunnel mismatch: got %q want %q", got, want)
	}

	withChain := withTunnel.WithChainName("chain")
	if got, want := withTunnel.GetFullTopic(), "A TUNNEL_ID-1"; got != want {
		t.Fatalf("withTunnel mutated after WithChainName: got %q want %q", got, want)
	}
	if got, want := withChain.GetFullTopic(), "A TUNNEL_ID-1 CHAIN-chain"; got != want {
		t.Fatalf("withChain mismatch: got %q want %q", got, want)
	}
}

func TestChaining_OrderIndependence(t *testing.T) {
	a := alert.NewTopic("BASE").
		WithChainName("c").
		WithTunnelID(5).
		WithEndpoint("e")

	b := alert.NewTopic("BASE").
		WithEndpoint("e").
		WithTunnelID(5).
		WithChainName("c")

	ga, gb := a.GetFullTopic(), b.GetFullTopic()
	want := "BASE TUNNEL_ID-5 CHAIN-c ENDPOINT-e"

	if ga != want {
		t.Fatalf("order A got %q want %q", ga, want)
	}
	if gb != want {
		t.Fatalf("order B got %q want %q", gb, want)
	}
}
