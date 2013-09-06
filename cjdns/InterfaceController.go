package cjdns

import (
	"strings"
	"time"
)

type PeerState int

// Peer state values
const (
	Unauthenticated = iota
	Handshake
	Established
	Unresponsive
)

var (
	peerStateStrings = [4]string{
		"UNAUTHENTICATED",
		"HANDSHAKE",
		"ESTABLISHED",
		"UNRESPONSIVE",
	}
)

func (s PeerState) String() string {
	if s.Int() < 0 || s.Int() > len(peerStateStrings)-1 {
		return "INVALID"
	}
	return peerStateStrings[s]
}

func (s PeerState) Int() int {
	return int(s)
}

// Peer statistics
type PeerStats struct {
	PublicKey   string    // Public key of peer
	SwitchLabel string    // Internal switch label for reaching the peer
	IsIncoming  bool      // Is the peer connected to us, or us to them
	BytesOut    int64     // Total number of bytes sent
	BytesIn     int64     // Total number of bytes received
	State       PeerState // Peer connection state
	Last        time.Time // Last time a packet was received from the peer
}

// Returns stats on currently connected peers
func InterfaceController_peerStats(user *Conn, page int) (
	response []PeerStats, err error) {

	args := make(map[string]interface{})

	args["page"] = page

	data, err := SendCmd(user, "InterfaceController_peerStats", args)
	if err != nil {
		return
	}

	// Convert the map to a slice of structs.
	// This should be fixed so ALL functions return structs... eventually...
	response = make([]PeerStats, 0)
	for _, peer := range data["peers"].([]interface{}) {
		info := peer.(map[string]interface{})

		// Convert the int to a bool
		var incoming bool
		if info["isIncoming"].(int64) > 0 {
			incoming = true
		}

		// Convert connection state to an int
		var state PeerState
		tu := strings.ToUpper(info["state"].(string))
		for i, name := range peerStateStrings {
			if name == tu {
				state = PeerState(i)
			}
		}
		// Convert the last packet received timestamp to a time.Time
		last := time.Unix(info["last"].(int64), 0)

		peer := PeerStats{
			Last:        last,
			BytesIn:     info["bytesIn"].(int64),
			BytesOut:    info["bytesOut"].(int64),
			IsIncoming:  incoming,
			State:       state,
			PublicKey:   info["publicKey"].(string),
			SwitchLabel: info["switchLabel"].(string),
		}
		response = append(response, peer)
	}
	return
}

//InterfaceController_disconnectPeer(pubkey)