package cmd

import (
	"github.com/spf13/viper"
	"github.com/spf13/cobra"	
)
const (
	flagHome = "home"
	flagFile = "file"
)

func configInitFlags(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	fileFlag(v, cmd)
	return cmd
}

func fileFlag(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringP(flagFile, "f", "", "fetch toml data from specified file")
	if err := v.BindPFlag(flagFile, cmd.Flags().Lookup(flagFile)); err != nil {
		panic(err)
	}
	return cmd
}