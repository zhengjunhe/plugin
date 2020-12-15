// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.8

package main

import (
	_ "github.com/33cn/dplatformos/system"
	"github.com/33cn/plugin/cli/buildflags"
	_ "github.com/33cn/plugin/plugin"

	"github.com/33cn/dplatformos/util/cli"
)

func main() {
	if buildflags.RPCAddr == "" {
		buildflags.RPCAddr = "http://localhost:28803"
	}
	cli.Run(buildflags.RPCAddr, buildflags.ParaName, "")
}
