package utils

func Uint8ToString(uint8s []uint8) string {
	var bytes []byte
	for _, b := range uint8s {
		bytes = append(bytes, byte(b))
	}
	return string(bytes)
}

func StringToUint8(str string) uint8 {
	bytes := []byte(str)
	if len(bytes) != 1 {
		panic("length of string is too long to uint8")
	}
	return uint8(bytes[0])
}
