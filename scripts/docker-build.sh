#!/bin/bash

docker::build_and_push() {
  common::envvars

  local image_name="${1}"
  local file_name="${2:-Dockerfile}"
  local bid="$BID"
  local branch="${TRAVIS_BRANCH}"

  docker::build_file "${file_name}" "${image_name}" "${bid}" && \
  docker::push_and_tag "${image_name}" "${bid}" "${branch}"
}

docker::build() {
  docker::build_file "Dockerfile" "$@"
}

docker::build_file() {
  local file_name="$1"
  local image_name="$2"
  local bid="$3"
  local branch="${TRAVIS_BRANCH}"
   
  local buildTypeId="${TRAVIS_BUILD_ID}"

  if [ -n "$buildTypeId" ]; then
      cp ${file_name} ${file_name}_backup
      GIT_REPO=$(git config --get remote.origin.url)
      echo "" >> ${file_name}
      echo "LABEL com.tibco.flogo.ci.buildNumber=\"${bid}\" \\" >> ${file_name}
      echo " com.tibco.flogo.ci.buildTypeId=\"${buildTypeId}\" \\" >> ${file_name}
      echo " com.tibco.flogo.ci.url=\"https://travis-ci.org/${TRAVIS_REPO_SLUG}/builds/${buildTypeId}\"" >> ${file_name}
      if [ -n "${branch}" ]; then
          echo "LABEL com.tibco.flogo.git.repo=\"$GIT_REPO\" com.tibco.flogo.git.commit=\"${TRAVIS_COMMIT}\"" >> ${file_name}
      fi
  fi

  local latest='latest'
  if [[ -n ${branch} && ( ${branch} != 'master' ) ]]; then
    latest="latest-${branch}"
    bid="${bid}-${branch}"
  fi

  local latest='latest'

  docker build --force-rm=true --rm=true -t $image_name:${bid:-latest} -f ${file_name} .
  rc=$?

  if [ -e "${file_name}_backup" ]; then
      rm ${file_name}
      mv ${file_name}_backup ${file_name}
  fi

  if [ $rc -ne 0 ]; then
    echo "Build failed"
    exit $rc
  fi
}

docker::pull_and_tag() {
  local base_name="$1"
  local docker_registry=$( [ -n "${DOCKER_REGISTRY}" ] && echo "${DOCKER_REGISTRY}/" || echo "" ) 

  if [ -n "$docker_registry" ]; then
    docker pull ${docker_registry}${base_name} && \
    docker tag ${docker_registry}${base_name} ${base_name} && \
    docker rmi ${docker_registry}${base_name}
  fi
}

docker::push_and_tag() {
  common::envvars
  local image_name="$1"
  local bid="$2"
  local docker_registry=$( [ -n "${DOCKER_REGISTRY}" ] && echo "${DOCKER_REGISTRY}/" || echo "" ) 
  local branch="${TRAVIS_BRANCH}"
  local latest='latest'
  if [[ -n ${branch} && ( ${branch} != 'master' ) ]]; then
    latest="latest-${branch}"
    bid="${bid}-${branch}"
  fi

  if [ -n "${bid}" -a -n "${branch}" ]; then
    echo "Publishing image..."
    docker tag ${image_name}:${bid} ${image_name}:${latest} && \
    docker tag ${image_name}:${bid} ${docker_registry}${image_name}:${bid} && \
    docker tag ${image_name}:${bid} ${docker_registry}${image_name}:${latest} && \
    docker images | grep ${image_name} >> images.txt && \
    docker push ${docker_registry}${image_name}:${latest} && \
    docker push ${docker_registry}${image_name}:${bid} && \
    docker rmi ${docker_registry}${image_name}:${latest} && \
    docker rmi ${docker_registry}${image_name}:${bid} && \
    echo "Done."
  fi
}


docker::copy_tag_and_push() {
  common::envvars
  local src_image_name="$1"
  local dest_image_name="$2"
  local bid="${3:-${BID}}"
  local docker_registry=$( [ -n "${DOCKER_REGISTRY}" ] && echo "${DOCKER_REGISTRY}/" || echo "" ) 
  local branch="${TRAVIS_BRANCH}"
  local latest='latest'
  # non-branch-aware TeamCity jobs won't have $IS_MASTER at all
  if [[ -n ${branch} && ( ${branch} != 'master' ) ]]; then
    latest="latest-${branch}"
    bid="${bid}-${branch}"
  fi
  

  if [ -n "${bid}" -a -n "${branch}" ]; then
    echo "Retagging image from: ${src_image_name}:${bid} to: ${dest_image_name}:${bid} ..."
    docker tag ${src_image_name}:${bid} ${dest_image_name}:${latest} && \
    docker tag ${src_image_name}:${bid} ${docker_registry}${dest_image_name}:${bid} && \
    docker tag ${src_image_name}:${bid} ${docker_registry}${dest_image_name}:${latest} && \
    docker push ${docker_registry}${dest_image_name}:${latest} && \
    docker push ${docker_registry}${dest_image_name}:${bid} && \
    docker rmi ${docker_registry}${dest_image_name}:${latest} && \
    docker rmi ${docker_registry}${dest_image_name}:${bid} && \
    echo "Done."
  else
     # no bid and docker registry i.e. local machine
     docker tag  ${src_image_name}:${latest} ${dest_image_name}:${latest}
  fi
}
