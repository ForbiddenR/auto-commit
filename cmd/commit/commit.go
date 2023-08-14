package commit

import (
	cmdutil "github.com/ForbiddenR/auto-commit/cmd/util"
	"github.com/spf13/cobra"
)

type CommitOptions struct {
	username string
	email    string
	message  string
	author   string
	version  string
	mode     string
	flag     int
}

func NewCommitOptions() *CommitOptions {
	return &CommitOptions{
		username: cmdutil.GetVariableFromGit("user.name"),
		email:    cmdutil.GetVariableFromGit("user.email"),
	}
}

func NewCmdCommit(f cmdutil.Factory) *cobra.Command {
	o := NewCommitOptions()

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit changes",
		Long:  `Commit changes to the repository. If you have multiple changes, separate them with a comma.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(f)
		},
	}

	cmd.Flags().StringVarP(&o.message, "message", "m", "", "Motifiying message")
	cmd.Flags().StringVarP(&o.author, "author", "a", "", "Author")
	cmd.Flags().StringVarP(&o.version, "version", "v", "", "Version")
	cmd.Flags().StringVarP(&o.mode, "mode", "d", "dockerfile", "Mode")
	cmd.Flags().IntVarP(&o.flag, "flat", "f", 0, "Flat")

	return cmd
}

func (o *CommitOptions) Run(f cmdutil.Factory) error {
	r := f.NewBuilder().
		Param(o.mode, o.message, o.author, o.username, o.email, o.flag, o.version).
		Do()

	if err := r.Err(); err != nil {
		return err
	}

	return r.Visitor().Visit()
}
