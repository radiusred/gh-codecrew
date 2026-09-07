package tracker

import "strings"

// NormalizeLineEndings rewrites CRLF as LF. Every record the tracker reads
// is read line by line — the `(?m)` scans anchor with `$`, which in Go
// matches only before `\n`, and the paragraph split looks for a blank line
// — so a body GitHub's web editor saved with CRLF defeated all of them: an
// `## Adopts` section that yields its refs under LF yielded none under
// CRLF, and `task finish` would have closed nothing (#296).
//
// It is applied in two places, and needs both:
//
//   - Where a body enters the package, in the GitHub-backed readers that
//     return one (IssueBody, Comments). That is the boundary, and it is
//     where a CRLF body actually arrives.
//   - At every exported scanner's entry, because Tracker is an interface:
//     a string crossing that seam carries no promise about its line
//     endings, and a scanner reached with a body from another backend, a
//     fake tracker or a caller's own hand must read it the same way. The
//     second pass is free when there is nothing to replace — strings.Replace
//     returns its input unchanged when it finds no match — so it costs a
//     scan, not an allocation.
//
// Nothing below those two layers repeats it: line endings are not part of
// the record grammar (SPEC §4), and one rule wants one place to hold it.
func NormalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// StripCode blanks Markdown code out of record text in all three of its
// forms: inline spans (a backtick run closed by a run of the same length,
// as CommonMark reads them), fenced blocks (a line opening with three or
// more backticks or tildes, closed by a fence of the same character at
// least as long) and indented blocks (a run of lines indented four columns
// or more, opening wherever it would not interrupt a paragraph). An
// unclosed fence runs to the end of the text; an unclosed backtick run is
// literal text. A span is replaced with a space so the words on either side
// do not fuse; a block's lines are dropped.
//
// One rule, one implementation: the citation walk (milestone evidence) and
// the verdict scan (ParseVerdicts) both read a comment through it, so a URL
// or a verdict quoted in code is content in both — the record-reading rule
// and the citation-reading rule are the same rule (M13-R6). It lived in
// internal/cli beside the citation walk until the verdict scan needed it.
func StripCode(text string) string {
	var out strings.Builder
	var fence string  // the opening fence of the block being skipped
	indented := false // inside an indented code block
	// canOpen: the line just read is one no paragraph can hold, so an
	// indented block may open on the next. CommonMark 4.4 restricts
	// indented code in one way only — it may not interrupt a paragraph —
	// so a heading, a thematic break or a fence needs no blank line after
	// it, and the start of the text opens like a blank line (#288).
	canOpen := true
	for _, line := range strings.SplitAfter(text, "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		if fence != "" {
			if strings.HasPrefix(trimmed, fence) && strings.Trim(trimmed, fence[:1]+" \t\r\n") == "" {
				fence = ""
			}
			canOpen = true // no line of a fenced block is a paragraph's
			continue
		}
		blank := strings.TrimRight(line, " \t\r\n") == ""
		if indented {
			// Blank lines inside the run belong to the block. Dropping them
			// costs at most an empty line: whatever opened the block — a
			// blank line, a heading, a thematic break, a fence — was
			// written before it, so the newline parting what follows from
			// what precedes survives, and nothing downstream reads
			// paragraph structure out of this output.
			if blank || indentWidth(line) >= 4 {
				continue
			}
			indented = false
		}
		// Tested before the fence, as CommonMark orders them: four columns
		// where a block may open is an indented block whatever it holds.
		if !blank && canOpen && indentWidth(line) >= 4 {
			indented = true
			continue
		}
		if f := fenceOpener(trimmed); f != "" {
			fence = f
			canOpen = true
			continue
		}
		canOpen = blank || notParagraph(line, trimmed)
		out.WriteString(stripSpans(line))
	}
	return out.String()
}

// notParagraph reports whether a line is one no paragraph can hold, so an
// indented code block may open on the line after it: an ATX heading or a
// thematic break. The scanner recognises the other two such lines itself —
// a blank line, and any line of a fenced block. Four columns of
// indentation or more disqualifies both, as CommonMark disqualifies them:
// at four columns a line is code where a block may open, and a paragraph's
// lazy continuation where one may not, but never a heading or a break.
func notParagraph(line, trimmed string) bool {
	if indentWidth(line) >= 4 {
		return false
	}
	return atxHeading(trimmed) || thematicBreak(trimmed)
}

// atxHeading: one to six #, then a space, a tab, or the line's end.
func atxHeading(s string) bool {
	n := 0
	for n < len(s) && s[n] == '#' {
		n++
	}
	if n == 0 || n > 6 {
		return false
	}
	rest := strings.TrimRight(s[n:], "\r\n")
	return rest == "" || rest[0] == ' ' || rest[0] == '\t'
}

// thematicBreak: three or more of *, - or _, the same character
// throughout, spaces and tabs allowed between them and nothing else on the
// line. A --- closing a paragraph is a setext heading underline rather
// than a break, and opens a block either way: neither is a paragraph.
func thematicBreak(s string) bool {
	s = strings.TrimRight(s, " \t\r\n")
	if s == "" || (s[0] != '*' && s[0] != '-' && s[0] != '_') {
		return false
	}
	c, n := s[0], 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case c:
			n++
		case ' ', '\t':
		default:
			return false
		}
	}
	return n >= 3
}

// indentWidth is a line's indentation in columns, a tab advancing to the
// next multiple of four as CommonMark expands one. Four or more opens an
// indented code block where one may open; a line that continues a paragraph
// or a list item is measured too, but never opens one, because a paragraph
// is already open. The measure is from column 0: a list item's own content
// column is not tracked (#285).
func indentWidth(line string) int {
	n := 0
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case ' ':
			n++
		case '\t':
			n += 4 - n%4
		default:
			return n
		}
		if n >= 4 {
			return n
		}
	}
	return n
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
