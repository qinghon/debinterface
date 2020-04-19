package debinterface

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Interface map[string]interface{}

func (adapter Interface) Set(key string, v interface{}) {
	adapter[key] = v
}
func (adapter Interface) SetAuto(f bool) {
	adapter.Set("auto", f)
}
func (adapter Interface) SetHotplug(f bool) {
	adapter.Set("allow-hotplug", f)
}
func (adapter Interface) SetName(name string) {
	adapter.Set("name", name)
}
func (adapter Interface) SetAddrFam(addrFam string) {
	switch addrFam {
	case "inet", "inet6", "ipx", "can":
		adapter.Set("addrFam", addrFam)
	}
}
func (adapter Interface) SetAddrSource(methodName string) {
	switch methodName {
	case "dhcp", "static", "loopback", "manual",
		"bootp", "ppp", "wvdial", "dynamic",
		"ipv4ll", "v4tunnel", "auto", "6to4", "tunnel":
		adapter.Set("method_name", methodName)
	}
}

// The loopback Method
func (adapter Interface) SetLookBack(turn bool) {
	if turn {
		adapter.Set("lookback", "on")
	} else {
		adapter.Set("lookback", "off")
	}
}
// "static" Method
func (adapter Interface) SetAddress(address net.IP) {
	adapter.Set("address", address)
}
func (adapter Interface) SetNetmask(mask net.IPMask) {
	adapter.Set("netmask", mask)
}

func (adapter Interface) SetBroadcast(broadcast net.IP) {
	adapter.Set("broadcast", broadcast)
}
// Metric for added routes (dhclient)
func (adapter Interface) SetMetric(metric int) {
	adapter.Set("metric", metric)
}
func (adapter Interface) SetGateWay(ip net.IP) {
	adapter.Set("gateway", ip)
}

func (adapter Interface) SetPoinToPoint(ip net.IP) error {
	adapter.Set("pointopoint", ip)
	return nil
}

// "manual" "static"  "dhcp" Method
// Hardware address.
func (adapter Interface) SetHwaddress(hwaddress net.HardwareAddr) error {
	adapter.Set("hwaddress", hwaddress)
	return nil
}

// "manual" "static" "dhcp" Method
func (adapter Interface) SetMtu(mtu int) {
	adapter.Set("mtu", mtu)
}

func (adapter Interface) SetScope(scope string) error {
	switch scope {
	case "global", "link", "host":
		adapter["scope"] = scope
	default:
		return os.ErrInvalid
	}
	return nil
}
// "bootp" Method
func (adapter Interface) SetServer(server net.IP) {
	adapter.Set("server", server)
}
func (adapter Interface) SetDstAddr(dstaddr net.IP) {
	adapter.Set("dstaddr", dstaddr)
}
func (adapter Interface) SetLocal(local net.IP) {
	adapter.Set("local", local)
}
func (adapter Interface) SetDnsNameServer(ipList []net.IP) {
	adapter.Set("dns-nameserver", ipList)
}
func (adapter Interface) SetDnsSearch(ipList []string) {
	adapter.Set("dns-search", ipList)
}

// IFACE OPTIONS
// The following "command" options are available for every family and method.  Each of these options can be given multiple times in a  single  stanza,
// in  which  case the commands are executed in the order in which they appear in the stanza.  (You can ensure a command never fails by suffixing them
// with "|| true".)
func (adapter Interface) SetScript(action string, script string) {
	switch action {
	case "up":
		adapter.SetUpScript(script)
	case "down":
		adapter.SetDownScript(script)
	case "pre-up":
		adapter.SetPreUpScript(script)
	case "pre-down":
		adapter.SetPreDownScript(script)
	case "post-up":
		adapter.SetPostUpScript(script)
	case "post-down":
		adapter.SetPostDownScript(script)
	}
}

// Run command before bringing the interface up.  If this command fails then ifup aborts, refraining from marking the interface as  configured,
// prints an error message, and exits with status 0.  This behavior may change in the future.
func (adapter Interface) SetPreUpScript(up string) {
	if adapter["pre-up"] == nil {
		adapter["pre-up"] = []string{}
	}
	adapter["pre-up"] = append(adapter["pre-up"].([]string), up)
}
func (adapter Interface) SetUpScript(up string) {
	if adapter["up"] == nil {
		adapter["up"] = []string{}
	}
	adapter["up"] = append(adapter["up"].([]string), up)
}
// Run  command  after  bringing the interface up.  If this command fails then ifup aborts, refraining from marking the interface as configured
// (even though it has really been configured), prints an error message, and exits with status 0.  This behavior may change in the future.
func (adapter Interface) SetPostUpScript(up string) {
	if adapter["post-up"] == nil {
		adapter["post-up"] = []string{}
	}
	adapter["post-up"] = append(adapter["post-up"].([]string), up)
}
func (adapter Interface) SetDownScript(up string) {
	if adapter["down"] == nil {
		adapter["down"] = []string{}
	}
	adapter["down"] = append(adapter["down"].([]string), up)
}
func (adapter Interface) SetPreDownScript(up string) {
	if adapter["pre-down"] == nil {
		adapter["pre-down"] = []string{}
	}
	adapter["pre-down"] = append(adapter["pre-down"].([]string), up)
}
func (adapter Interface) SetPostDownScript(up string) {
	if adapter["post-down"] == nil {
		adapter["post-down"] = []string{}
	}
	adapter["post-down"] = append(adapter["post-down"].([]string), up)
}

