package commit

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/spf13/cobra"
)

var (
	User   string
	Email  string
	Author string
)

var message string

func NewCmdCommit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit changes",
		Long:  `Commit changes to the repository. If you have multiple changes, separate them with a comma.`,
		Run: func(cmd *cobra.Command, args []string) {
			// wd is the working directory.
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			dockerfile := path.Join(wd, "/Dockerfile")
			version, err := getTagFromDockerfile(dockerfile)
			if err != nil {
				panic(err)
			}

			mdp := path.Join(wd, "/Version.md")
			vf, new, err := getOrCreateVersionFile(mdp)
			if err != nil {
				panic(err)
			}
			defer vf.Close()

			if !new {
				vf.WriteString("\n")
			}
			vf.WriteString(fmt.Sprintf("### %s\n", version))
			vf.WriteString(fmt.Sprintf("+ Author %s %s\n", Author, time.Now().Format("2006.1.2")))
			for _, v := range strings.Split(message, ",") {
				vf.WriteString(fmt.Sprintf("+ %s\n", v))
			}

			commits := "[ci-build]" + version

			r, err := git.PlainOpen(wd)
			if err != nil {
				panic(err)
			}

			w, err := r.Worktree()
			if err != nil {
				panic(err)
			}

			status, err := w.Status()
			if err != nil {
				panic(err)
			}

			var modified []string
			for k := range status {
				if k == "Dockerfile" || k == "Version.md" || k == "go.mod" || k == "go.sum" {
					continue
				}
				modified = append(modified, k)
			}

			vf.WriteString(fmt.Sprintf("+ Modified: %s\n", strings.Join(modified, ", ")))

			// Add all files to the staging area.
			_, err = w.Add(".")
			if err != nil {
				panic(err)
			}

			// Commit the changes to the local repository.
			commit, err := w.Commit(commits, &git.CommitOptions{
				Author: &object.Signature{
					Name:  User,
					Email: Email,
					When:  time.Now(),
				},
			})

			if err != nil {
				panic(err)
			}

			obj, err := r.CommitObject(commit)
			if err != nil {
				panic(err)
			}
			fmt.Println(obj)
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Motifiying message")

	return cmd
}

func getTagFromDockerfile(dockerfile string) (string, error) {
	df, err := os.Open(dockerfile)
	if err != nil {
		return "", err
	}
	defer df.Close()
	result, err := parser.Parse(df)
	if err != nil {
		return "", err
	}

	var version string
	for _, child := range result.AST.Children {
		if child.Next.Value == "TAG" {
			if child.Next.Next == nil {
				return "", fmt.Errorf("TAG is not set")
			}
			version = strings.Split(child.Next.Next.Value, ":")[1]
			break
		}
	}
	return version, nil
}

func getOrCreateVersionFile(vfp string) (*os.File, bool, error) {
	_, err := os.Stat(vfp)
	if os.IsNotExist(err) {
		mf, err := os.Create(vfp)
		if err != nil {
			return nil, true, err
		}
		return mf, true, nil
	} else {
		mf, err := os.OpenFile(vfp, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, false, err
		}
		return mf, false, nil
	}
}
