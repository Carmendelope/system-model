/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package main

import (
	"github.com/nalej/system-model/cmd/system-model-cli/commands"
	"github.com/nalej/system-model/version"
)

// MainVersion with the application version.
var MainVersion string

// MainCommit with the commit id.
var MainCommit string

func main() {
	version.AppVersion = MainVersion
	version.Commit = MainCommit
	commands.Execute()
}
