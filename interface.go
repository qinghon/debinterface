package debinterface

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Interface map[string]interface{}

type FloatKey struct {
	Option []string
	Key    string
}

// return only iface line Interface
func NewAdapter(name, AddrFam, AddrSource string) Interface {
	var adapter = make(Interface)
	adapter.SetName(name)
	adapter.SetAddrFam(AddrFam)
	adapter.SetAddrSource(AddrSource)
	return adapter
}

func (adapter Interface) Set(key string, v interface{}) {
	adapter[key] = v
}

// Set Selection for iface, like "auto", "iface", "allow-hotplug", "no-auto-down",
// "no-scripts" or custom ...
func (adapter Interface) SetSelection(sel string) {
	if adapter["selection"] == nil {
		adapter["selection"] = []string{}
	}
	adapter["selection"] = append(adapter["selection"].([]string), sel)
}
func (adapter Interface) SetAuto(f bool) {
	adapter.SetSelection("auto")
}
func (adapter Interface) SetHotplug(f bool) {
	adapter.SetSelection("allow-hotplug")
}
func (adapter Interface) SetName(name string) {
	adapter.Set("name", name)
}

// Set iface address family
// The interface name is followed by the name of the address family that  the  interface  uses.   This
// will be "inet" for TCP/IP networking, but there is also some support for IPX networking ("ipx"),
// and IPv6 networking ("inet6").  Following that is the name of the method used to configure the interface.
func (adapter Interface) SetAddrFam(addrFam string) {
	switch addrFam {
	case "inet", "inet6", "ipx", "can":
		adapter.Set("addrFam", addrFam)
	}
}

// Set address method
// The method of the interface (e.g., static), or "none" (see below).
func (adapter Interface) SetAddrSource(methodName string) {
	switch methodName {
	case "dhcp", "static", "loopback", "manual",
		"bootp", "ppp", "wvdial", "dynamic",
		"ipv4ll", "v4tunnel", "auto", "6to4", "tunnel":
		adapter.Set("method_name", methodName)
	}
}

// The loopback method
func (adapter Interface) SetLookBack(turn bool) {
	if turn {
		adapter.Set("lookback", "on")
	} else {
		adapter.Set("lookback", "off")
	}
}

// "static" Method
// Set static ip address
func (adapter Interface) SetAddress(address net.IP) {
	adapter.Set("address", address)
}

// Set netmask
func (adapter Interface) SetNetmask(mask net.IPMask) {
	adapter.Set("netmask", mask)
}

// Set Broadcast address (dotted quad, + or -). Default value: "+"
func (adapter Interface) SetBroadcast(broadcast net.IP) {
	adapter.Set("broadcast", broadcast)
}

// Metric for added routes (dhclient)
func (adapter Interface) SetMetric(metric int) {
	adapter.Set("metric", metric)
}

// Set Default gateway (dotted quad)
func (adapter Interface) SetGateWay(ip net.IP) {
	adapter.Set("gateway", ip)
}

// Address of other end point (dotted quad). Note the spelling of "point-to".
func (adapter Interface) SetPoinToPoint(ip net.IP) error {
	adapter.Set("pointopoint", ip)
	return nil
}

// Hardware address.
// "manual" "static"  "dhcp" Method
func (adapter Interface) SetHwAddress(hwaddress net.HardwareAddr) error {
	adapter.Set("hwaddress", hwaddress)
	return nil
}

// Set MTU size
//
// In "manual" "static" "dhcp" Method
func (adapter Interface) SetMtu(mtu int) {
	adapter.Set("mtu", mtu)
}

// Set Scope: "global", "link", "host"
func (adapter Interface) SetScope(scope string) error {
	switch scope {
	case "global", "link", "host":
		adapter["scope"] = scope
	default:
		return os.ErrInvalid
	}
	return nil
}

// Set Use the IP address address to communicate with the server.
// In "bootp" Method
func (adapter Interface) SetServer(server net.IP) {
	adapter.Set("server", server)
}
func (adapter Interface) SetDstAddr(dstaddr net.IP) {
	adapter.Set("dstaddr", dstaddr)
}
func (adapter Interface) SetLocal(local net.IP) {
	adapter.Set("local", local)
}

// Set "dns-nameserver"
func (adapter Interface) SetDnsNameServer(ipList []net.IP) {
	adapter.Set("dns-nameserver", ipList)
}

// Set "dns-search" ,like "local", "lan"
func (adapter Interface) SetDnsSearch(ipList []string) {
	adapter.Set("dns-search", ipList)
}

