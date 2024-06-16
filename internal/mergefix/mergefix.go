package mergefix

import (
	"bytes"
	"errors"
	"io"
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

func FixConflicts(r io.Reader) ([]byte, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if !hasConflictMarkers(buf) {
		return nil, ErrorNoConflicts
	}

	lines := bytes.Split(buf, []byte("\n"))
	var out []byte

	// 0: normal line
	// 1: inside parent conflict marker
	// 2: inside ours conflict marker
	// 3: inside theirs conflict marker
	var state int = 0
	for _, line := range lines {
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
				// dir := bytes.Split(line, []byte(" "))[0]
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
