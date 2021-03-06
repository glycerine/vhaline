package vhaline

import (
	"fmt"
)

var LAST_GIT_COMMIT_HASH string
var NEAREST_GIT_TAG string
var GIT_BRANCH string
var GO_VERSION string

var LibName = "github.com/glycerine/vhaline"

func GoVersion() string {
	return fmt.Sprintf("%s commit: %s / nearest-git-tag: %s / branch: %s / go version: %s\n",
		LibName, LAST_GIT_COMMIT_HASH, NEAREST_GIT_TAG, GIT_BRANCH, GO_VERSION)
}
