package cmd

import "github.com/spf13/cobra"

const (
	flagHome              = "home"
	flagLogLevel          = "log-level"
	flagLogFormat         = "log-format"
	flagFile              = "file"
	flagPrivateKey        = "private-key"
	flagMnemonic          = "mnemonic"
	flagCoinType          = "coin-type"
	flagWalletAccount     = "account"
	flagWalletIndex       = "index"
	flagRemoteAddress     = "remote-address"
	flagRemoteUrl         = "remote-url"
	flagMetricsListenAddr = "metrics-listen-addr"
	flagTunnelCreator     = "tunnel-creator"
	flagTunnelIds         = "tunnel-ids"
	flagForce             = "force"
)

// registerCommonFlags registers the common flags for the command.
func registerCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(flagHome, defaultHome, "set home directory")
	cmd.PersistentFlags().String(flagLogLevel, "", "log level format (info, debug, warn, error, panic or fatal)")
	cmd.PersistentFlags().String(flagLogFormat, "auto", "log output format (auto, logfmt, json, or console)")
}
