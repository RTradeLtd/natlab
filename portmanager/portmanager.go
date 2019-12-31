package portmanager

import (
	"math/rand"
	"net"
)

// PortManager allocates WAN UDP ports according to a configurable policy.
type PortManager struct {
	wanIPs []net.IP
	rng    *rand.Rand
}

// ipRefcount holds an IP address and a reference count.
type ipRefcount struct {
	ip     net.IP
	refcnt int
}

func New(wanIPs []net.IP) *PortManager {
	return &PortManager{
		wanIPs: append([]net.IP(nil), wanIPs...),
		rng:    NewRandom(),
	}
}

// Allocate tries to allocate a WAN ip:port for the given clientAddr.
func (p *PortManager) Allocate(clientAddr *net.UDPAddr) (port *net.UDPAddr, close func(), err error) {
	// TODO: more policies, right now we're just doing fully
	// randomized allocation.
	for {
		publicIP := p.wanIPs[p.rng.Intn(len(p.wanIPs))]
		port := 1024 + p.rng.Intn(64511)
		conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: publicIP, Port: port})
		if err != nil {
			// TODO: log?
			continue
		}

		addr := conn.LocalAddr().(*net.UDPAddr)
		return addr, func() { conn.Close() }, nil
	}
}