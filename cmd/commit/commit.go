package commit

import (
	"fmt"
	"os"

	"github.com/ForbiddenR/auto-commit/pkg/parser"
	"github.com/spf13/cobra"
)

type CommitOptions struct {
	message string
	prefix  string
}

func NewCommitOptions() *CommitOptions {
	return &CommitOptions{}
}

func NewCmdCommit() *cobra.Command {
	o := NewCommitOptions()
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit changes",
		Long:  `Commit changes to the repository. If you have multiple changes, separate them with spaces.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
		SilenceUsage: true,
	}
	cmd.Flags().StringVarP(&o.message, "message", "m", "", "Motifiying message")
	cmd.Flags().StringVarP(&o.prefix, "prefix", "p", "", "Prefix")
	return cmd
}

func (o *CommitOptions) Run() error {
	df, err := os.Open("Dockerfile")
	if err != nil {
		return err
	}
	defer df.Close()
	dp := parser.DockerfileParser{}
	err = dp.Parse(df)
	if err != nil {
		return err
	}
	p := parser.NewVersionParser(dp.String())
	err = func() error {
		file, err := os.Open("Version.md")
		if err != nil {
			return err
		}
		defer file.Close()
		err = p.Parse(file)
		if err != nil {
			return err
		}
		err = p.AddRecord(o.message)
		return err
	}()
	if err != nil {
		return err
	}
	file, err := os.OpenFile("Version.md", os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(p.String())
	if err != nil {
		return err
	}
	fmt.Printf("%s %s %s\n", o.prefix, dp.String(), o.message)
	return nil
}
