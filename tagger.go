package prj

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shabbyrobe/golib/bytescan"
)

const tagFileName = ".prjtags"

var validTag = regexp.MustCompile(`^([^#\-\s]+)$`)

type fileTagger struct {
	file string
}

var _ Tagger = &fileTagger{}

func fileTaggerFromDir(dir string) *fileTagger {
	return &fileTagger{
		file: filepath.Join(dir, tagFileName),
	}
}

func (t *fileTagger) Tags() ([]string, error) {
	bts, err := ioutil.ReadFile(t.file)
	if err != nil {
		return nil, err
	}

	var tags []string
	var seen = map[string]struct{}{}
	scn := bytescan.NewScanner(bts)
	line := 0
	for scn.Scan() {
		line++
		tag := strings.TrimSpace(scn.Text())
		if len(tag) == 0 || tag[0] == '#' {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		if !validTag.MatchString(tag) {
			return nil, fmt.Errorf("prj: invalid tag on line %d of tag file %q: %q", line, t.file, tag)
		}
		tags = append(tags, tag)
	}

	return tags, scn.Err()
}

func (t *fileTagger) Tag(with ...string) error {
	bts, err := ioutil.ReadFile(t.file)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	withIdx := make(map[string]int, len(with))
	for _, w := range with {
		w = strings.TrimSpace(w)
		if !validTag.MatchString(w) {
			return fmt.Errorf("prj: invalid tag %q", w)
		}
		withIdx[w] = 0
	}

	scn := bytescan.NewScanner(bts)
	hasHuman := false
	line := 0
	for scn.Scan() {
		tag := strings.TrimSpace(scn.Text())
		line++
		if len(tag) == 0 || tag[0] == '#' {
			hasHuman = true
			continue
		}
		if !validTag.MatchString(tag) {
			return fmt.Errorf("prj: invalid tag on line %d of tag file %q: %q", line, t.file, tag)
		}
		if n, ok := withIdx[tag]; ok {
			withIdx[tag] = n + 1
		}
	}
	if err := scn.Err(); err != nil {
		return err
	}

	if len(bts) > 0 && bts[len(bts)-1] != '\n' {
		bts = append(bts, '\n')
	}
	if hasHuman {
		if len(bts) > 0 {
			bts = append(bts, '\n')
		}
	}

	for _, w := range with {
		if withIdx[w] <= 0 {
			bts = append(bts, w...)
			bts = append(bts, '\n')
		}
	}

	if err := ioutil.WriteFile(t.file, bts, 0600); err != nil {
		return err
	}

	return nil
}

func (t *fileTagger) Untag(with ...string) error {
	bts, err := ioutil.ReadFile(t.file)
	if err != nil {
		return err
	}

	withIdx := make(map[string]struct{}, len(with))
	for _, w := range with {
		w = strings.TrimSpace(w)
		if !validTag.MatchString(w) {
			return fmt.Errorf("prj: invalid tag %q", w)
		}
		withIdx[w] = struct{}{}
	}

	scn := bytescan.NewScanner(bts)

	var out bytes.Buffer
	var line = 0
	for scn.Scan() {
		line++
		tag := strings.TrimSpace(scn.Text())
		if !validTag.MatchString(tag) {
			return fmt.Errorf("prj: invalid tag on line %d of tag file %q: %q", line, t.file, tag)
		}
		if len(tag) == 0 || tag[0] == '#' {
			out.WriteString(tag)
			out.WriteByte('\n')

		} else if _, ok := withIdx[tag]; !ok {
			out.WriteString(tag)
			out.WriteByte('\n')
		}
	}

	if err := scn.Err(); err != nil {
		return err
	}
	if err := ioutil.WriteFile(t.file, out.Bytes(), 0600); err != nil {
		return err
	}

	return nil
}
