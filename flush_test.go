package spectrum

import (
	"bytes"
	"errors"
	"net"
	"os"
	"testing"
	"time"

	spectrumprotocol "github.com/cooldogedev/spectrum-df/protocol"
	spectrumpacket "github.com/cooldogedev/spectrum-df/server/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// testConn returns a conn that writes to a pipe, and the pipe's other end, which plays the proxy.
func testConn(t *testing.T) (*conn, net.Conn) {
	server, proxy := net.Pipe()
	c := &conn{conn: server, writer: spectrumprotocol.NewWriter(server), pool: packet.NewClientPool(), closed: make(chan struct{})}
	t.Cleanup(func() {
		_ = c.Close()
		_ = proxy.Close()
	})
	return c, proxy
}

// readID reads one frame and returns its packet ID, or -1 when nothing comes within the timeout.
func readID(t *testing.T, r *spectrumprotocol.Reader, proxy net.Conn, timeout time.Duration) int {
	_ = proxy.SetReadDeadline(time.Now().Add(timeout))
	payload, err := r.ReadPacket()
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return -1
		}
		t.Fatalf("read frame: %v", err)
	}
	header := &packet.Header{}
	if err := header.Read(bytes.NewBuffer(payload[1:])); err != nil {
		t.Fatalf("read header: %v", err)
	}
	return int(header.PacketID)
}

func TestFlushFollowsWrittenPackets(t *testing.T) {
	c, proxy := testConn(t)
	r := spectrumprotocol.NewReader(proxy)
	go c.flushLoop()

	go func() { _ = c.WritePacket(&packet.SetTime{Time: 1}) }()
	if id := readID(t, r, proxy, time.Second); id != packet.IDSetTime {
		t.Fatalf("first frame is packet %d, want SetTime", id)
	}
	if id := readID(t, r, proxy, time.Second); id != spectrumpacket.IDFlush {
		t.Fatalf("second frame is packet %d, want Flush", id)
	}
	if id := readID(t, r, proxy, 3*flushInterval); id != -1 {
		t.Fatalf("got packet %d with nothing written, want no frame", id)
	}
}

func TestFlushSkipsSpectrumPackets(t *testing.T) {
	c, proxy := testConn(t)
	r := spectrumprotocol.NewReader(proxy)
	go c.flushLoop()

	go func() { _ = c.WritePacket(&spectrumpacket.Latency{}) }()
	if id := readID(t, r, proxy, time.Second); id != spectrumpacket.IDLatency {
		t.Fatalf("first frame is packet %d, want Latency", id)
	}
	if id := readID(t, r, proxy, 3*flushInterval); id != -1 {
		t.Fatalf("got packet %d after a Spectrum packet, want no Flush", id)
	}
}

func TestFlushWithNothingWritten(t *testing.T) {
	c, proxy := testConn(t)
	if err := c.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if id := readID(t, spectrumprotocol.NewReader(proxy), proxy, 2*flushInterval); id != -1 {
		t.Fatalf("got packet %d, want no frame", id)
	}
}

func TestFlushLoopStopsOnClose(t *testing.T) {
	c, _ := testConn(t)
	done := make(chan struct{})
	go func() {
		c.flushLoop()
		close(done)
	}()
	_ = c.Close()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("flushLoop still runs after Close")
	}
}

func TestDisconnectFlushes(t *testing.T) {
	c, proxy := testConn(t)
	r := spectrumprotocol.NewReader(proxy)
	done := make(chan struct{})
	go func() {
		_ = (&Listener{}).Disconnect(c, "bye")
		close(done)
	}()
	if id := readID(t, r, proxy, time.Second); id != packet.IDDisconnect {
		t.Fatalf("first frame is packet %d, want Disconnect", id)
	}
	if id := readID(t, r, proxy, time.Second); id != spectrumpacket.IDFlush {
		t.Fatalf("second frame is packet %d, want Flush", id)
	}
	<-done
}
