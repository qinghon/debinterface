package debinterface

import (
	"net"
	"testing"
)

func TestInterface_String(t *testing.T) {
	var adapter =Interface{}
	adapter.SetName("eth0")
	adapter.SetAuto(true)
	adapter.SetAddrSource("static")
	adapter.SetAddrFam("inet")
	adapter.SetAddress(net.ParseIP("192.168.4.2"))
	adapter.SetGateWay(net.ParseIP("192.168.4.1"))
	adapter.SetMtu(60)
	t.Log(adapter.Export())
}
