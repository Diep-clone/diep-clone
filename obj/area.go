package obj

import (
	"app/lib"
	"encoding/binary"
)

// TODO : Area

var AreaList Area

type Area struct {
	X    float64
	Y    float64
	W    float64
	H    float64
	Team string
}

func AreaSocket() []byte {
	var b []byte = make([]byte, 17)

	b[0] = byte(1)
	binary.BigEndian.PutUint32(b[1:5], uint32(-lib.GameSetting.MapSize.X))
	binary.BigEndian.PutUint32(b[5:9], uint32(-lib.GameSetting.MapSize.Y))
	binary.BigEndian.PutUint32(b[9:13], uint32(lib.GameSetting.MapSize.X*2))
	binary.BigEndian.PutUint32(b[13:17], uint32(lib.GameSetting.MapSize.Y*2))

	return b
}
