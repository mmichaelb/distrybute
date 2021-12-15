package util

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

var (
	GitBranch    string
	GitTag       string
	GitCommitSha string
	Version      = fmt.Sprintf("%s/%s/%s", GitBranch, GitTag, GitCommitSha)
)

var Authors = []*cli.Author{{Name: "mmichaelb", Email: "me@mmichaelb.pw"}}
var GeneralApp = cli.App{
	Authors: Authors,
	Version: Version,
}
