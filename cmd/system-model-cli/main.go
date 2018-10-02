/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package main

import (
	"github.com/nalej/system-model/cmd/system-model-cli/commands"
	"github.com/nalej/system-model/version"
)

var MainVersion string

var MainCommit string

func main() {
	version.AppVersion = MainVersion
	version.Commit = MainCommit
	commands.Execute()
}
