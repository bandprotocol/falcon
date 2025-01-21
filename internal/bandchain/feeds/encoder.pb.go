package types

import "github.com/cosmos/gogoproto/proto"

// Encoder is an enumerator that defines the mode of encoding message in tss module.
type Encoder int32

const (
	// ENCODER_UNSPECIFIED is an unspecified encoder mode.
	ENCODER_UNSPECIFIED Encoder = 0
	// ENCODER_FIXED_POINT_ABI is a fixed-point price abi encoder (price * 10^9).
	ENCODER_FIXED_POINT_ABI Encoder = 1
	// ENCODER_TICK_ABI is a tick abi encoder.
	ENCODER_TICK_ABI Encoder = 2
)

var Encoder_name = map[int32]string{
	0: "ENCODER_UNSPECIFIED",
	1: "ENCODER_FIXED_POINT_ABI",
	2: "ENCODER_TICK_ABI",
}

var Encoder_value = map[string]int32{
	"ENCODER_UNSPECIFIED":     0,
	"ENCODER_FIXED_POINT_ABI": 1,
	"ENCODER_TICK_ABI":        2,
}

func (x Encoder) String() string {
	return proto.EnumName(Encoder_name, int32(x))
}
