package tracker

import "strings"

// StripCode blanks Markdown code out of record text: fenced blocks (a line
// opening with three or more backticks or tildes, closed by a fence of the
// same character at least as long) and inline spans (a backtick run closed
// by a run of the same length, as CommonMark reads them). An unclosed fence
// runs to the end of the text; an unclosed backtick run is literal text.
// Replaced with a space so words on either side do not fuse.
//
// One rule, one implementation: the citation walk (milestone evidence) and
// the verdict scan (ParseVerdicts) both read a comment through it, so a URL
// or a verdict quoted in code is content in both — the record-reading rule
// and the citation-reading rule are the same rule (M13-R6). It lived in
// internal/cli beside the citation walk until the verdict scan needed it.
func StripCode(text string) string {
	var out strings.Builder
	var fence string // the opening fence of the block being skipped
	for _, line := range strings.SplitAfter(text, "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		if fence != "" {
			if strings.HasPrefix(trimmed, fence) && strings.Trim(trimmed, fence[:1]+" \t\r\n") == "" {
				fence = ""
			}
			continue
		}
		if f := fenceOpener(trimmed); f != "" {
			fence = f
			continue
		}
		out.WriteString(stripSpans(line))
	}
	return out.String()
}

// fenceOpener returns the fence run that opens a code block on this line
// — three or more backticks or tildes at its start, an info string such as
// a language tag allowed after them — or "" if none. A fence may open on
// the same line as a list marker (`- ```sh`): the item's content starts
// after the marker, so the marker is skipped first.
func fenceOpener(line string) string {
	line = strings.TrimLeft(skipListMarker(line), " \t")
	if line == "" || (line[0] != '`' && line[0] != '~') {
		return ""
	}
	n := 0
	for n < len(line) && line[n] == line[0] {
		n++
	}
	if n < 3 {
		return ""
	}
	return line[:n]
}

// skipListMarker drops a leading bullet (-, *, +) or ordered marker (1.,
// 1)) and the space after it; a line with no marker is returned as is.
func skipListMarker(line string) string {
	rest := line
	if len(rest) > 0 && strings.ContainsRune("-*+", rune(rest[0])) {
		rest = rest[1:]
	} else {
		n := 0
		for n < len(rest) && rest[n] >= '0' && rest[n] <= '9' {
			n++
		}
		if n == 0 || n > 9 || n >= len(rest) || (rest[n] != '.' && rest[n] != ')') {
			return line
		}
		rest = rest[n+1:]
	}
	if rest == "" || (rest[0] != ' ' && rest[0] != '\t') {
		return line
	}
	return rest
}

// stripSpans replaces every closed inline code span on one line with a
// space; a backtick run with no matching closer stays as it is.
func stripSpans(line string) string {
	var out strings.Builder
	for i := 0; i < len(line); {
		if line[i] != '`' {
			out.WriteByte(line[i])
			i++
			continue
		}
		n := 0
		for i+n < len(line) && line[i+n] == '`' {
			n++
		}
		run := line[i : i+n]
		end := -1
		for j := i + n; j < len(line); {
			k := strings.Index(line[j:], run)
			if k < 0 {
				break
			}
			j += k
			m := j
			for m < len(line) && line[m] == '`' {
				m++
			}
			if m-j == n {
				end = j
				break
			}
			j = m
		}
		if end < 0 {
			out.WriteString(run)
			i += n
			continue
		}
		out.WriteByte(' ')
		i = end + n
	}
	return out.String()
}
