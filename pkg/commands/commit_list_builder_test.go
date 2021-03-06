package commands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

// NewDummyCommitListBuilder creates a new dummy CommitListBuilder for testing
func NewDummyCommitListBuilder() *CommitListBuilder {
	osCommand := NewDummyOSCommand()

	return &CommitListBuilder{
		Log:        NewDummyLog(),
		GitCommand: NewDummyGitCommandWithOSCommand(osCommand),
		OSCommand:  osCommand,
		Tr:         i18n.NewLocalizer(NewDummyLog()),
	}
}

// TestCommitListBuilderGetMergeBase is a function.
func TestCommitListBuilderGetMergeBase(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func(string, error)
	}

	scenarios := []scenario{
		{
			"swallows an error if the call to merge-base returns an error",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)

				switch args[0] {
				case "symbolic-ref":
					assert.EqualValues(t, []string{"symbolic-ref", "--short", "HEAD"}, args)
					return exec.Command("echo", "master")
				case "merge-base":
					assert.EqualValues(t, []string{"merge-base", "HEAD", "master"}, args)
					return exec.Command("test")
				}
				return nil
			},
			func(output string, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "", output)
			},
		},
		{
			"returns the commit when master",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)

				switch args[0] {
				case "symbolic-ref":
					assert.EqualValues(t, []string{"symbolic-ref", "--short", "HEAD"}, args)
					return exec.Command("echo", "master")
				case "merge-base":
					assert.EqualValues(t, []string{"merge-base", "HEAD", "master"}, args)
					return exec.Command("echo", "blah")
				}
				return nil
			},
			func(output string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "blah", output)
			},
		},
		{
			"checks against develop when a feature branch",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)

				switch args[0] {
				case "symbolic-ref":
					assert.EqualValues(t, []string{"symbolic-ref", "--short", "HEAD"}, args)
					return exec.Command("echo", "feature/test")
				case "merge-base":
					assert.EqualValues(t, []string{"merge-base", "HEAD", "develop"}, args)
					return exec.Command("echo", "blah")
				}
				return nil
			},
			func(output string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "blah", output)
			},
		},
		{
			"bubbles up error if there is one",
			func(cmd string, args ...string) *exec.Cmd {
				return exec.Command("test")
			},
			func(output string, err error) {
				assert.Error(t, err)
				assert.Equal(t, "", output)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			c := NewDummyCommitListBuilder()
			c.OSCommand.SetCommand(s.command)
			s.test(c.getMergeBase("HEAD"))
		})
	}
}
