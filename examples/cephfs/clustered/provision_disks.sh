#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
set -x

KUBE_MIG="kubernetes-minion-group"
DISK_SIZE="2TB"

KUBE_NODES=`gcloud compute instance-groups managed list-instances -q --format json kubernetes-minion-group \
    | jq  -r '.[].instance' \
    | tr '\n' ' '`
KUBE_DISKS=`echo "$KUBE_NODES" | sed 's/[[:space:]]/-ceph-disk /g'`

gcloud compute disks create $KUBE_DISKS --size $DISK_SIZE

for node in $KUBE_NODES; do
  gcloud compute instances attach-disk "$node" \
    --device-name "sdb" \
    --disk "$node-ceph-disk" \
    --mode rw &
done;

wait

read block

for node in $KUBE_NODES; do
  gcloud compute instances detach-disk "$node" \
    --disk "$node-ceph-disk" &
done;

wait

gcloud compute disks delete $KUBE_DISKS