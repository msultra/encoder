package encoder_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/msultra/encoder"
)

type MessageHeader struct {
	Signature   [8]byte `smb:"len:8"`
	MessageType uint32
}

type ChallengeMessage struct {
	MessageHeader
	TargetName struct {
		Lenght uint16
		MaxLen uint16
		Offset uint32
	}
	NegotiateFlags    uint32
	ServerChallenge   [8]byte
	Reserved          [8]byte
	TargetInformation struct {
		Lenght uint16
		MaxLen uint16
		Offset uint32
	}
	Version [8]byte
	Payload []byte
}

// TestUnmarshal tests the unmarshal function
func TestUnmarshal(t *testing.T) {
	data := "4e544c4d53535000020000000600060038000000358289e23ce65bea9b2dc11000000000000000005e005e003e0000000a0063450000000f4c0041004200020006004c0041004200010004004400430004000e006c00610062002e006c0061006e0003001400440043002e006c00610062002e006c0061006e0005000e006c00610062002e006c0061006e000700080084aa2fbb90ecd80100000000"
	bs, err := hex.DecodeString(data)
	if err != nil {
		t.Fatal(err)
	}

	var msg ChallengeMessage
	if err := encoder.Unmarshal(bs, &msg); err != nil {
		t.Fatal(err)
	}

	targetName := encoder.UTF16ToStr(msg.Payload[msg.TargetName.Offset-uint32(56) : msg.TargetName.Offset-uint32(56)+uint32(msg.TargetName.Lenght)])
	if targetName != "LAB" {
		t.Fatalf("targetName mismatch: %v != LAB", targetName)
	}

	if msg.NegotiateFlags != 0xe2898235 {
		t.Fatalf("NegotiateFlags mismatch: %v != 0xe2898235", msg.NegotiateFlags)
	}

	if msg.ServerChallenge != [8]byte{0x3c, 0xe6, 0x5b, 0xea, 0x9b, 0x2d, 0xc1, 0x10} {
		t.Fatalf("ServerChallenge mismatch: %v != [0x3c, 0xe6, 0x5b, 0xea, 0x9b, 0x2d, 0xc1, 0x10]", msg.ServerChallenge)
	}

	targetInfoBytes := msg.Payload[msg.TargetInformation.Offset-uint32(56) : msg.TargetInformation.Offset-uint32(56)+uint32(msg.TargetInformation.Lenght)]
	expectedTgtInfo := []byte{0x2, 0x0, 0x6, 0x0, 0x4c, 0x0, 0x41, 0x0, 0x42, 0x0, 0x1, 0x0, 0x4, 0x0, 0x44, 0x0, 0x43, 0x0, 0x4, 0x0, 0xe, 0x0, 0x6c, 0x0, 0x61, 0x0, 0x62, 0x0, 0x2e, 0x0, 0x6c, 0x0, 0x61, 0x0, 0x6e, 0x0, 0x3, 0x0, 0x14, 0x0, 0x44, 0x0, 0x43, 0x0, 0x2e, 0x0, 0x6c, 0x0, 0x61, 0x0, 0x62, 0x0, 0x2e, 0x0, 0x6c, 0x0, 0x61, 0x0, 0x6e, 0x0, 0x5, 0x0, 0xe, 0x0, 0x6c, 0x0, 0x61, 0x0, 0x62, 0x0, 0x2e, 0x0, 0x6c, 0x0, 0x61, 0x0, 0x6e, 0x0, 0x7, 0x0, 0x8, 0x0, 0x84, 0xaa, 0x2f, 0xbb, 0x90, 0xec, 0xd8, 0x1, 0x0, 0x0, 0x0, 0x0}
	if slices.Compare(targetInfoBytes, expectedTgtInfo) != 0 {
		t.Fatalf("targetInfoBytes mismatch: %v != %v", targetInfoBytes, expectedTgtInfo)
	}
}

func TestMarshal(t *testing.T) {
	msg := ChallengeMessage{
		MessageHeader: MessageHeader{
			MessageType: 2,
			Signature:   [8]byte{0x4e, 0x54, 0x4c, 0x4d, 0x53, 0x53, 0x50, 0x0},
		},
		NegotiateFlags:  0xe2898235,
		ServerChallenge: [8]byte{0x3c, 0xe6, 0x5b, 0xea, 0x9b, 0x2d, 0xc1, 0x10},
	}

	bs, err := encoder.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.EncodeToString(bs))
}
