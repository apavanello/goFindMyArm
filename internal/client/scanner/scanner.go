package scanner

import (
	"fmt"
	"net"
	"time"

	"github.com/apavanello/goFindMyArm/internal/shared/crypto"
	"github.com/apavanello/goFindMyArm/pkg/protocol"
	"google.golang.org/protobuf/proto"
)

type Device struct {
	Hostname   string `json:"hostname"`
	IP         string `json:"ip"`
	MAC        string `json:"mac"`
	OS         string `json:"os"`
	Version    string `json:"version"`
	LastSeenAt string `json:"last_seen"` // formatted time
}

// ScanResults maps IP to Device
type ScanResults map[string]Device

// ScanNetwork broadcasts a probe and listens for responses for the given duration
func ScanNetwork(port int, timeout time.Duration, password string) (ScanResults, error) {
	results := make(ScanResults)

	// 1. Setup UDP Broadcast
	addr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: port,
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial udp: %v", err)
	}
	defer conn.Close()

	// 2. Send Probe
	probe := &protocol.DiscoveryPacket{
		MagicHeader: "GFMA_PROBE",
		Nonce:       []byte("TODO_RANDOM"), // In real crypto, nonce is part of the sealed box, here just a check
	}
	probeBytes, _ := proto.Marshal(probe)

	if _, err := conn.Write(probeBytes); err != nil {
		return nil, fmt.Errorf("failed to write udp: %v", err)
	}

	// 3. Listen for responses (we need a separate listener because DialUDP acts as connected socket depending on usage,
	// but for broadcast response we should listen on a port or use the same conn if it's bound correctly)
	// Actually, typically we bind to 0.0.0.0 to listen for responses.
	// Let's create a listener.

	listenAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0} // Ephemeral port
	listener, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return nil, err
	}
	defer listener.Close()

	// We need to re-send the packet with the Source Port of our listener so agents reply to US.
	// But Wait, broadcast is one-way usually unless we fake the source.
	// Easier Strategy: Agents reply to the SENDER's IP/Port.
	// So we should send FROM the listener connection.

	dstAddr := &net.UDPAddr{IP: net.IPv4bcast, Port: port}
	if _, err := listener.WriteToUDP(probeBytes, dstAddr); err != nil {
		return nil, err
	}

	fmt.Printf("[Scanner] Sending probe to %v from %v\n", dstAddr, listener.LocalAddr())

	deadline := time.Now().Add(timeout)
	listener.SetReadDeadline(deadline)

	buf := make([]byte, 4096)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(buf)
		if err != nil {
			// Timeout expected
			break
		}

		fmt.Printf("[Scanner] Received %d bytes from %v\n", n, remoteAddr)

		// 4. Decrypt
		decrypted, err := crypto.Decrypt(buf[:n], password)
		if err != nil {
			fmt.Printf("[Scanner] Decrypt failed from %s: %v\n", remoteAddr.String(), err)
			continue
		}

		// 5. Unmarshal
		var info protocol.DeviceInfo
		if err := proto.Unmarshal(decrypted, &info); err != nil {
			fmt.Printf("[Scanner] Unmarshal failed: %v\n", err)
			continue
		}

		results[info.IpAddress] = Device{
			Hostname:   info.Hostname,
			IP:         info.IpAddress,
			MAC:        info.MacAddress,
			OS:         info.OsDistro,
			Version:    info.AgentVersion,
			LastSeenAt: time.Now().Format(time.TimeOnly),
		}
	}

	return results, nil
}
