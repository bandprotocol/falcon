package tunnel

import (
	"github.com/cosmos/gogoproto/proto"

	bandtsstypes "github.com/bandprotocol/falcon/internal/bandchain/bandtss"
	feedstypes "github.com/bandprotocol/falcon/internal/bandchain/feeds"
)

// register concrete types.
func init() {
	proto.RegisterType((*TSSRoute)(nil), "band.tunnel.v1beta1.TSSRoute")
	proto.RegisterType((*TSSPacketReceipt)(nil), "band.tunnel.v1beta1.TSSPacketReceipt")
	proto.RegisterType((*IBCRoute)(nil), "band.tunnel.v1beta1.IBCRoute")
	proto.RegisterType((*IBCPacketReceipt)(nil), "band.tunnel.v1beta1.IBCPacketReceipt")
	proto.RegisterType((*TunnelPricesPacketData)(nil), "band.tunnel.v1beta1.TunnelPricesPacketData")
}

// TSSRoute represents a route for TSS packets and implements the RouteI interface.
type TSSRoute struct {
	// destination_chain_id is the destination chain ID
	DestinationChainID string `protobuf:"bytes,1,opt,name=destination_chain_id,json=destinationChainId,proto3"                 json:"destination_chain_id,omitempty"`
	// destination_contract_address is the destination contract address
	DestinationContractAddress string `protobuf:"bytes,2,opt,name=destination_contract_address,json=destinationContractAddress,proto3" json:"destination_contract_address,omitempty"`
	// encoder is the mode of encoding packet data.
	Encoder feedstypes.Encoder `protobuf:"varint,3,opt,name=encoder,proto3,enum=band.feeds.v1beta1.Encoder"                     json:"encoder,omitempty"`
}

func (m *TSSRoute) Reset()         { *m = TSSRoute{} }
func (m *TSSRoute) String() string { return proto.CompactTextString(m) }
func (*TSSRoute) ProtoMessage()    {}

// TSSPacketReceipt represents a receipt for a TSS packet and implements the PacketReceiptI interface.
type TSSPacketReceipt struct {
	// signing_id is the signing ID
	SigningID bandtsstypes.SigningID `protobuf:"varint,1,opt,name=signing_id,json=signingId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/bandtss.SigningID" json:"signing_id,omitempty"`
}

func (m *TSSPacketReceipt) Reset()         { *m = TSSPacketReceipt{} }
func (m *TSSPacketReceipt) String() string { return proto.CompactTextString(m) }
func (*TSSPacketReceipt) ProtoMessage()    {}

// IBCRoute represents a route for IBC packets and implements the RouteI interface.
type IBCRoute struct {
	// channel_id is the IBC channel ID
	ChannelID string `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
}

func (m *IBCRoute) Reset()         { *m = IBCRoute{} }
func (m *IBCRoute) String() string { return proto.CompactTextString(m) }
func (*IBCRoute) ProtoMessage()    {}

// IBCPacketReceipt represents a receipt for a IBC packet and implements the PacketReceiptI interface.
type IBCPacketReceipt struct {
	// sequence is representing the sequence of the IBC packet.
	Sequence uint64 `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

func (m *IBCPacketReceipt) Reset()         { *m = IBCPacketReceipt{} }
func (m *IBCPacketReceipt) String() string { return proto.CompactTextString(m) }
func (*IBCPacketReceipt) ProtoMessage()    {}

// TunnelPricesPacketData represents the IBC packet payload for the tunnel packet.
type TunnelPricesPacketData struct {
	// tunnel_id is the tunnel ID
	TunnelID uint64 `protobuf:"varint,1,opt,name=tunnel_id,json=tunnelId,proto3"   json:"tunnel_id,omitempty"`
	// sequence is representing the sequence of the tunnel packet.
	Sequence uint64 `protobuf:"varint,2,opt,name=sequence,proto3"                  json:"sequence,omitempty"`
	// prices is the list of prices information from feeds module.
	Prices []feedstypes.Price `protobuf:"bytes,3,rep,name=prices,proto3"                     json:"prices"`
	// created_at is the timestamp when the packet is created
	CreatedAt int64 `protobuf:"varint,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (m *TunnelPricesPacketData) Reset()         { *m = TunnelPricesPacketData{} }
func (m *TunnelPricesPacketData) String() string { return proto.CompactTextString(m) }
func (*TunnelPricesPacketData) ProtoMessage()    {}