// Set "dns-domain" ,like "example.com"
func (adapter Interface) SetDnsDomain(domainList []string) {
	adapter.Set("dns-domain", domainList)
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

// Run command before bringing the interface up.  If this command fails then ifup aborts,
// refraining from marking the interface as  configured,
// prints an error message, and exits with status 0.  This behavior may change in the future.
func (adapter Interface) SetPreUpScript(command string) {
	if adapter["pre-up"] == nil {
		adapter["pre-up"] = []string{}
	}
	adapter["pre-up"] = append(adapter["pre-up"].([]string), command)
}

// up command
func (adapter Interface) SetUpScript(command string) {
	if adapter["up"] == nil {
		adapter["up"] = []string{}
	}
	adapter["up"] = append(adapter["up"].([]string), command)
}

// Run  command  after bringing the interface up.  If this command fails then ifup aborts, refraining from
// marking the interface as configured (even though it has really been configured), prints an  error
// message, and exits with status 0.  This behavior may change in the future.
func (adapter Interface) SetPostUpScript(command string) {
	if adapter["post-up"] == nil {
		adapter["post-up"] = []string{}
	}
	adapter["post-up"] = append(adapter["post-up"].([]string), command)
}
func (adapter Interface) SetDownScript(command string) {
	if adapter["down"] == nil {
		adapter["down"] = []string{}
	}
	adapter["down"] = append(adapter["down"].([]string), command)
}

// Run  command  before  taking the interface down.  If this command fails then ifdown
// aborts, marks the interface as deconfigured (even though it  has  not  really  been
// deconfigured), and exits with status 0.  This behavior may change in the future.
func (adapter Interface) SetPreDownScript(command string) {
	if adapter["pre-down"] == nil {
		adapter["pre-down"] = []string{}
	}
	adapter["pre-down"] = append(adapter["pre-down"].([]string), command)
}

// Run  command  after  taking  the interface down.  If this command fails then ifdown
// aborts, marks the interface as deconfigured, and exits with status 0.
// This  behavior may change in the future.
func (adapter Interface) SetPostDownScript(command string) {
	if adapter["post-down"] == nil {
		adapter["post-down"] = []string{}
	}
	adapter["post-down"] = append(adapter["post-down"].([]string), command)
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

// The bootp Method.
func (adapter Interface) SetBootFile(bootfile string) {
	adapter.Set("bootfile", bootfile)
}

// The ppp Method.
// Use name as the provider (from /etc/ppp/peers).
func (adapter Interface) SetProvider(provider string) {
	adapter.Set("provider", provider)
}

// Set bridge Spanning Tree Protocol
func (adapter Interface) SetBridgeStp(off bool) {
	adapter.Set("bridge_stp", off)
}

// Set delay before a port becomes available
func (adapter Interface) SetBridgeWaitPort(delay int) {
	adapter.Set("bridge_waitport", delay)
}

// Forwarding delay
func (adapter Interface) SetBridgeFd(delay int) {
	adapter.Set("bridge_fd", delay)
}

// Set bridge adapter bridge to define ports.
// The option:
//	"none": not  bind to any ports.
//  "regex": use a regular expression to define ports
// Default option is ""
func (adapter Interface) SetBridgePorts(option, eth string) {
	if option == "" && eth == "" {
		return
	}
	adapter.Set("bridge_ports", FloatKey{Option: []string{option}, Key: eth})
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

// Set unknow key, "Export()" together
func (adapter Interface) SetUnknown(unknown string) {
	if adapter["unknown"] == nil {
		adapter["unknown"] = []string{}
	}
	adapter["unknown"] = append(adapter["unknown"].([]string), unknown)
}

// Get adapter name ,like "eth0", "br0"
func (adapter Interface) GetName() string {
	if adapter["name"] == nil {
		return "<nil>"
	}
	return adapter["name"].(string)
}

// convert Interface to debian network file format
func (adapter Interface) Export() string {
	var output string

	output += fmt.Sprintf("iface %s %s %s\n", adapter["name"], adapter["addrFam"], adapter["method_name"])
	for k, v := range adapter {
		if k == "method_name" || k == "addrFam" || k == "name" || k == "fromfile" {
			continue
		}
		switch v.(type) {
		case string:
			output += fmt.Sprintf("\t%s %s\n", k, v)
		case net.IP:
			output += fmt.Sprintf("\t%s %s\n", k, v.(net.IP).String())
		case net.IPMask:
			output += fmt.Sprintf("\t%s %s\n", k, ToIPv4(v.(net.IPMask)))
		case []string:
			switch k {
			case "pre-up", "up", "post-up", "down", "pre-down", "post-down":
				for _, str := range v.([]string) {
					output += fmt.Sprintf("\t%s %s\n", k, str)
				}
			case "unknown":
				for _, str := range v.([]string) {
					output += fmt.Sprintf("\t%s\n", str)
				}
			case "selection":
				for _, str := range v.([]string) {
					output = fmt.Sprintf("%s %s\n", str, adapter["name"]) + output
				}
			default:
				output += fmt.Sprintf("\t%s %s\n", k, strings.Join(v.([]string), ""))
			}
		case []net.IP:
			ipList := []string{}
			for _, ip := range v.([]net.IP) {
				ipList = append(ipList, ip.String())
			}
			output += fmt.Sprintf("\t%s %s\n", k, strings.Join(ipList, " "))
		case int:
			output += fmt.Sprintf("\t%s %d\n", k, v)
		case FloatKey:
			output += fmt.Sprintf("\t%s %s\n", k, v.(FloatKey).String())
		case bool:
			if v.(bool) {
				output += fmt.Sprintf("\t%s on\n", k)
			} else {
				output += fmt.Sprintf("\t%s off\n", k)
			}
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

func (fk FloatKey) String() string {
	return strings.Join(fk.Option, " ") + " " + fk.Key
}
