#!/bin/bash

readonly BUILD_ROOT=$(
  unset CDPATH
  build_root=$(dirname "${BASH_SOURCE}")/..
  cd "${build_root}"
  pwd
)

source "${BUILD_ROOT}/scripts/common.sh"
source "${BUILD_ROOT}/scripts/docker-build.sh"
