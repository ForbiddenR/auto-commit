package commit

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5"
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
			df, err := os.Open(dockerfile)
			if err != nil {
				panic(err)
			}
			defer df.Close()
			result, err := parser.Parse(df)
			if err != nil {
				panic(err)
			}

			var version string
			for _, child := range result.AST.Children {
				if child.Next.Value == "TAG" {
					if child.Next.Next == nil {
						panic("TAG value is not set")
					}
					version = strings.Split(child.Next.Next.Value, ":")[1]
					break
				}
			}

			mdp := path.Join(wd, "/Version.md")
			var mf *os.File
			var new bool
			_, err = os.Stat(mdp)
			if os.IsNotExist(err) {
				mf, err = os.Create(mdp)
				if err != nil {
					panic(err)
				}
				new = true
			} else {
				mf, err = os.OpenFile(mdp, os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					panic(err)
				}
			}
			defer mf.Close()
			if !new {
				mf.WriteString("\n")
			}
			mf.WriteString(fmt.Sprintf("### %s\n", version))
			mf.WriteString(fmt.Sprintf("+ Author %s %s\n", Author, time.Now().Format("2006.1.2")))
			for _, v := range strings.Split(message, ",") {
				mf.WriteString(fmt.Sprintf("+ %s\n", v))
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
			
			status, err  := w.Status()
			if err != nil {
				panic(err)
			}

			var modified []string
			for k := range status {
				modified = append(modified, k)
			}

			mf.WriteString(fmt.Sprintf("+ Modified: %s\n", strings.Join(modified, ", ")))

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
