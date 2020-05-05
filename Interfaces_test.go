package debinterface

import (
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"testing"
)

func TestMktmp1(t *testing.T) {
	target, err := Mktmp(false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(target)
	cmd := exec.Command("ls", "/tmp/tmp.*")
	t.Log(cmd.Output())
	defer os.Remove(target)
}

func TestMktmp2(t *testing.T) {
	target, err := Mktmp(true)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(target)
	defer os.Remove(target)
}
func TestGenTestFile(t *testing.T) {
	ioutil.WriteFile("/tmp/testInterface1", []byte(`# madwifi-ng WDS Bridge
#source /etc/network/interfaces.d/*
auto lo
iface lo inet lookback

auto br0
hotplug br0
iface br0 inet static
       address 192.168.1.2/24
       #netmask 255.255.255.0
       network 192.168.1.0
       broadcast 192.168.1.255
       gateway 192.168.1.1
       bridge_ports eth0 ath0 ath1
		bridge_waitport 0    # no delay before a port becomes available
        bridge_fd 0          # no forwarding delay
        bridge_ports none    # if you do not want to bind to any ports
        bridge_ports regex eth* # use a regular expression to define ports
       pre-up wlanconfig ath0 create wlandev wifi0 wlanmode ap
       pre-up wlanconfig ath1 create wlandev wifi0 wlanmode wds
       pre-up iwpriv ath0 mode 11g
       pre-up iwpriv ath0 bintval 1000
       pre-up iwconfig ath0 essid "voyage-wds" channel 1
       up ifconfig ath0 down ; ifconfig ath0 up # this is a workaround
       up iwpriv ath1 wds 1
       #up iwpriv ath1 wds_add AA:BB:CC:DD:EE:FF
       post-up ifconfig ath1 down ; ifconfig ath1 up # this is a workaround
       post-down wlanconfig ath0 destroy
       pre-down wlanconfig ath1 destroy
		# anscoasnc
# 1test
# asncoaisncoasnmcomas

iface br1 inet6 auto
	dhcp 1
	request_prefix 1
`), 0644)
}
func TestInterfaces_Del(t *testing.T) {
	var adapter = Interface{}
	adapter.SetName("br0")
	adapter.SetAuto(true)
	adapter.SetAddrSource("static")
	adapter.SetAddrFam("inet")
	adapter.SetAddress(net.ParseIP("192.168.4.2"))
	adapter.SetGateWay(net.ParseIP("192.168.4.1"))
	adapter.SetMtu(60)
	adapter.SetUnkonw("server 192.168.4.1")
	TestGenTestFile(t)
	var adapters = Interfaces{}
	adapters.FilePath = "/tmp/testInterface1"
	//adapters.Adapters=append(adapters.Adapters,adapter)
	err := adapters.Del(adapter)
	if err != nil {
		t.Error(err)
	}
}
func TestInterfaces_Update(t *testing.T) {
	var adapter = Interface{}
	adapter.SetName("br0")
	adapter.SetAuto(true)
	adapter.SetAddrSource("static")
	adapter.SetAddrFam("inet")
	adapter.SetAddress(net.ParseIP("192.168.4.2"))
	adapter.SetGateWay(net.ParseIP("192.168.4.1"))
	adapter.SetMtu(60)
	adapter.SetUnkonw("server 192.168.4.1")
	TestGenTestFile(t)
	var adapters = Interfaces{}
	adapters.FilePath = "/tmp/testInterface1"
	//adapters.Adapters=append(adapters.Adapters,adapter)
	err := adapters.Update(adapter)
	if err != nil {
		t.Error(err)
	}
}
func TestInterfaces_Add(t *testing.T) {
	var adapter = Interface{}
	adapter.SetName("br2")
	adapter.SetAuto(true)
	adapter.SetAddrSource("static")
	adapter.SetAddrFam("inet")
	adapter.SetAddress(net.ParseIP("192.168.4.2"))
	adapter.SetGateWay(net.ParseIP("192.168.4.1"))
	adapter.SetMtu(60)
	adapter.SetUnkonw("server 192.168.4.1")
	TestGenTestFile(t)
	var adapters = Interfaces{}
	adapters.FilePath = "/tmp/testInterface1"
	//adapters.Adapters=append(adapters.Adapters,adapter)
	err := adapters.Add(adapter)
	if err != nil {
		t.Error(err)
	}
}
