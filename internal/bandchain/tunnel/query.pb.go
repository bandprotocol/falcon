package tunnel

import (
	"fmt"
	"io"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/gogoproto/proto"
)

// TunnelStatusFilter defines a filter for tunnel status.
type TunnelStatusFilter int32

const (
	// TUNNEL_STATUS_FILTER_UNSPECIFIED defines an unspecified status.
	TUNNEL_STATUS_FILTER_UNSPECIFIED TunnelStatusFilter = 0
	// TUNNEL_STATUS_FILTER_ACTIVE defines an active tunnel.
	TUNNEL_STATUS_FILTER_ACTIVE TunnelStatusFilter = 1
	// TUNNEL_STATUS_FILTER_INACTIVE defines an inactive tunnel.
	TUNNEL_STATUS_FILTER_INACTIVE TunnelStatusFilter = 2
)

var TunnelStatusFilter_name = map[int32]string{
	0: "TUNNEL_STATUS_FILTER_UNSPECIFIED",
	1: "TUNNEL_STATUS_FILTER_ACTIVE",
	2: "TUNNEL_STATUS_FILTER_INACTIVE",
}

var TunnelStatusFilter_value = map[string]int32{
	"TUNNEL_STATUS_FILTER_UNSPECIFIED": 0,
	"TUNNEL_STATUS_FILTER_ACTIVE":      1,
	"TUNNEL_STATUS_FILTER_INACTIVE":    2,
}

func (x TunnelStatusFilter) String() string {
	return proto.EnumName(TunnelStatusFilter_name, int32(x))
}

// QueryTunnelsRequest is the request type for the Query/Tunnels RPC method.
type QueryTunnelsRequest struct {
	// status_filter is a flag to filter tunnels by status.
	StatusFilter TunnelStatusFilter `protobuf:"varint,1,opt,name=status_filter,json=statusFilter,proto3,enum=band.tunnel.v1beta1.TunnelStatusFilter" json:"status_filter,omitempty"`
	// pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3"                                                                   json:"pagination,omitempty"`
}

func (m *QueryTunnelsRequest) Reset()         { *m = QueryTunnelsRequest{} }
func (m *QueryTunnelsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTunnelsRequest) ProtoMessage()    {}
func (m *QueryTunnelsRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryTunnelsResponse is the response type for the Query/Tunnels RPC method.
type QueryTunnelsResponse struct {
	// Tunnels is a list of tunnels.
	Tunnels []*Tunnel `protobuf:"bytes,1,rep,name=tunnels,proto3"    json:"tunnels,omitempty"`
	// pagination defines an optional pagination for the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryTunnelsResponse) Reset()         { *m = QueryTunnelsResponse{} }
func (m *QueryTunnelsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTunnelsResponse) ProtoMessage()    {}
func (m *QueryTunnelsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTunnelsResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryTunnelRequest is the request type for the Query/Tunnel RPC method.
type QueryTunnelRequest struct {
	// tunnel_id is the ID of the tunnel to query.
	TunnelId uint64 `protobuf:"varint,1,opt,name=tunnel_id,json=tunnelId,proto3" json:"tunnel_id,omitempty"`
}

func (m *QueryTunnelRequest) Reset()         { *m = QueryTunnelRequest{} }
func (m *QueryTunnelRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTunnelRequest) ProtoMessage()    {}

// QueryTunnelResponse is the response type for the Query/Tunnel RPC method.
type QueryTunnelResponse struct {
	// tunnel is the tunnel with the given ID.
	Tunnel Tunnel `protobuf:"bytes,1,opt,name=tunnel,proto3" json:"tunnel"`
}

func (m *QueryTunnelResponse) Reset()         { *m = QueryTunnelResponse{} }
func (m *QueryTunnelResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTunnelResponse) ProtoMessage()    {}
func (m *QueryTunnelResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

// QueryPacketRequest is the request type for the Query/Packet RPC method.
type QueryPacketRequest struct {
	// tunnel_id is the ID of the tunnel to query packets.
	TunnelId uint64 `protobuf:"varint,1,opt,name=tunnel_id,json=tunnelId,proto3" json:"tunnel_id,omitempty"`
	// sequence is the sequence of the packet to query.
	Sequence uint64 `protobuf:"varint,2,opt,name=sequence,proto3"                json:"sequence,omitempty"`
}

func (m *QueryPacketRequest) Reset()         { *m = QueryPacketRequest{} }
func (m *QueryPacketRequest) String() string { return proto.CompactTextString(m) }
func (*QueryPacketRequest) ProtoMessage()    {}

// QueryPacketResponse is the response type for the Query/Packet RPC method.
type QueryPacketResponse struct {
	// packet is the packet with the given tunnel ID and sequence.
	Packet *Packet `protobuf:"bytes,1,opt,name=packet,proto3" json:"packet,omitempty"`
}

func (m *QueryPacketResponse) Reset()         { *m = QueryPacketResponse{} }
func (m *QueryPacketResponse) String() string { return proto.CompactTextString(m) }
func (*QueryPacketResponse) ProtoMessage()    {}
func (m *QueryPacketResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *QueryTunnelsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryTunnelsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTunnelsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tunnels", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tunnels = append(m.Tunnels, &Tunnel{})
			if err := m.Tunnels[len(m.Tunnels)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (m *QueryTunnelResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryTunnelResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTunnelResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tunnel", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Tunnel.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (m *QueryPacketResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryPacketResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPacketResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Packet", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Packet == nil {
				m.Packet = &Packet{}
			}
			if err := m.Packet.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
