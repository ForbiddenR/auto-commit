package commit

import (
	"os/exec"

	cmdutil "github.com/ForbiddenR/auto-commit/cmd/util"
	"github.com/spf13/cobra"
)

type CommitOptions struct {
	username string
	email    string
	message  string
	author   string
}

func NewCommitOptions() *CommitOptions {
	out, err := exec.Command("git", "config", "--get", "user.name").Output()
	if err != nil {
		panic(err)
	}
	username := string(out)[:len(out)-1]

	out, err = exec.Command("git", "config", "--get", "user.email").Output()
	if err != nil {
		panic(err)
	}
	email := string(out)[:len(out)-1]
	return &CommitOptions{
		username: username,
		email:    email,
	}
}

func NewCmdCommit(f cmdutil.Factory) *cobra.Command {
	o := NewCommitOptions()

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit changes",
		Long:  `Commit changes to the repository. If you have multiple changes, separate them with a comma.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// wd is the working directory.
			return o.Run(f)
		},
	}

	cmd.Flags().StringVarP(&o.message, "message", "m", "", "Motifiying message")
	cmd.Flags().StringVarP(&o.author, "author", "a", "", "Author")

	return cmd
}

func (o *CommitOptions) Run(f cmdutil.Factory) error {
	r := f.NewBuilder().
		Param("dockerfile", o.message, o.author, o.username, o.email).
		Do()

	if err := r.Err(); err != nil {
		return err
	}

	return r.Visitor().Visit()
}
