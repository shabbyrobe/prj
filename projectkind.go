package prj

import (
	"fmt"
	"strings"
)

type ProjectKind int

const (
	ProjectSimple ProjectKind = iota + 1
	ProjectGit
	ProjectHg

	ProjectKindSize int = iota + 1
)

func (k ProjectKind) String() string {
	switch k {
	case ProjectSimple:
		return "prj"
	case ProjectGit:
		return "git"
	case ProjectHg:
		return "hg"
	default:
		return "unknown"
	}
}

func (k ProjectKind) IsValid() bool {
	switch k {
	case ProjectSimple,
		ProjectGit,
		ProjectHg:
		return true
	default:
		return false
	}
}

func (k *ProjectKind) Set(s string) error {
	switch s {
	case "prj":
		*k = ProjectSimple
	case "git":
		*k = ProjectGit
	case "hg":
		*k = ProjectHg
	default:
		return fmt.Errorf("unknown project kind %q", s)
	}
	return nil
}

type ProjectKindSet [ProjectKindSize]bool

func (p *ProjectKindSet) SetAll() {
	for i := ProjectKind(0); i < ProjectKind(ProjectKindSize); i++ {
		if i > 0 {
			p[i] = true
		}
	}
}

func (p *ProjectKindSet) Count() (n int) {
	for i := 1; i < ProjectKindSize; i++ {
		if p[ProjectKind(i)] {
			n++
		}
	}
	return n
}

func (p *ProjectKindSet) Set(s string) error {
	var k ProjectKind
	if err := k.Set(s); err != nil {
		return err
	}
	p[k] = true
	return nil
}

func (p *ProjectKindSet) String() string {
	var bits []string
	for i := 0; i < ProjectKindSize; i++ {
		k := ProjectKind(i)
		if p[k] {
			bits = append(bits, k.String())
		}
	}
	return strings.Join(bits, ",")
}
