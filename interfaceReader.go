// debian network config parse package
package debinterface

import (
	"bufio"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var INTERFACES_FILE = "/etc/network/interfaces"

type InterfaceReader struct {
	Adapters       []Interface
	Interfaces     []Interfaces
	context        int
	headerComments string
	selection      map[string][]string
	filepath       string
}

func NewReader(_filepath string) *InterfaceReader {
	var reader = InterfaceReader{}
	reader.context = -1
	if _filepath == "" {
		reader.filepath = INTERFACES_FILE
	} else {
		reader.filepath = _filepath
	}
	reader.selection = make(map[string][]string)
	return &reader
}
func (reader *InterfaceReader) Read() error {
	return reader.ReadLines(reader.filepath)
}
func (reader *InterfaceReader) ReadLines(_filepath string) error {
	fp, err := os.Open(_filepath)
	if err != nil {
		return err
	}
	defer fp.Close()
	var i Interfaces
	i.FilePath = _filepath
	reader.Interfaces = append(reader.Interfaces, i)
	br := bufio.NewReader(fp)
	var header = false
	var lineStr = ""
	var lineByte []byte
	var e error
	for {
		lineByte, _, e = br.ReadLine()
		if e == io.EOF {
			break
		}
		//println(string(lineByte))
		lineStr = string([]rune(string(lineByte)))
		if strings.TrimSpace(lineStr) == "" {
			continue
		}
		if strings.TrimSpace(lineStr)[0:1] == "#" {
			if !header {
				reader.headerComments += lineStr
			}
			continue
		}
		header = false
		// 去掉行尾 "#"
		if strings.Index(lineStr, "#") != -1 {
			lineStr = lineStr[:strings.Index(lineStr, "#")]
		}
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		reader.readAuto(lineStr)
		reader.parseIface(lineStr)
		if err = reader.parseDetails(lineStr); err != nil {
			return err
		}
	}
	reader.parseAuto()
	return nil
}
func (reader *InterfaceReader) ReadFromRD(rd io.Reader) error {
	br := bufio.NewReader(rd)
	var header = false
	var lineStr = ""
	var lineByte []byte
	var e, err error
	for {
		lineByte, _, e = br.ReadLine()
		if e == io.EOF {
			break
		}
		//println(string(lineByte))
		lineStr = string([]rune(string(lineByte)))
		if strings.TrimSpace(lineStr) == "" {
			continue
		}
		if strings.TrimSpace(lineStr)[0:1] == "#" {
			if !header {
				reader.headerComments += lineStr
			}
			continue
		}
		header = false
		// 去掉行尾 "#"
		if strings.Index(lineStr, "#") != -1 {
			lineStr = lineStr[:strings.Index(lineStr, "#")]
		}
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		reader.readAuto(lineStr)

		reader.parseIface(lineStr)
		if err = reader.parseDetails(lineStr); err != nil {
			return err
		}
	}
	reader.parseAuto()
	return nil
}
func (reader *InterfaceReader) parseIface(line string) {
	sline := strings.Fields(line)
	if sline[0] == "iface" {
		reader.context = reader.context + 1
		var adapter = Interface{}
		adapter.SetName(sline[1])
		adapter.SetAddrFam(sline[2])
		adapter.SetAddrSource(sline[3])
		//adapter.SetFromFile(filepath)
		reader.Adapters = append(reader.Adapters, adapter)

	}
}
func (reader *InterfaceReader) readAuto(line string) {
	sline := strings.Fields(line)
	switch sline[0] {
	//case "auto":
	//	reader.auto = append(reader.auto, sline[1])
	//case "hotplug":
	//	reader.hotplug = append(reader.hotplug, sline[1])
	case "auto", "hotplug", "allow-hotplug", "no-auto-down", "no-scripts":
		reader.selection[sline[1]] = append(reader.selection[sline[1]], sline[0])
	}
}

//
func (reader *InterfaceReader) parseAuto() {
	for i, a := range reader.Adapters {
		if len(reader.selection[a.GetName()]) != 0 {
			for _, v := range reader.selection[a.GetName()] {
				reader.Adapters[i].SetSelection(v)
			}
		}
	}
}
func (reader *InterfaceReader) parseDetails(line string) error {
	sline := strings.Fields(line)
	//log.Println(sline,len(sline))
	if len(sline) < 2 {
		return os.ErrInvalid
	}

	switch sline[0] {
	case "auto", "iface", "allow-hotplug", "no-auto-down", "no-scripts":
		return nil
	case "address":
		if strings.Contains(sline[1], "/") {
			ip, ipnet, err := net.ParseCIDR(sline[1])
			if err == nil {
				reader.Adapters[reader.context].SetAddress(ip)
				reader.Adapters[reader.context].SetNetmask(ipnet.Mask)
			} else {
				//log.Println(err)
				return err
			}
		} else {
			reader.Adapters[reader.context].SetAddress(net.ParseIP(sline[1]))
		}
	case "netmask":
		ip := net.ParseIP(sline[1])
		// todo ipv6 兼容
		ipmask := net.IPv4Mask(ip[len(ip)-4], ip[len(ip)-3], ip[len(ip)-2], ip[len(ip)-1])
		reader.Adapters[reader.context].SetNetmask(ipmask)
	case "gateway":
		reader.Adapters[reader.context].SetGateWay(net.ParseIP(sline[1]))
	case "mtu":
		mtu, err := strconv.Atoi(sline[1])
		if err != nil {
			return err
		}
		reader.Adapters[reader.context].SetMtu(mtu)
	case "scope":
		reader.Adapters[reader.context].SetScope(sline[1])
	case "hwaddress":
		mac, err := net.ParseMAC(sline[1])
		if err == nil {
			reader.Adapters[reader.context].SetHwAddress(mac)
		} else {
			//log.Println(err)
			return err
		}
	case "dns-search":
		reader.Adapters[reader.context].SetDnsSearch(sline[1:])
	case "dns-nameservers":
		var ipList []net.IP
		for _, p := range sline[1:] {
			ipList = append(ipList, net.ParseIP(p))
		}
		reader.Adapters[reader.context].SetDnsNameServer(ipList)
	case "broadcast":
		reader.Adapters[reader.context].SetBroadcast(net.ParseIP(sline[1]))
	case "up", "down", "pre-up", "pre-down", "post-up", "post-down":
		reader.Adapters[reader.context].SetScript(sline[0], strings.Join(sline[1:], " "))
	case "source", "source-directory":
		if sline[0] == "source-directory" {
			//todo source-directory parse
			return nil
		}
		files, err := filepath.Glob(sline[1])
		if err != nil {
			//log.Println(err)
			return err
		}
		for _, f := range files {
			reader.ReadLines(f)
		}
	case "pointopoint":
		reader.Adapters[reader.context].SetPoinToPoint(net.ParseIP(sline[1]))
	case "provider":
		reader.Adapters[reader.context].SetProvider(sline[1])
	case "bridge_ports":
		switch sline[1] {
		case "none":
			reader.Adapters[reader.context].SetBridgePorts(sline[1], "")
		case "regex":
			reader.Adapters[reader.context].SetBridgePorts(sline[1], sline[2])
		default:
			reader.Adapters[reader.context].SetBridgePorts("", sline[1])
		}
	case "bridge_fd":
		fd, err := strconv.Atoi(sline[1])
		if err == nil {
			reader.Adapters[reader.context].SetBridgeFd(fd)
		} else {
			reader.Adapters[reader.context].SetUnknown(line)
		}
	case "bridge_waitport":
		wp, err := strconv.Atoi(sline[1])
		if err == nil {
			reader.Adapters[reader.context].SetBridgeWaitPort(wp)
		} else {
			reader.Adapters[reader.context].SetUnknown(line)
		}
	case "bridge_stp":
		if sline[1] == "off" {
			reader.Adapters[reader.context].SetBridgeStp(false)
		} else if sline[1] == "on" {
			reader.Adapters[reader.context].SetBridgeStp(true)
		}
	default:
		//log.Println("Unknow options",sline)
		if reader.context != -1 {
			reader.Adapters[reader.context].SetUnknown(line)
		}
	}
	return nil
}
