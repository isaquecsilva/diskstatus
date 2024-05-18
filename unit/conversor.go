package unit

import (
	"math"
	"strconv"
)

type ByteUnit int

var (
	Byte  = ByteUnit(1)
	KByte = ByteUnit(1024)
	MByte = ByteUnit(math.Pow(1024, 2))
	GByte = ByteUnit(math.Pow(1024, 3))
	TByte = ByteUnit(math.Pow(1024, 4))
)

type Conversor bool

func (Conversor) ToByte(str string, targetUnity ByteUnit) (string, error) {
	num, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(num) / int(targetUnity)), nil
}
