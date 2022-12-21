package mn

type MessageType uint8

const (
	MessageTypeNone MessageType = iota // 0
	MessageTypeConnect
	MessageTypeConnectAck
	MessageTypeConnectRefuse
	MessageTypePing
	MessageTypePong // 5
	MessageTypeCodePush
	MessageTypeCodeAck
	MessageTypeCodeGet
	MessageTypeCodeBack
	MessageTypeDataPush // 10
	MessageTypeDataAck
	MessageTypeDataGet
	MessageTypeDataBack
	MessageTypeTaskPush
	MessageTypeTaskAck
)

type EncryptType uint8

const (
	EncryptTypeNone EncryptType = iota
	EncryptTypeBase64
	EncryptTypeAesCbc
)
