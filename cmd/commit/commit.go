package commit

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var (
	User  string
	Email string
)

var message string

func NewCmdCommit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit changes",
		Long:  `Commit changes to the repository.`,
		Run: func(cmd *cobra.Command, args []string) {
			// print the local directory.
			fmt.Println("Committing changes...")
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			fmt.Println(wd)

			r, err := git.PlainOpen(wd)
			if err != nil {
				panic(err)
			}

			w, err := r.Worktree()
			if err != nil {
				panic(err)
			}

			// Add all files to the staging area.
			_, err = w.Add(".")
			if err != nil {
				panic(err)
			}

			status, err := w.Status()
			if err != nil {
				panic(err)
			}

			fmt.Println(status)

			// Commit the changes to the local repository.
			commit, err := w.Commit(message, &git.CommitOptions{
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

	cmd.Flags().StringVarP(&message, "message", "m", "", "Commit message")

	return cmd
}
