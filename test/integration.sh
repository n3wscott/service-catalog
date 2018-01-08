#!/bin/bash

# Copyright 2017 The Kubernetes Authors.
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

# This emulates "readlink -f" which is not available on MacOS X.
function sc::readlinkdashf {
  # run in a subshell for simpler 'cd'
  (
    if [[ -d "$1" ]]; then # This also catch symlinks to dirs.
      cd "$1"
      pwd -P
    else
      cd $(dirname "$1")
      local f
      f=$(basename "$1")
      if [[ -L "$f" ]]; then
        readlink "$f"
      else
        echo "$(pwd -P)/${f}"
      fi
    fi
  )
}

# This emulates "realpath" which is not available on MacOS X
sc::realpath() {
  if [[ ! -e "$1" ]]; then
    echo "$1: No such file or directory" >&2
    return 1
  fi
  sc::readlinkdashf "$1"
}


# this script resides in the `test/` folder at the root of the project
KUBE_ROOT=$(sc::realpath $(dirname "${BASH_SOURCE}")/../pkg/kubernetes)
source "${KUBE_ROOT}/hack/lib/init.sh"
GOROOT=$(go env GOROOT)

runTests() {
  # kube::etcd::start

  if [[ -w ${GOROOT}/pkg ]]; then
    FLAGS="-i"
  elif [[ -n ${PKGDIR:-} ]]; then
    FLAGS="-pkgdir $PKGDIR"
  else
    FLAGS=""
  fi

  go test -race $FLAGS github.com/kubernetes-incubator/service-catalog/test/integration/... -c \
      && ./integration.test -test.v $@
}

# Run cleanup to stop etcd on interrupt or other kill signal.
# trap kube::etcd::cleanup EXIT

runTests $@

