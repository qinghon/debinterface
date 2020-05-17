package debinterface

import "net"

type Bridge struct {
	adp Interface
	Script
}
type Script struct {
	Interface
	PreUp, Up, PostUp       []string
	PreDown, Down, PostDown []string
}
type Address struct {
	Interface
	Address net.IP
	Mask    net.IPMask
}
type PPP Interface

func (b *Bridge) SetBridgeOpts(opts []string) {
	b.Interface.SetBridgePorts("", "")

}
func (b *Bridge) GetBridgeOpts(opts []string) {

}
func (p *PPP) New() {

}
