#!/bin/bash

export PATH=/go/bin:/usr/local/go/bin:${PATH}

#store the present working directory
PRE_PWD=`pwd`
CONF_WD=$PRE_PWD/internal/config
LOG_WD=$PRE_PWD/internal/logger
MET_WD=$PRE_PWD/internal/metric
UTL_WD=$PRE_PWD/internal/utils

#goprobe path
SDK_PATH="src/github.com/CA-Application-Performance-Management/goprobe"

export GOPATH=${WORKSPACE}
go get github.com/satori/go.uuid

for package in $PRE_PWD $CONF_WD $LOG_WD $MET_WD $UTL_WD; do
  cd $package
  #This builds the packages of GO SDK
  go build
  EXIT_STATUS=$?
  if [ $EXIT_STATUS == 0 ]; then
  echo "Build" $package "Finished"
  else
    echo "Build Failed"
    exit $EXIT_STATUS
  fi
  #This runs the unittests of the packages of the SDK
  go test
  EXIT_STATUS=$?
  if [ $EXIT_STATUS == 0 ]; then
    echo "Test" $package "Finished"
  else
    echo "Test Failed"
    exit $EXIT_STATUS
  fi
done

exit $EXIT_STATUS
