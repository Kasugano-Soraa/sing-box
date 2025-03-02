#!/usr/bin/env bash

PROJECTS=$(dirname "$0")/../..

go get -x github.com/sagernet/sing@$(git -C $PROJECTS/sing rev-parse HEAD)
go get -x github.com/sagernet/sing-dns@$(git -C $PROJECTS/sing-dns rev-parse HEAD)
go get -x github.com/sagernet/sing-tun@$(git -C $PROJECTS/sing-tun rev-parse HEAD)
go get -x github.com/sagernet/sing-shadowsocks@$(git -C $PROJECTS/sing-shadowsocks rev-parse HEAD)
go get -x github.com/sagernet/sing-vmess@$(git -C $PROJECTS/sing-vmess rev-parse HEAD)
go mod tidy
