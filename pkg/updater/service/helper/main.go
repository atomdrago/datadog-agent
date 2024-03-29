// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package main is a package that allows dd-updater
// to execute a subset of priviledged commands
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unicode"
)

var (
	installPath string
	systemdPath = "/lib/systemd/system" // todo load it at build time from omnibus
	pkgDir      = "/opt/datadog-packages"
	testSkipUID = ""
)

func enforceUID() bool {
	return testSkipUID != "true"
}

type privilegeCommand struct {
	Command string `json:"command,omitempty"`
	Unit    string `json:"unit,omitempty"`
	Path    string `json:"path,omitempty"`
	Target  string `json:"target,omitempty"`
}

func isValidUnitChar(c rune) bool {
	return unicode.IsLower(c) || c == '.' || c == '-'
}

func isValidUnitString(s string) bool {
	for _, c := range s {
		if !isValidUnitChar(c) {
			return false
		}
	}
	return true
}

func buildHelperPath(target string) (path string, err error) {
	updaterHelperPath := filepath.Join(pkgDir, "updater", target, "bin/updater/updater-helper")
	info, err := os.Stat(updaterHelperPath)
	if err != nil {
		return "", err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return "", fmt.Errorf("couldn't get update helper stats: %w", err)
	}
	ddUpdaterUser, err := user.Lookup("dd-updater")
	if err != nil {
		return "", fmt.Errorf("failed to lookup dd-agent user: %s", err)
	}
	if ddUpdaterUser.Uid != strconv.Itoa(int(stat.Uid)) {
		return "", fmt.Errorf("updater-helper should be owned by dd-agent")
	}

	if stat.Mode != 750 {
		return "", fmt.Errorf("updater-helper should only be executable by the user")
	}
	return updaterHelperPath, nil
}

func buildCommand(inputCommand privilegeCommand) (*exec.Cmd, error) {
	if inputCommand.Command == "systemd-reload" {
		return exec.Command("systemctl", "daemon-reload"), nil
	}

	if inputCommand.Unit != "" {
		return buildUnitCommand(inputCommand)
	}
	if inputCommand.Target != "" {

	}
	return buildPathCommand(inputCommand)
}

func buildUnitCommand(inputCommand privilegeCommand) (*exec.Cmd, error) {
	command := inputCommand.Command
	unit := inputCommand.Unit
	if !strings.HasPrefix(unit, "datadog-") || !isValidUnitString(unit) {
		return nil, fmt.Errorf("invalid unit")
	}
	switch command {
	case "stop", "enable", "disable":
		return exec.Command("systemctl", command, unit), nil
	case "start":
		// --no-block is used to avoid waiting on oneshot executions
		return exec.Command("systemctl", command, unit, "--no-block"), nil
	case "load-unit":
		return exec.Command("cp", filepath.Join(installPath, "systemd", unit), filepath.Join(systemdPath, unit)), nil
	case "remove-unit":
		return exec.Command("rm", filepath.Join(systemdPath, unit)), nil
	default:
		return nil, fmt.Errorf("invalid command")
	}
}

func buildPathCommand(inputCommand privilegeCommand) (*exec.Cmd, error) {
	path := inputCommand.Path
	// detects symlinks and ..
	absPath, err := filepath.Abs(path)
	if absPath != path || err != nil {
		return nil, fmt.Errorf("invalid path")
	}
	if !strings.HasPrefix(path, pkgDir) {
		return nil, fmt.Errorf("invalid path")
	}
	switch inputCommand.Command {
	case "chown dd-agent":
		return exec.Command("chown", "-R", "dd-agent:dd-agent", path), nil
	case "rm":
		return exec.Command("rm", "-rf", path), nil
	default:
		return nil, fmt.Errorf("invalid command")
	}
}

func buildTargetCommand(inputCommand privilegeCommand) (*exec.Cmd, error) {
	target := inputCommand.Target
	switch inputCommand.Command {
	case "setcap cap_setuid+ep":
		helperPath, err := buildHelperPath(target)
		if err != nil {
			return nil, err
		}
		return exec.Command("setcap", "cap_setuid+ep", helperPath), nil
	default:
		return nil, fmt.Errorf("invalid command")
	}
}

func executeCommand() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("wrong number of arguments")
	}
	inputCommand := os.Args[1]

	var pc privilegeCommand
	err := json.Unmarshal([]byte(inputCommand), &pc)
	if err != nil {
		return fmt.Errorf("decoding command")
	}

	currentUser := syscall.Getuid()
	command, err := buildCommand(pc)
	if err != nil {
		return err
	}

	// only root or dd-updater can execute this command
	if currentUser != 0 && enforceUID() {
		ddUpdaterUser, err := user.Lookup("dd-updater")
		if err != nil {
			return fmt.Errorf("failed to lookup dd-updater user: %s", err)
		}
		if strconv.Itoa(currentUser) != ddUpdaterUser.Uid {
			return fmt.Errorf("only root or dd-updater can execute this command")
		}
		if err := syscall.Setuid(0); err != nil {
			return fmt.Errorf("failed to setuid: %s", err)
		}
		defer func() {
			err := syscall.Setuid(currentUser)
			if err != nil {
				log.Printf("Failed to set back to current user: %s", err)
			}
		}()
	}

	log.Printf("Running command: %s", command.String())
	return command.Run()
}

func main() {
	log.SetOutput(os.Stdout)
	err := executeCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
