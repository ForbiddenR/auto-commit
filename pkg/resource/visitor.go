package resource

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type EagerVisitorList []Visitor

func (l EagerVisitorList) Visit() error {
	var errs []error
	for i := range l {
		err := l[i].Visit()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.New(errs[0].Error())
}

type DockerfileVisitor struct {
	*VersionFileVisitor
	Message  string
	Author   string
	Username string
	Email    string
}

func (v *DockerfileVisitor) Visit() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	dockerfile := path.Join(wd, "/Dockerfile")
	df, err := os.Open(dockerfile)
	if err != nil {
		return err
	}
	defer df.Close()
	result, err := parser.Parse(df)
	if err != nil {
		return err
	}

	var version string
	for _, child := range result.AST.Children {
		if child.Next.Value == "TAG" {
			if child.Next.Next == nil {
				return fmt.Errorf("TAG is not set")
			}
			version = strings.Split(child.Next.Next.Value, ":")[1]
			break
		}
	}

	v.VersionFileVisitor = &VersionFileVisitor{
		GitVisitor: &GitVisitor{},
		Wd:         wd,
		Version:    version,
		Message:    v.Message,
		Author:     v.Author,
		Username:   v.Username,
		Email:      v.Email,
	}
	return v.VersionFileVisitor.Visit()
}

type VersionFileVisitor struct {
	*GitVisitor
	Wd       string
	Version  string
	Message  string
	Author   string
	Username string
	Email    string
}

func (v *VersionFileVisitor) Visit() error {
	vfp := path.Join(v.Wd, "/Version.md")
	_, err := os.Stat(vfp)
	var mf *os.File
	var new bool
	if os.IsNotExist(err) {
		mf, err = os.Create(vfp)
		if err != nil {
			return err
		}
	} else {
		mf, err = os.OpenFile(vfp, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	defer mf.Close()

	if !new {
		mf.WriteString("\n")
	}
	mf.WriteString(fmt.Sprintf("### %s\n", v.Version))
	var author string
	if v.Author != "" {
		author = v.Author
	} else {
		author = v.Username
	}
	mf.WriteString(fmt.Sprintf("+ Author %s %s\n", author, time.Now().Format("2006.1.2")))
	for _, v := range strings.Split(v.Message, ",") {
		mf.WriteString(fmt.Sprintf("+ %s\n", v))
	}
	v.GitVisitor = &GitVisitor{
		VfFd:     mf,
		Wd:       v.Wd,
		Version:  v.Version,
		Username: v.Username,
		Email:    v.Email,
	}
	return v.GitVisitor.Visit()
}

type GitVisitor struct {
	*ModifyVisitor
	VfFd     *os.File
	Wd       string
	Version  string
	Username string
	Email    string
}

func (v *GitVisitor) Visit() error {
	commits := fmt.Sprintf("[ci-build]%s", v.Version)

	r, err := git.PlainOpen(v.Wd)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	status, err := w.Status()
	if err != nil {
		return err
	}

	var modified []string

	for k := range status {
		if k == "Dockerfile" || k == "Version.md" || k == "go.mod" || k == "go.sum" {
			continue
		}
		modified = append(modified, k)
	}

	v.ModifyVisitor = &ModifyVisitor{
		Modifiled: modified,
		VfFd:      v.VfFd,
	}

	err = v.ModifyVisitor.Visit()
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	commit, err := w.Commit(commits, &git.CommitOptions{
		Author: &object.Signature{
			Name:  v.Username,
			Email: v.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	_, err = r.CommitObject(commit)
	if err != nil {
		return err
	}

	return nil
}

type ModifyVisitor struct {
	Modifiled []string
	VfFd      *os.File
}

func (v *ModifyVisitor) Visit() error {
	_, err := v.VfFd.WriteString(fmt.Sprintf("+ Modified: %s\n", strings.Join(v.Modifiled, ", ")))
	return err
}
