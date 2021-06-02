package network

import (
	"fmt"
	"net"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestNext(t *testing.T) {

	networks := make([]*net.IPNet, 0)
	for _, cidr := range []string{
		"172.18.1.0/24",
		"172.17.1.0/24",
		"172.18.2.0/24",
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
	_, myNext, err := net.ParseCIDR("172.18.3.0/24")
	assert.NoError(t, err)
	assert.Equal(t, myNext, next)
}