// "dhcp" Method (pump, dhcpcd, udhcpc)
// Hostname to be requested (pump, dhcpcd, udhcpc)
func (adapter Interface) SetHostName(hostname string) {
	adapter.Set("hostname", hostname)
}

// Preferred lease time in hours (pump)
func (adapter Interface) SetLeaseHours(leasehours string) {
	adapter.Set("leasehours", leasehours)
}

// Preferred lease time in seconds (dhcpcd)
func (adapter Interface) SetLeaseTime(leasetime string) {
	adapter.Set("leasetime", leasetime)
}

// Vendor class identifier (dhcpcd)
func (adapter Interface) SetVendor(vendor string) {
	adapter.Set("vendor", vendor)
}

// Client identifier (dhcpcd)
func (adapter Interface) SetClient(client string) {
	adapter.Set("client", client)
}

// The bootp Method
func (adapter Interface) SetBootFile(bootfile string) {
	adapter.Set("bootfile", bootfile)
}

// The ppp Method
// Use name as the provider (from /etc/ppp/peers).
func (adapter Interface) SetProvider(provider string) {
	adapter.Set("provider", provider)
}

// Use number as the ppp unit number.
func (adapter Interface) SetUnit(unit string) {
	adapter.Set("unit", unit)
}

// Pass string as additional options to pon.
func (adapter Interface) SetOptions(options string) {
	adapter.Set("options", options)
}

// parsed from file ,only this package
func (adapter Interface) SetFromFile(_filepath string) {
	adapter["fromfile"] = _filepath
}

func (adapter Interface) SetUnkonw(unkonw string) {
	if adapter["unknow"] == nil {
		adapter["unknow"] = []string{}
	}
	adapter["unknow"] = append(adapter["unknow"].([]string), unkonw)
}

func (adapter Interface) GetName() string {
	if adapter["name"]==nil {
		return "<nil>"
	}
	return adapter["name"].(string)
}
// convert Interface to debian network file format
func (adapter Interface) Export() (string)  {
	var output string

	if adapter["auto"]==true {
		output+=fmt.Sprintf("auto %s\n",adapter["name"])
	}
	if adapter["hotplug"]==true {
		output+=fmt.Sprintf("hotplug %s\n",adapter["name"])
	}
	output+=fmt.Sprintf("iface %s %s %s\n",adapter["name"],adapter["addrFam"],adapter["method_name"])
	for k,v:=range adapter {
		if k == "method_name" || k == "addrFam" || k == "name"||k=="fromfile" {
			continue
		}
		switch v.(type) {
		case string:
			output+=fmt.Sprintf("\t%s %s\n",k,v)
		case net.IP:
			output+=fmt.Sprintf("\t%s %s\n",k,v.(net.IP).String())
		case net.IPMask:
			output+=fmt.Sprintf("\t%s %s\n",k,ToIPv4(v.(net.IPMask)))
		case []string:
			switch k {
			case "pre-up", "up", "post-up", "down", "pre-down", "post-down":
				for _,str:=range v.([]string) {
					output+=fmt.Sprintf("\t%s %s\n",k,str)
				}
			case "unknow":
				for _,str:=range v.([]string) {
					output+=fmt.Sprintf("\t%s\n",str)
				}
			default:
				output+=fmt.Sprintf("\t%s %s\n",k,strings.Join(v.([]string),""))
			}
		case []net.IP:
			ipList:=[]string{}
			for _,ip:=range v.([]net.IP) {
				ipList=append(ipList,ip.String())
			}
			output+=fmt.Sprintf("\t%s %s\n",k,strings.Join(ipList," "))
		case int:
			output+=fmt.Sprintf("\t%s %d\n",k,v)

		}
	}
	return output
}
// Convert ipv4 mask to ip format, skip ipv6 convert
func ToIPv4(m net.IPMask) string {
	p := m[len(m)-4:]

	if len(m) == 0 {
		return "<nil>"
	}
	// If IPv4, use dotted notation.
	if len(p) == net.IPv4len {
		const maxIPv4StringLen = len("255.255.255.255")
		b := make([]byte, maxIPv4StringLen)

		n := ubtoa(b, 0, p[0])
		b[n] = '.'
		n++

		n += ubtoa(b, n, p[1])
		b[n] = '.'
		n++

		n += ubtoa(b, n, p[2])
		b[n] = '.'
		n++

		n += ubtoa(b, n, p[3])
		return string(b[:n])
	}
	return m.String()
}

// ubtoa encodes the string form of the integer v to dst[start:] and
// returns the number of bytes written to dst. The caller must ensure
// that dst has sufficient length.
func ubtoa(dst []byte, start int, v byte) int {
	if v < 10 {
		dst[start] = v + '0'
		return 1
	} else if v < 100 {
		dst[start+1] = v%10 + '0'
		dst[start] = v/10 + '0'
		return 2
	}

	dst[start+2] = v%10 + '0'
	dst[start+1] = (v/10)%10 + '0'
	dst[start] = v/100 + '0'
	return 3
}