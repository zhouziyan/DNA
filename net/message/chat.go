package message

import (
	"bytes"
	"encoding/binary"

	"DNA/core/ledger"
	"DNA/events"
	. "DNA/net/protocol"
)

type chat struct {
	msgHdr
	content []byte
}

func NewChatMsg(msg []byte) ([]byte, error) {
	var c chat
	c.content = msg
	b := new(bytes.Buffer)
	if err := binary.Write(b, binary.LittleEndian, c.content); err != nil {
		return nil, err
	}
	s := checkSum(b.Bytes())
	c.init("chat", s, uint32(len(b.Bytes())))
	m, err := c.Serialization()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (msg chat) Verify(buf []byte) error {
	return msg.msgHdr.Verify(buf)
}

func (msg chat) Handle(node Noder) error {
	ledger.DefaultLedger.Blockchain.BCEvents.Notify(events.EventChatMessage, string(msg.content))
	return nil
}

func (msg chat) Serialization() ([]byte, error) {
	hdrBuf, err := msg.msgHdr.Serialization()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(hdrBuf)
	if err := binary.Write(buf, binary.LittleEndian, msg.content); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (msg *chat) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	if err := binary.Read(buf, binary.LittleEndian, &(msg.msgHdr)); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.LittleEndian, msg.content); err != nil {
		return err
	}

	return nil
}
