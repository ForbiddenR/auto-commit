package resource

import (
	"errors"
)

type Builder struct {
	errs []error

	paths []Visitor

	versionFile string
	dockerfile  string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Dockerfile(path string) *Builder {
	b.dockerfile = path
	return b
}

func (b *Builder) VersionFile(path string) *Builder {
	b.versionFile = path
	return b
}

func (b *Builder) Do() *Result {
	r := b.visitorResult()
	if r.err != nil {
		return r
	}

	return r
}

func (b *Builder) visitorResult() *Result {
	if b.errs != nil {
		return &Result{
			err: errors.New(b.errs[0].Error()),
		}
	}

	if len(b.paths) > 0 {
		return b.visitorByPaths()
	}

	return &Result{err: errors.New("no paths")}
}

func (b *Builder) visitorByPaths() *Result {
	result := &Result{}

	visitors := EagerVisitorList(b.paths)

	result.visitor = visitors

	return result
}

func (b *Builder) Param(mode string, message, author, username, email string) *Builder {
	switch mode {
	case "dockerfile":
		b.paths = append(b.paths, &DockerfileVisitor{
			Message:  message,
			Author:   author,
			Username: username,
			Email:    email,
		})
	default:
		b.errs = append(b.errs, errors.New("invalid mode"))
	}
	return b
}
