#!/bin/bash

common::envvars() {
  if [ -z "${BUILD_NUMBER}" -a -n "${TRAVIS_BUILD_NUMBER}" ]; then
    BUILD_NUMBER=${TRAVIS_BUILD_NUMBER}
  fi
  if [ -z "${BUILD_BRANCH}" -a -n "${TRAVIS_BRANCH}" ]; then
    BUILD_BRANCH=${TRAVIS_BRANCH}
  fi
  if [ -z "${BUILD_TYPE_ID}" -a -n "${TRAVIS_BUILD_ID}" ]; then
    BUILD_TYPE_ID=${TRAVIS_BUILD_ID}
  fi
  if [ -z "${BUILD_URL}" -a -n "${TRAVIS_BUILD_ID}" ]; then
    BUILD_URL="https://travis-ci.com/${TRAVIS_REPO_SLUG}/builds/${TRAVIS_BUILD_ID}"
  fi
  if [ -z "${BUILD_GIT_COMMIT}" -a -n "${TRAVIS_COMMIT}" ]; then
    BUILD_GIT_COMMIT=${TRAVIS_COMMIT}
  fi
  if [ -z "${BUILD_GIT_URL}"  ]; then
    BUILD_GIT_URL=$(git config --get remote.origin.url)
  fi
}