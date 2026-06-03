package protocol

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

func ReadPacket(conn net.Conn) ([]byte, error) {
	var buf [4]byte
	_, err := io.ReadFull(conn, buf[:])
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(buf[:])
	if length > 1<<20 {
		return nil, fmt.Errorf("packet too large")
	}
	payload := make([]byte, length)
	_, err = io.ReadFull(conn, payload)
	return payload, err
}
func SendPacket(conn net.Conn, data interface{}) error {
	jsonBytes, _ := json.Marshal(data)
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(len(jsonBytes)))
	_, err := conn.Write(buf[:])
	if err != nil {
		return err
	}
	_, err = conn.Write(jsonBytes)
	return err
}
