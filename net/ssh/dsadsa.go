package main

func main() {
	pin = "boss-01"
}

type VpnTunnel struct {
	Id              string     `json:"id"`
	Version         uint64     `json:"version"`
	TunnelVni       uint32     `json:"tunnelVni"`
	BgwVni          uint32     `json:"bgwVni"`
	Psk             string     `json:"psk"`
	LocalId         string     `json:"localId"`
	RemoteId        string     `json:"remoteId"`
	LocalFip        string     `json:"localFip"`
	RemoteFip       string     `json:"remoteFip"`
	TunnelLocalIp   string     `json:"tunnelLocalIp"`
	TunnelLocalMac  string     `json:"tunnelLocalMac"`
	TunnelRemoteIp  string     `json:"tunnelRemoteIp"`
	TunnelRemoteMac string     `json:"tunnelRemoteMac"`
	BandwidthIn     uint32     `json:"bandwidthIn"`
	BandwidthOut    uint32     `json:"bandwidthOut"`
	Status          uint32     `json:"status"`
	AdminStatus     uint32     `json:"adminStatus"`
	IkeConn         *IkeConn   `json:"ikeConn"`
	IpsecConn       *IpsecConn `json:"ipsecConn"`
}

type IkeConn struct {
	IkeVersion   uint32 `json:"ikeVersion"`
	IkeMode      uint32 `json:"ikeMode"`
	IkeDpd       uint32 `json:"ikeDpd"`
	IkeProposals uint32 `json:"ikeProposals"`
}

type IpsecConn struct {
	RekeyTime      uint32 `json:"rekeyTime"`
	RekeyPackets   uint64 `json:"rekeyPackets"`
	RekeyBytes     uint64 `json:"rekeyBytes"`
	Priority       uint32 `json:"proirity"`
	IpsecProposals uint32 `json:"ipsecProposals"`
}
