#!/bin/bash

common::envvars() {
  if [ -z "$BID" -a -n "${TRAVIS_BUILD_NUMBER}" ]; then
    BID=${TRAVIS_BUILD_NUMBER}
  fi
}