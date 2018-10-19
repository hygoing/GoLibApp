package main

import (
	"bytes"
	nlb2 "jd.com/lb/jstack-lb-common/message/nlb"
	"fmt"
	"encoding/json"
)

func main() {
	map2 := map[string]*nlb2.Nlb{
		"one" :&nlb2.Nlb{
			Id:      "netlb-1q2w3e4r5t",
			Version: 100,
			Vni:     100,
			Vip:     "10.0.0.1",
			Mac:     "11:22:33:44:55:66",
			GwIp:    "10.0.0.100",
			Cidr:    "10.0.0.1/24",
			Locals: []*nlb2.Local{
				&nlb2.Local{
					Ip:       "1.1.1.1",
					Mac:      "22:22:22:22:22:22",
					HostName: "a-1",
				},
			},
		},
	}

	map1 := make(map[string]*nlb2.Nlb)
	map1["one"] = &nlb2.Nlb{}
	err := deepCopy(map1["one"], map2["one"])
	if err!=nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n",map2["one"])
	fmt.Printf("%+v\n",map1["one"])

	map2["one"].Vip = "99.99.99.99"

	fmt.Printf("%+v\n",map2["one"])
	fmt.Printf("%+v",map1["one"])
}

func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return json.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
