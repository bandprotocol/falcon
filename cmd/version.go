package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
)

var (
	// Version defines the application version (defined at compile time)
	Version = ""
	Commit  = ""
	Dirty   = ""
)

type versionInfo struct {
	Version string `json:"version" yaml:"version"`
	Commit  string `json:"commit"  yaml:"commit"`
	Go      string `json:"go"      yaml:"go"`
}

// VersionCmd returns a command that prints the falcon version information.
func VersionCmd(app *relayer.App) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   fmt.Sprintf("Print the %s version info", app.Name),
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s version
$ %s v`,
			app.Name, app.Name,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			commit := Commit
			if Dirty != "0" {
				commit += " (dirty)"
			}

			verInfo := versionInfo{
				Version: Version,
				Commit:  commit,
				Go:      fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
			}

			bz, err := toml.Marshal(&verInfo)

			fmt.Fprintln(cmd.OutOrStdout(), string(bz))
			return err
		},
	}

	return versionCmd
}
