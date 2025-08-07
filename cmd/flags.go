package cmd

import "github.com/spf13/cobra"

const (
	FlagHome      = "home"
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"

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
func registerCommonFlags(cmd *cobra.Command, defaultHome string) {
	cmd.PersistentFlags().String(FlagHome, defaultHome, "set home directory")
	cmd.PersistentFlags().String(FlagLogLevel, "", "log level format (info, debug, warn, error, panic or fatal)")
	cmd.PersistentFlags().String(FlagLogFormat, "auto", "log output format (auto, logfmt, json, or console)")
}
