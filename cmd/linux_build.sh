#!/usr/bin/env bash
# linux server build infra-aws-operator images
# author: weiqiang; date: 2022-12
set -x
set -e

export GOARCH=amd64
export GOOS=linux
export GCCGO=gc

registrieAddress="harbor.xxx.com/devops"
version=$1
if [ -z $version ]; then
    version=v0.5
fi

imageTagVersion=${version}
servicename=infra-aws-operator
pkgname=infra-aws-operator

go build -o ${pkgname} main.go
chmod u+x ${pkgname}
tar -zcvf ${pkgname}-linux-amd64-${version}.tar.gz \
  ${pkgname} configs/config.yaml configs/${pkgname}.service README.md README-cn.md web/

# docker build -f build/dockerfile -t ${registrieAddress}/${servicename}:${imageTagVersion} .
# docker push ${registrieAddress}/${servicename}:${imageTagVersion}
