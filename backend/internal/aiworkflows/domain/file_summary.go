package domain

import (
	"errors"
	"strings"
)

// FileSummary is an immutable value-object pairing a repository file
// path with the LLM-produced summary text for that file.
type FileSummary struct {
	filename string
	summary  string
}

// NewFileSummary constructs a FileSummary. An empty filename is rejected;
// an empty summary is allowed (the LLM may produce no useful output for
// some files and the workflow records that fact rather than dropping it).
func NewFileSummary(filename, summary string) (FileSummary, error) {
	if strings.TrimSpace(filename) == "" {
		return FileSummary{}, errors.New("file summary filename must not be empty")
	}
	return FileSummary{filename: filename, summary: summary}, nil
}

func (f FileSummary) Filename() string { return f.filename }
func (f FileSummary) Summary() string  { return f.summary }
