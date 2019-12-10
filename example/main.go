package main

import (
	"fmt"

	updater "github.com/drbh/github-commit-updater"
)

func main() {
	curr := updater.CheckCurrentStoredVersion("../log")
	fmt.Println(curr)

	version := updater.CheckCurrentGithubParent("drbh/github-commit-updater")
	fmt.Println(version)

	shouldUpdate := updater.CompareStoredVerionAndGihubVersion(
		"drbh/github-commit-updater", "../log")
	fmt.Println(shouldUpdate)
}
