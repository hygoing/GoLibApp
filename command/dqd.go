package main

import (
	"encoding/binary"
	"net"
	"fmt"
)

func main() {
	var vni uint32 = 100
	vni64 := uint64(vni)
	fmt.Println(matchMacWithMask(vni64,24,24))
}

func matchMacWithMask(data uint64, from, len uint) string {
	// get 0x80, return 00:00:80:00:00:00/ff:ff:ff:00:00:00
	if from+len > 48 {
		panic("bad mac match from/to > 48")
	}
	to := from + len
	a, b := uint64(data), (uint64(0xffffffffffffffff)<<from)^(uint64(0xffffffffffffffff)<<(to))
	a = a << from
	bt := make([]byte, 8)
	binary.BigEndian.PutUint64(bt, a)
	str := net.HardwareAddr(bt[2:]).String()
	binary.BigEndian.PutUint64(bt, b)
	str += "/" + net.HardwareAddr(bt[2:]).String()
	return str
}