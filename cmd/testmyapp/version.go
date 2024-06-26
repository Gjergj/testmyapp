package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3/ffcli"
)

// set by goreleaser at build time ,https://goreleaser.com/cookbooks/using-main.version/
var version string

func versionCommand() *ffcli.Command {
	fs := flag.NewFlagSet("testmyapp version", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "version",
		FlagSet:    fs,
		Exec: func(_ context.Context, args []string) error {
			fmt.Println(version)
			return nil
		},
	}
}
