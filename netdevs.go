// Package netdevs maps hardware and IP addresses to local network interfaces.
package netdevs

import (
	"fmt"
	. "net"
	"time"
)

var (
	// ifmap maps hardware and IP addresses to local network interfaces.
	ifmap map[string]*Interface

	// MinScanInterval is the minimum time between ifmap updates.
	MinScanInterval = 100 * time.Millisecond

	// deadline is the time before which ifmap will not be updated.
	deadline time.Time
)

// InterfaceWithIP returns the local network interface configured with
// the specified IP address, or nil, nil if no interface is configured
// with the specified address.
func InterfaceWithIP(ip IP) (*Interface, error) {
	k, err := keyForIP(ip)
	if err != nil {
		return nil, err
	}
	return interfaceWithKey(k)
}

// InterfaceWithMAC returns the local network interface with the
// specified hardware address, or nil, nil if no interface has the
// specified address.
func InterfaceWithMAC(hw HardwareAddr) (*Interface, error) {
	k, err := keyForMAC(hw)
	if err != nil {
		return nil, err
	}
	return interfaceWithKey(k)
}

// keyForIP returns a unique key for the given IP address.
func keyForIP(ip IP) (string, error) {
	if ipv4 := ip.To4(); ipv4 != nil {
		return "i4" + string(ipv4[:IPv4len]), nil
	}
	if ipv6 := ip.To16(); ipv6 != nil {
		return "i6" + string(ipv6[:IPv6len]), nil
	}
	return "", fmt.Errorf("unhandled IP address %q", []byte(ip))
}

// keyForMAC returns a unique key for the given hardware address.
func keyForMAC(hw HardwareAddr) (string, error) {
	if len(hw) == 6 {
		return "m6" + string(hw), nil
	}
	return "", fmt.Errorf("unhandled hardware address %q", []byte(hw))
}

// interfaceWithKey returns the local network interface with the
// specified key, or nil, nil if no interface has the specified key.
func interfaceWithKey(k string) (*Interface, error) {
	if ifi, found := ifmap[k]; found {
		return ifi, nil
	}

	if err := maybeUpdateIfMap(); err != nil {
		return nil, err
	}

	return ifmap[k], nil
}

// maybeUpdateIfMap replaces ifmap with an updated table if the current
// time is after deadline.  Does nothing if called before deadline.
func maybeUpdateIfMap() error {
	ts := time.Now()
	if ts.Before(deadline) {
		return nil
	}

	ifs, err := Interfaces()
	if err != nil {
		return err
	}

	ifmap = make(map[string]*Interface)
	for _, ifi := range ifs {
		if hw := ifi.HardwareAddr; len(hw) != 0 {
			k, err := keyForMAC(hw)
			if err != nil {
				return err
			}
			if _, found := ifmap[k]; found {
				return fmt.Errorf("duplicate hardware address %q", hw)
			}
			ifmap[k] = &ifi
		}

		addrs, err := ifi.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			netaddr, ok := addr.(*IPNet)
			if !ok {
				continue
			}

			ip := netaddr.IP
			k, err := keyForIP(ip)
			if err != nil {
				return err
			}
			if _, found := ifmap[k]; found {
				return fmt.Errorf("duplicate IP address %q", ip)
			}
			ifmap[k] = &ifi
		}
	}

	deadline = ts.Add(MinScanInterval)
	return nil
}
