package websocket

func parseBinaryPacket(pkt []byte) (msgName string, data []byte) {
	return "protobuf.PbMessage", pkt
}

func composeBinaryPacket(msgName string, data []byte) []byte {
	return data
}
