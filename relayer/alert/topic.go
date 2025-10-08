package alert

import "fmt"

type Topic struct {
	base      string
	tunnelID  *uint64
	chainName *string
	endpoint  *string
}

func NewTopic(base string) *Topic {
	return &Topic{
		base: base,
	}
}

func (t *Topic) WithTunnelID(tunnelID uint64) *Topic {
	return &Topic{
		base:      t.base,
		tunnelID:  &tunnelID,
		chainName: t.chainName,
		endpoint:  t.endpoint,
	}
}

func (t *Topic) WithChainName(chainName string) *Topic {
	return &Topic{
		base:      t.base,
		tunnelID:  t.tunnelID,
		chainName: &chainName,
		endpoint:  t.endpoint,
	}
}

func (t *Topic) WithEndpoint(endpoint string) *Topic {
	return &Topic{
		base:      t.base,
		tunnelID:  t.tunnelID,
		chainName: t.chainName,
		endpoint:  &endpoint,
	}
}

// GetFullTopic append the topic string with tunnel ID and chain name.
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
