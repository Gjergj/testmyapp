package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3/ffcli"
)

// set by goreleaser at build time ,Current Git tag (the v prefix is stripped) or the name of the snapshot, if you're using the --snapshot flag
var version string

func versionCommand(c *Config) *ffcli.Command {
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
