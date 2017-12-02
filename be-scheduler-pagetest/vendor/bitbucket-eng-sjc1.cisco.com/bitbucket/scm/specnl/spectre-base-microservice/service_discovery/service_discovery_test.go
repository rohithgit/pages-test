package service_discovery

import (
	"errors"
	"strconv"
	"testing"

	"github.com/benschw/dns-clb-go/dns"
	"github.com/stretchr/testify/assert"
)

const (
	name string = "search.service.consul"
)

type MockCLB struct {
	Address string
	Port    uint16
	Err     error
}

func (ref *MockCLB) GetAddress(name string) (dns.Address, error) {

	if ref.Err == nil {
		addr := dns.Address{
			ref.Address,
			ref.Port,
		}
		return addr, nil
	}
	return dns.Address{}, ref.Err
}

func NewMockCLB(addr string, port uint16, err error) *MockCLB {
	return &MockCLB{
		addr,
		port,
		err,
	}
}

func TestGetServiceAddressValid(t *testing.T) {
	mockclb := NewMockCLB("127.0.0.1", 8500, nil)
	saddr := NewServiceAddress(mockclb)
	err := saddr.GetServiceAddress(name)
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", saddr.IP)
	assert.Equal(t, uint16(8500), saddr.Port)
	assert.Equal(t, "127.0.0.1:8500", saddr.Address)
	assert.Equal(t, name+":"+strconv.Itoa(int(saddr.Port)), saddr.NameAndPort)

	mockclb = NewMockCLB("127.0.0.1", 8500, errors.New("it errored"))
	saddr = NewServiceAddress(mockclb)
	err = saddr.GetServiceAddress(name)
	assert.Error(t, err)
}

func TestParseServerAddressValidWithIP(t *testing.T) {
	mockclb := NewMockCLB("", 0, nil)
	saddr := NewServiceAddress(mockclb)
	err := saddr.ParseServerAddress("127.0.0.1:8500")
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", saddr.IP)
	assert.Equal(t, uint16(8500), saddr.Port)
	assert.Equal(t, "127.0.0.1:8500", saddr.Address)
}

func TestParseServerAddressValidWithoutPort(t *testing.T) {
	mockclb := NewMockCLB("", 0, nil)
	saddr := NewServiceAddress(mockclb)
	err := saddr.ParseServerAddress("127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", saddr.IP)
	assert.Equal(t, uint16(0), saddr.Port)
	assert.Equal(t, "127.0.0.1", saddr.Address)
}

func TestParseServerAddressValidWithName(t *testing.T) {
	mockclb := NewMockCLB("127.0.0.1", 8500, nil)
	saddr := NewServiceAddress(mockclb)
	err := saddr.ParseServerAddress(name)
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", saddr.IP)
	assert.Equal(t, uint16(8500), saddr.Port)
	assert.Equal(t, "127.0.0.1:8500", saddr.Address)
	assert.Equal(t, name, saddr.Name)
}
