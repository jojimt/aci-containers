// Copyright 2016 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"github.com/noironetworks/aci-containers/pkg/ipam"
)

type IPAMMetaData struct {
	UUID string    `json:"uuid,omitempty"`
	Map  []*Subnet `json:"map,omitempty"`
}

type Subnet struct {
	Network string `json:"network,omitempty"`
	VTEP    string `json:"vtep,omitempty"`
}

func WriteIPAMMetaData(file, vtep string, ipRanges *NetIps) error {
	data := genIPAMMeta(vtep, ipRanges)
	data.UUID = fmt.Sprintf("%s-%s", file, vtep)
	datacont, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, datacont, 0644)
}

func genIPAMMeta(vtep string, ipRanges *NetIps) *IPAMMetaData {
	result := &IPAMMetaData{}
	var subnets []*net.IPNet
	for _, n := range ipRanges.V4 {
		subnets = append(subnets, ipam.Range2Cidr(n.Start, n.End)...)
	}

	var sn *Subnet
	for _, nw := range subnets {
		sn = new(Subnet)
		sn.Network = nw.String()
		sn.VTEP = vtep
		result.Map = append(result.Map, sn)
	}

	return result
}
