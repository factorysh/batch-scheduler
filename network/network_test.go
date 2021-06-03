package network

import (
	"fmt"
	"net"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPow(t *testing.T) {
	assert.Equal(t, 1, intPow(2, 0))
	assert.Equal(t, 0, intPow(0, 8))
	assert.Equal(t, 256, intPow(2, 8))
}
func TestByNetwork(t *testing.T) {
	networks := make(ByNetwork, 0)
	for _, cidr := range []string{
		"172.18.1.0/24",
		"172.17.1.0/24",
		"172.18.2.0/24",
	} {
		_, n, err := net.ParseCIDR(cidr)
		assert.NoError(t, err)
		networks = append(networks, n)
	}
	sort.Sort(networks)
	fmt.Println(networks)
	txt := make([]string, len(networks))
	for i, n := range networks {
		txt[i] = n.String()
	}
	assert.Equal(t, []string{"172.17.1.0/24", "172.18.1.0/24", "172.18.2.0/24"}, txt)
}

func TestLast(t *testing.T) {
	_, n, err := net.ParseCIDR("172.18.1.0/24")
	assert.NoError(t, err)
	first, last := FirstLast(n)
	assert.Equal(t, net.IP{172, 18, 1, 0}, first)
	assert.Equal(t, net.IP{172, 18, 1, 255}, last)
}

func TestDistance(t *testing.T) {
	d := Distance(net.IP{172, 18, 1, 0}, net.IP{172, 18, 1, 255})
	assert.Equal(t, 255, d)
	d = Distance(net.IP{172, 17, 0, 42}, net.IP{172, 17, 0, 0})
	assert.Equal(t, -42, d)
	d = Distance(net.IP{172, 18, 1, 0}, net.IP{172, 18, 1, 0})
	assert.Equal(t, 0, d)
}

func TestNetDistance(t *testing.T) {
	_, a, err := net.ParseCIDR("172.17.0.0/24")
	assert.NoError(t, err)
	_, b, err := net.ParseCIDR("172.17.1.0/24")
	assert.NoError(t, err)
	d := NetDistance(a, b)
	assert.Equal(t, 256, d)
	d = NetDistance(b, a)
	assert.Equal(t, 256, d)
	d = NetDistance(a, a)
	assert.Equal(t, 0, d)
	_, c, err := net.ParseCIDR("172.17.0.0/16")
	assert.NoError(t, err)
	d = NetDistance(a, c)
	assert.True(t, d < 0)
}

func TestNext(t *testing.T) {
	_, min, err := net.ParseCIDR("172.17.0.0/24")
	assert.NoError(t, err)
	_, max, err := net.ParseCIDR("172.18.32.0/24")
	assert.NoError(t, err)
	networks := make([]*net.IPNet, 0)

	next, err := NextAvailableNetwork(networks, min, max, net.IPMask{255, 255, 255, 0})
	assert.NoError(t, err)
	assert.Equal(t, min, next)

	for _, cidr := range []string{
		"172.18.1.0/24",
		"172.17.1.0/24",
		"172.18.2.0/24",
	} {
		_, n, err := net.ParseCIDR(cidr)
		assert.NoError(t, err)
		networks = append(networks, n)
	}
	next, err = NextAvailableNetwork(networks, min, max, net.IPMask{255, 255, 255, 0})
	assert.NoError(t, err)
	_, myNext, err := net.ParseCIDR("172.18.3.0/24")
	assert.NoError(t, err)
	assert.Equal(t, myNext, next)
	networks = append(networks, next)

	for i := 0; i < 20; i++ {
		next, err = NextAvailableNetwork(networks, min, max, net.IPMask{255, 255, 255, 0})
		assert.NoError(t, err)
		networks = append(networks, next)
	}

	bag := make(map[string]interface{})
	for _, n := range networks {
		bag[n.String()] = new(interface{})
	}

	fmt.Println(bag)
	for i := 4; i < 24; i++ {
		k := fmt.Sprintf("172.18.%d.0/24", i)
		_, ok := bag[k]
		assert.True(t, ok, "Oups %s", k)
	}
}

func TestNextFirst(t *testing.T) {
	networks := make([]*net.IPNet, 0)
	for _, cidr := range []string{
		"172.18.1.0/24",
		"172.17.1.0/24",
		"172.18.32.0/24",
	} {
		_, n, err := net.ParseCIDR(cidr)
		assert.NoError(t, err)
		networks = append(networks, n)
	}
	_, min, err := net.ParseCIDR("172.17.0.0/24")
	assert.NoError(t, err)
	_, max, err := net.ParseCIDR("172.18.32.0/24")
	assert.NoError(t, err)
	next, err := NextAvailableNetwork(networks, min, max, net.IPMask{255, 255, 255, 0})
	assert.NoError(t, err)
	_, myNext, err := net.ParseCIDR("172.17.0.0/24")
	assert.NoError(t, err)
	assert.Equal(t, myNext, next)
}

func TestNextHole(t *testing.T) {
	networks := make([]*net.IPNet, 0)
	for _, cidr := range []string{
		"172.18.1.0/24",
		"172.17.0.0/24",
		"172.18.32.0/24", // last slot is used
	} {
		_, n, err := net.ParseCIDR(cidr)
		assert.NoError(t, err)
		networks = append(networks, n)
	}
	_, min, err := net.ParseCIDR("172.17.0.0/24")
	assert.NoError(t, err)
	_, max, err := net.ParseCIDR("172.18.32.0/24")
	assert.NoError(t, err)

	next, err := NextAvailableNetwork(networks, min, max, net.IPMask{255, 255, 255, 0})
	assert.NoError(t, err)
	_, myNext, err := net.ParseCIDR("172.17.1.0/24")
	assert.NoError(t, err)
	assert.Equal(t, myNext, next)

	networks = append(networks, next)

	next, err = NextAvailableNetwork(networks, min, max, net.IPMask{255, 255, 255, 0})
	assert.NoError(t, err)
	_, myNext, err = net.ParseCIDR("172.17.2.0/24")
	assert.NoError(t, err)
	assert.Equal(t, myNext, next)
}
