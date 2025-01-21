package relayertest

import (
	"bytes"
	"context"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/bandprotocol/falcon/cmd"
)

// System is a system under test.
type System struct {
	// Temporary directory to be injected as --home argument.
	HomeDir string
}

// NewSystem creates a new system with a home dir associated with a temp dir belonging to t.
//
// The returned System does not store a reference to t;
// some of its methods expect a *testing.T as an argument.
// This allows creating one instance of System to be shared with subtests.
func NewSystem(t *testing.T) *System {
	t.Helper()

	homeDir := t.TempDir()

	return &System{
		HomeDir: homeDir,
	}
}

// RunResult is the stdout and stderr resulting from a call to (*System).Run,
// and any error that was returned.
type RunResult struct {
	Stdout, Stderr bytes.Buffer

	Err error
}

func (s *System) RunWithInput(t *testing.T, args ...string) RunResult {
	rootCmd := cmd.NewRootCmd(zaptest.NewLogger(t))

	rootCmd.SilenceUsage = true

	ctx := context.Background()

	var res RunResult
	rootCmd.SetOut(&res.Stdout)
	rootCmd.SetErr(&res.Stderr)

	// Prepend the system's home directory to any provided args.
	args = append([]string{"--home", s.HomeDir}, args...)
	rootCmd.SetArgs(args)

	res.Err = rootCmd.ExecuteContext(ctx)

	return res
}
