package types

import "github.com/hhatto/gocloc"

type Report struct {
	Path        string
	Inputgocloc *gocloc.Result
}
