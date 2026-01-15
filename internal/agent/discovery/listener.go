package discovery

import (
	"fmt"
	"net"
	"os"
	"runtime"

	"github.com/apavanello/goFindMyArm/internal/shared/crypto"
	"github.com/apavanello/goFindMyArm/pkg/protocol"
	"google.golang.org/protobuf/proto"
)

const MagicHeader = "GFMA_PROBE"

// Serve starts the UDP listener
func Serve(port int, password string, disableRemote bool, agentVersion string) error {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("[Discovery] Listening on UDP %d\n", port)

	buf := make([]byte, 2048)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("[Discovery] Error reading UDP: %v\n", err)
			continue
		}

		go handlePacket(conn, remoteAddr, buf[:n], password, disableRemote, agentVersion)
	}
}

func handlePacket(conn *net.UDPConn, remoteAddr *net.UDPAddr, data []byte, password string, disableRemote bool, version string) {
	// 1. Unmarshal Request
	var req protocol.DiscoveryPacket
	if err := proto.Unmarshal(data, &req); err != nil {
		// Ignore invalid packets to avoid log spam from random network noise
		return
	}

	if req.MagicHeader != MagicHeader {
		return
	}

	// 2. Build Response
	hostname, _ := os.Hostname()
	ip := getOutboundIP() // Helper to find real IP

	resp := &protocol.DeviceInfo{
		Hostname:              hostname,
		IpAddress:             ip.String(),
		MacAddress:            getMacAddr(),
		OsDistro:              runtime.GOOS,   // Simple for now
		KernelVersion:         runtime.GOARCH, // Simple for now
		AgentVersion:          version,
		RemoteCommandsEnabled: !disableRemote,
	}

	respBytes, err := proto.Marshal(resp)
	if err != nil {
		fmt.Printf("[Discovery] Error marshaling response: %v\n", err)
		return
	}

	// 3. Encrypt Response
	encryptedResp, err := crypto.Encrypt(respBytes, password)
	if err != nil {
		fmt.Printf("[Discovery] Error encrypting response: %v\n", err)
		return
	}

	// 4. Send Response
	_, err = conn.WriteToUDP(encryptedResp, remoteAddr)
	if err != nil {
		fmt.Printf("[Discovery] Error sending response: %v\n", err)
	} else {
		fmt.Printf("[Discovery] Responded to %s\n", remoteAddr.String())
	}
}

// getOutboundIP prefers UDP connection to find preferred outbound IP
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.ParseIP("127.0.0.1")
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func getMacAddr() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && i.HardwareAddr != nil {
			// Naive: returns the first valid MAC
			return i.HardwareAddr.String()
		}
	}
	return ""
}
