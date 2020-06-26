// Copyright 2020 Cisco Systems, Inc.
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

package hostagent

import (
	"fmt"
	"net"
	osexec "os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"

	md "github.com/noironetworks/aci-containers/pkg/metadata"
)

var ipPath string
var bridgePath string

func initPaths() error {
	type Path struct {
		path *string
		cmd  string
	}

	pathList := []Path{
		{
			path: &ipPath,
			cmd:  "ip",
		},

		{
			path: &bridgePath,
			cmd:  "bridge",
		},
	}

	for _, p := range pathList {
		if *p.path != "" {
			continue
		}
		path, err := osexec.LookPath(p.cmd)
		if err != nil {
			return err
		}

		*p.path = path
	}

	return nil
}

func fetchIf(ipAddr string) (string, string, error) {
	var mac string
	var ifname string

	err := initPaths()
	if err != nil {
		return "", "", err
	}

	searchOut := func(data, key string, keyIndex, valIndex int) string {
		lines := strings.Split(data, "\n")
		for _, l := range lines {
			entries := strings.Split(l, " ")
			if entries[keyIndex] == key {
				if len(entries) > valIndex {
					return entries[valIndex]
				}
				return ""
			}
		}
		return ""
	}
	neighbors, err := osexec.Command(ipPath, "neighbor", "show").CombinedOutput()
	if err != nil {
		return "", "", err
	}

	mac = searchOut(string(neighbors), ipAddr, 0, 4)
	if mac == "" {
		return "", "", fmt.Errorf("mac not found")
	}

	fdb, err := osexec.Command(bridgePath, "fdb", "show").CombinedOutput()
	if err != nil {
		return "", "", err
	}

	ifname = searchOut(string(fdb), mac, 0, 2)
	return mac, ifname, nil
}

func (agent *HostAgent) snoopContainerIfaces(pod *v1.Pod) {
	if !agent.config.SnoopOnly {
		return
	}

	if pod.Status.PodIP == "" {
		return
	}

	logger := agent.log.WithFields(logrus.Fields{
		"pod":       pod.ObjectMeta.Name,
		"namespace": pod.ObjectMeta.Namespace,
	})

	logger.Debug("Snooping veth")
	epKey := pod.ObjectMeta.Namespace + "/" + pod.ObjectMeta.Name

	if _, ok := agent.epMetadata[epKey]; ok {
		return
	}

	mac, ifName, err := fetchIf(pod.Status.PodIP)
	if err != nil {
		logger.Errorf("Error snooping IP %s - %v", pod.Status.PodIP, err)
		return
	}

	ips := []md.ContainerIfaceIP{
		{Address: net.IPNet{
			IP: net.ParseIP(pod.Status.PodIP),
		},
		},
	}
	cid := strings.TrimPrefix(pod.Status.ContainerStatuses[0].ContainerID, "docker://")
	metadata := &md.ContainerMetadata{
		Id: md.ContainerId{
			Namespace: pod.ObjectMeta.Namespace,
			Pod:       pod.ObjectMeta.Name,
			ContId:    cid,
		},
		Ifaces: []*md.ContainerIfaceMd{
			{
				HostVethName: ifName,
				Mac:          mac,
				IPs:          ips,
			},
		},
	}
	agent.epMetadata[epKey] =
		make(map[string]*md.ContainerMetadata)
	agent.epMetadata[epKey][metadata.Id.ContId] = metadata
}
