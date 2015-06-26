#!/bin/bash

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

set -o errexit
set -o nounset
set -o pipefail

declare -r KUBE_RELEASE_BUCKET="https://storage.googleapis.com/kubernetes-release"
declare -r KUBE_BINARY_NAME="kubernetes.tar.gz"

usage() {
	echo "usage:
  $0 [stable|release|latest|latest-green]

        stable:        latest stable version
        release:       latest release candidate
        latest:        latest ci build
        latest-green:  latest ci build to pass gce e2e"
}

if [[ $# -ne 1 ]]; then
    usage
    exit 1
fi

if [[ "$1" == "latest" ]]; then
    # latest ci version is written in the latest.txt in the /ci dir
    KUBE_BINARY_TAR_RELATIVE_PATH=ci
    KUBE_VERSION_FILE=latest.txt
elif [[ "$1" == "latest-green" ]]; then
    # latest ci version to pass gce e2e is written in the latest-green.txt in the /ci dir
    KUBE_BINARY_TAR_RELATIVE_PATH=ci
    KUBE_VERSION_FILE=latest-green.txt
elif [[ "$1" == "stable" ]]; then
    # latest stable release version is written in the stable.txt file in the /release dir
    KUBE_BINARY_TAR_RELATIVE_PATH=release
    KUBE_VERSION_FILE=stable.txt
elif [[ "$1" == "release" ]]; then
    # latest release candidate version is written in latest.txt in the /release dir
    KUBE_BINARY_TAR_RELATIVE_PATH=release
    KUBE_VERSION_FILE=latest.txt
else
    usage
    exit 1
fi

KUBE_BINARY_DIRECTORY=$KUBE_RELEASE_BUCKET/$KUBE_BINARY_TAR_RELATIVE_PATH
KUBE_VERSION=$(curl --silent --fail $KUBE_BINARY_DIRECTORY/$KUBE_VERSION_FILE)
KUBE_BINARY_URL=$KUBE_BINARY_PATH/$KUBE_VERSION/$KUBE_BINARY_NAME

curl --fail -o kubernetes-$KUBE_VERSION.tar.gz $KUBE_BINARY_URL
