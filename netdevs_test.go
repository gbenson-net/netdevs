package netdevs

import (
	"net"
	"testing"

	"gotest.tools/v3/assert"
)

func TestKeyForIPv4(t *testing.T) {
	key, err := keyForIP(net.IPv4(97, 100, 100, 114))
	assert.NilError(t, err)
	assert.Equal(t, key, "i4addr")
}

func TestKeyForIPv6(t *testing.T) {
	ip := net.ParseIP("7468:6520:4269:6720:4c65:626f:7773:6b69")
	assert.Assert(t, ip != nil)
	key, err := keyForIP(ip)
	assert.NilError(t, err)
	assert.Equal(t, key, "i6the Big Lebowski")
}

func TestKeyForMAC(t *testing.T) {
	hw, err := net.ParseMAC("62:69:67:4d:61:63")
	assert.NilError(t, err)
	key, err := keyForMAC(hw)
	assert.NilError(t, err)
	assert.Equal(t, key, "m6bigMac")
}

func TestLookupKnownIP(t *testing.T) {
	ifi, err := InterfaceWithIP(net.IPv4(127, 0, 0, 1))
	assert.NilError(t, err)
	assert.Assert(t, ifi != nil)
	assert.Equal(t, ifi.Name, "lo")
}

func TestLookupUnknownIP(t *testing.T) {
	ifi, err := InterfaceWithIP(net.IPv4(127, 255, 255, 255))
	assert.NilError(t, err)
	assert.Assert(t, ifi == nil)
}

func TestUnsupportedIP(t *testing.T) {
	want := `unhandled IP address "hello"`
	ifi, err := InterfaceWithIP(net.IP([]byte("hello")))
	assert.Assert(t, ifi == nil)
	assert.Assert(t, err != nil)
	assert.Equal(t, err.Error(), want)
}

func TestLookupKnownMAC(t *testing.T) {
	ifs, err := net.Interfaces()
	assert.NilError(t, err)
	count := 0
	for _, want := range ifs {
		hw := want.HardwareAddr
		if len(hw) == 0 {
			continue
		}
		got, err := InterfaceWithMAC(hw)
		assert.NilError(t, err)
		assert.Equal(t, got.Name, want.Name)
		count++
	}
	assert.Assert(t, count > 0)
}

func TestUnsupportedMAC(t *testing.T) {
	want := `unhandled hardware address "hello"`
	ifi, err := InterfaceWithMAC(net.HardwareAddr([]byte("hello")))
	assert.Assert(t, ifi == nil)
	assert.Assert(t, err != nil)
	assert.Equal(t, err.Error(), want)
}
