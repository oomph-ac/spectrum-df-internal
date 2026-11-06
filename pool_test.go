package spectrum

import (
	"testing"

	spectrumpacket "github.com/cooldogedev/spectrum-df/server/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// otherRequest stands in for Spectrum's own ConnectionRequest, which its packet package registers in gophertunnel's pool.
type otherRequest struct {
	spectrumpacket.ConnectionRequest
}

func TestClientPoolKeepsOwnPackets(t *testing.T) {
	packet.RegisterPacketFromClient(spectrumpacket.IDConnectionRequest, func() packet.Packet { return &otherRequest{} })
	t.Cleanup(func() {
		packet.RegisterPacketFromClient(spectrumpacket.IDConnectionRequest, func() packet.Packet { return &spectrumpacket.ConnectionRequest{} })
	})

	pool := clientPool()
	if _, ok := pool[spectrumpacket.IDConnectionRequest]().(*spectrumpacket.ConnectionRequest); !ok {
		t.Fatal("ConnectionRequest is not this module's type")
	}
	if _, ok := pool[spectrumpacket.IDLatency]().(*spectrumpacket.Latency); !ok {
		t.Fatal("Latency is not this module's type")
	}
}
