package unit

type UnityConversionInterface interface {
	ToByte(string, ByteUnit) (string, error)
}
