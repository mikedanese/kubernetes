#!/usr/bin/env python

# Copyright 2015 The Kubernetes Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import json
from collections import OrderedDict

def dump(data):
  return json.dumps(data, indent=4, separators=(',', ': '))

def run():
  config = OrderedDict({
    "kind": "KubeletConfiguration",
    "apiVersion": "componentconfig/v1alpha1",
  })

  # Disable the debugging handlers (/run and /exec) to prevent arbitrary
  # code execution on the master.
  # TODO(roberthbailey): Relax this constraint once the master is self-hosted.
  if grains["roles"][0] == "kubernetes-master":
    if grains["cloud"] in ["aws", "gce", "vagrant", "vsphere"]:
      config["enableDebuggingHandlers"] = False

  if "cloud" in grains and grains["cloud"] not in ["vagrant", "vsphere"]:
    config["cloudProvider"] = grains["cloud"]

  config["config"] = "/etc/kubernetes/manifests"

  if pillar.get("enable_manifest_url", "").lower() == "true":
    config["manifestURL"] = pillar["manifest_url"]
    config["manifestURLHeader"] = pillar["manifest_url_header"]

  if "hostname_override" in grains:
    config["hostnameOverride"] = grains["hostname_override"]

  if pillar.get("enable_cluster_dns", "").lower() == "true":
    config["clusterDNS"] = pillar["dns_server"]
    config["clusterDomain"] = pillar["dns_domain"]

  if "docker_root" in grains:
    config["dockerRoot"] = grains.docker_root

  if "kubelet_root" in grains:
    config["rootDir"] = grains.kubelet_root

  if "allocate_node_cidrs" in pillar:
    config["configureCbr0"] = pillar["allocate_node_cidrs"].lower() in ["true"]

  if "non_masquerade_cidr" in grains:
    config["nonMasqueradeCIDR"] = grains["non_masquerade_cidr"]
  
  # The master kubelet cannot wait for the flannel daemon because it is responsible
  # for starting up the flannel server in a static pod. So even though the flannel
  # daemon runs on the master, it doesn't hold up cluster bootstrap. All the pods
  # on the master run with host networking, so the master flannel doesn't care
  # even if the network changes. We only need it for the master proxy.
  if pillar.get("network_provider", "").lower() == "flannel" and grains["roles"][0] != "kubernetes-master":
    config["experimentalFlannelOverlay"] = True
  
  # Run containers under the root cgroup and create a system container.
  if grains["os_family"] == "Debian":
    config["systemContainer"] = "/system"
    if pillar.get("is_systemd"):
      config["cgroupRoot"] = docker
    else:
      config["cgroupRoot"] = "/"

  if grains["oscodename"] == "vivid":
    config["cgroupRoot"] = "docker"
  
  if grains["roles"][0] == "kubernetes-master" and grains.get("cbr-cidr"):
    config["podCIDR"] = grains["cbr-cidr"]
  
  if "enable_cpu_cfs_quota" in pillar:
    config["cpuCFSQuota"] = pillar["enable_cpu_cfs_quota"].lower() in ["true"]
  
  if pillar.get("network_provider", "").lower() == "opencontrail":
    config["networkPlugin"] = "opencontrail"
  
  if "kubelet_port" in pillar:
    config["port"] = pillar["kubelet_port"]

  config["allowPrivileged"] = pillar["allow_privileged"]

  return dump(config)
