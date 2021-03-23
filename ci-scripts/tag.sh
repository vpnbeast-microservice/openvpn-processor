#!/bin/bash

function git_global_settings() {
    git config --global user.name ${USERNAME}
    git config --global user.email ${EMAIL}
}

function git_commit_and_push() {
    git --no-pager diff
    git add --all
    git commit -am "[ci-skip] version ${RELEASE_VERSION}.RELEASE"
    git tag -a "v${RELEASE_VERSION}" -m "v${RELEASE_VERSION} tagged"
    git status
    git push --follow-tags ${PUSH_URL} HEAD:${BRANCH}
}

function increment_minor_version() {
    local version_major=$(echo $1 | cut -d "." -f 1)
    local version_patch=$(echo $1 | cut -d "." -f 2)
    local version_minor=$(echo $1 | cut -d "." -f 3)
    version_minor=`expr ${version_minor} + 1`
    echo "${version_major}.${version_patch}.${version_minor}"
}

function set_version() {
    sed -i "s/${1}/${2}/g" version.properties
}

function set_chart_version() {
    sed -i "s/${CURRENT_VERSION}/${RELEASE_VERSION}/g" charts/${CHART_NAME}/Chart.yaml
}

set -ex
USERNAME=vpnbeast-ci
EMAIL=info@thevpnbeast.com
PROJECT_NAME=openvpn-processor
GIT_ACCESS_TOKEN=$1
CHART_NAME=openvpn-processor
BRANCH=master
CURRENT_VERSION=`grep RELEASE_VERSION version.properties | cut -d "=" -f2`
RELEASE_VERSION=$(increment_minor_version ${CURRENT_VERSION})
PUSH_URL=https://${USERNAME}:${GIT_ACCESS_TOKEN}@github.com/vpnbeast/${PROJECT_NAME}.git

set_version $CURRENT_VERSION $RELEASE_VERSION
set_chart_version
git_global_settings
git_commit_and_push
