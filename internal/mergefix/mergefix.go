package mergefix

import (
	"bytes"
	"errors"
	"regexp"
)

var (
	ParentMergeRE = regexp.MustCompile(`\|{7,}`)
	OursMergeRE   = regexp.MustCompile(`\<{7,}`)
	TheirsMergeRE = regexp.MustCompile(`\={7,}`)
	EndMergeRE    = regexp.MustCompile(`\>{7,}`)

	ErrorUnsupportedDirective = errors.New("unsupported directive. Please fix the conflicts manually.")
	ErrorNoConflicts          = errors.New("no conflicts to be fixed")
)

func hasConflictMarkers(b []byte) bool {
	return OursMergeRE.Match(b) && TheirsMergeRE.Match(b) && EndMergeRE.Match(b)
}

// RemoveConflictMarkers reads lines from the given input slice and removes any lines that include git conflict markers.
func RemoveConflictMarkers(data []byte) ([]byte, error) {
	if !hasConflictMarkers(data) {
		return nil, ErrorNoConflicts
	}

	// small state machine to track current line state
	// 0: normal line
	// 1: inside parent conflict marker
	// 2: inside ours conflict marker
	// 3: inside theirs conflict marker
	var state int = 0
	var out []byte
	for line := range bytes.Lines(data) {
		switch {
		case ParentMergeRE.Match(line):
			state = 1
		case OursMergeRE.Match(line):
			state = 2
		case TheirsMergeRE.Match(line):
			state = 3
		case EndMergeRE.Match(line):
			state = 0
		default:
			// report conflicts that should be fixed manually
			if state != 0 && isUnsupportedDirective(line) {
				return nil, ErrorUnsupportedDirective
			}
			out = append(out, line...)
			out = append(out, '\n')
		}
	}

	return out, nil
}

func isUnsupportedDirective(line []byte) bool {
	return bytes.HasPrefix(line, []byte("replace")) || bytes.HasPrefix(line, []byte("exclude"))
}
