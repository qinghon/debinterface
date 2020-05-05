package debinterface

import (
	"io/ioutil"
	"testing"
)

func TestInterfaceReader_ParseAuto(t *testing.T) {
	var reader InterfaceReader
	reader.context = 0
	reader.Adapters = []Interface{0: Interface{"test1": "set test1"}}
	reader.ParseAuto()
	t.Log(reader.Adapters)
}
func TestInterfaceReader_ReadLines(t *testing.T) {
	ioutil.WriteFile("/tmp/testInterface1", []byte(`
# madwifi-ng WDS Bridge
source /etc/network/interfaces.d/*
no-auto-down br0
auto br0
iface br0 inet static
       address 192.168.1.2/24
       netmask 255.255.255.0
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
       up iwpriv ath1 wds_add AA:BB:CC:DD:EE:FF
       post-up ifconfig ath1 down ; ifconfig ath1 up # this is a workaround
       post-down wlanconfig ath0 destroy
       pre-down wlanconfig ath1 destroy
iface br0 inet6 auto
	dhcp 1
	request_prefix 1
no-auto-down eth0
allow-auto eth0
no-scripts eth0
allow-hotplug eth0
iface eth0 inet dhcp
`), 0644)
	var reader = NewReader("/tmp/testInterface1")

	err := reader.ReadLines("/tmp/testInterface1")
	if err != nil {
		t.Error(err)
	}
	//js,_:=json.Marshal(reader.Adapters)
	//t.Log(string(js))
	t.Log(reader.Adapters)
	for _, v := range reader.Adapters {
		t.Log(v.Export())
	}
}
func TestInterfaceReader_ReadLines_2(t *testing.T) {
	ioutil.WriteFile("/tmp/testInterface2", []byte(`
no-auto-dow tap0
auto tap0
iface tap0 inet static
       address 192.168.1.2
       netmask 255.255.255.0
       network 192.168.1.0
       broadcast 192.168.1.255
`), 0644)
	var reader = NewReader("/tmp/testInterface2")

	err := reader.ReadLines("/tmp/testInterface2")
	if err != nil {
		t.Error(err)
	}
	//js,_:=json.Marshal(reader.Adapters)
	//t.Log(string(js))
	t.Log(reader.Adapters)
	for _, v := range reader.Adapters {
		t.Log(v.Export())
	}
}
func BenchmarkInterfaceReader_ReadLines(b *testing.B) {
	ioutil.WriteFile("/tmp/testInterface1", []byte(`
# madwifi-ng WDS Bridge
#source /etc/network/interfaces.d/*
auto br0
iface br0 inet static
       address 192.168.1.2/24
       netmask 255.255.255.0
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
       up iwpriv ath1 wds_add AA:BB:CC:DD:EE:FF
       post-up ifconfig ath1 down ; ifconfig ath1 up # this is a workaround
       post-down wlanconfig ath0 destroy
       pre-down wlanconfig ath1 destroy
iface br0 inet6 auto
	dhcp 1
	request_prefix 1

`), 0644)
	var reader = NewReader("/tmp/testInterface1")
	b.ResetTimer() // HL
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		// 被测试的功能
		err := reader.ReadLines("/tmp/testInterface1")
		if err != nil {
			b.Error(err)
		}
	}
}
