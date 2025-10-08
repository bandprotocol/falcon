package alert

import "fmt"

// Topic represents a hierarchical topic structure for alert.
type Topic struct {
	base      string
	tunnelID  *uint64
	chainName *string
	endpoint  *string
}

// NewTopic creates a new Topic with the given base topic.
func NewTopic(base string) *Topic {
	return &Topic{
		base: base,
	}
}

// WithTunnelID returns a new Topic with the tunnel ID appended.
func (t *Topic) WithTunnelID(tunnelID uint64) *Topic {
	return &Topic{
		base:      t.base,
		tunnelID:  &tunnelID,
		chainName: t.chainName,
		endpoint:  t.endpoint,
	}
}

// WithChainName returns a new Topic with the chain name appended.
func (t *Topic) WithChainName(chainName string) *Topic {
	return &Topic{
		base:      t.base,
		tunnelID:  t.tunnelID,
		chainName: &chainName,
		endpoint:  t.endpoint,
	}
}

// WithEndpoint returns a new Topic with the endpoint appended.
func (t *Topic) WithEndpoint(endpoint string) *Topic {
	return &Topic{
		base:      t.base,
		tunnelID:  t.tunnelID,
		chainName: t.chainName,
		endpoint:  &endpoint,
	}
}

// GetFullTopic constructs the full topic string by appending tunnel ID, chain name and endpoint if they exist.
func (t *Topic) GetFullTopic() string {
	fullTopic := t.base
	if t.tunnelID != nil {
		fullTopic = fmt.Sprintf("%s TUNNEL_ID-%d", fullTopic, *t.tunnelID)
	}

	if t.chainName != nil {
		fullTopic = fmt.Sprintf("%s CHAIN-%s", fullTopic, *t.chainName)
	}

	if t.endpoint != nil {
		fullTopic = fmt.Sprintf("%s ENDPOINT-%s", fullTopic, *t.endpoint)
	}
	return fullTopic
}
