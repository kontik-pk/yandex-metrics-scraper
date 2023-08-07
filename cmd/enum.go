package main

import "fmt"

const (
	buildVersion string = ""
	buildDate    string = ""
	buildCommit  string = ""
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
}
