package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	falcon "github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/store"
)

const (
	appName         = "falcon"
	defaultCoinType = 60

	PassphraseEnvKey = "PASSPHRASE"
)

var defaultHome = filepath.Join(os.Getenv("HOME"), ".falcon")

// NewRootCmd returns the root command for falcon.
func NewRootCmd(log *zap.Logger) *cobra.Command {
	passphrase := os.Getenv(PassphraseEnvKey)
	homePath := defaultHome
	app := falcon.NewApp(log, nil, passphrase, nil)

	// RootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   appName,
		Short: "Falcon relays tss tunnel messages from BandChain to destination chains/smart contracts",
		Long: strings.TrimSpace(`This application has:
   1. Configuration Management: Handles the configuration of the program.
   2. Key Management: Supports managing multiple keys across multiple chains.
   3. Transaction Execution: Enables executing transactions on destination chains.
   4. Query Functionality: Facilitates querying data from both source and destination chains.

   NOTE: Most of the commands have aliases that make typing them much quicker 
         (i.e. 'falcon tx', 'falcon q', etc...)`),
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) (err error) {
		// set up store
		app.Store, err = store.NewFileSystem(homePath)
		if err != nil {
			return err
		}

		// load configuration
		app.Config, err = app.Store.GetConfig()
		if err != nil {
			return err
		}

		// retrieve log level from config
		configLogLevel := ""
		if app.Config != nil {
			configLogLevel = app.Config.Global.LogLevel
		}

		// init log object
		app.Log, err = initLogger(configLogLevel)
		if err != nil {
			return err
		}

		return app.Init(rootCmd.Context())
	}

	rootCmd.PersistentPostRun = func(cmd *cobra.Command, _ []string) {
		// Force syncing the logs before exit, if anything is buffered.
		// check error of log.Sync() https://github.com/uber-go/zap/issues/991#issuecomment-962098428
		if err := app.Log.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			fmt.Fprintf(os.Stderr, "failed to sync logs: %v\n", err)
		}
	}

	// Register --home flag
	rootCmd.PersistentFlags().StringVar(&homePath, flagHome, defaultHome, "set home directory")
	if err := viper.BindPFlag(flagHome, rootCmd.PersistentFlags().Lookup(flagHome)); err != nil {
		panic(err)
	}

	// Register --log-format flag
	rootCmd.PersistentFlags().String("log-format", "auto", "log output format (auto, logfmt, json, or console)")
	if err := viper.BindPFlag("log-format", rootCmd.PersistentFlags().Lookup("log-format")); err != nil {
		panic(err)
	}

	// Register --log-level flag
	rootCmd.PersistentFlags().String("log-level", "", "log level format (info, debug, warn, error, panic or fatal)")
	if err := viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		panic(err)
	}

	// Register subcommands
	rootCmd.AddCommand(
		ConfigCmd(app),
		ChainsCmd(app),
		KeysCmd(app),

		lineBreakCommand(),
		TransactionCmd(app),
		QueryCmd(app),
		StartCmd(app),

		lineBreakCommand(),
		VersionCmd(app),
	)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd(nil)
	rootCmd.SilenceUsage = true

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		// Wait for interrupt signal.
		sig := <-sigCh

		// Cancel the context to signal to the rest of the application to shut down.
		cancel()

		// Short delay before printing the received signal message.
		// This should result in cleaner output from non-interactive commands that stop quickly.
		time.Sleep(250 * time.Millisecond)
		fmt.Fprintf(
			os.Stderr,
			"Received signal %v. Attempting clean shutdown. Send interrupt again to force hard shutdown.\n",
			sig,
		)

		// Dump all goroutines on panic, not just the current one.
		debug.SetTraceback("all")

		// Block waiting for a second interrupt or a timeout.
		// The main goroutine ought to finish before either case is reached.
		// But if a case is reached, panic so that we get a non-zero exit and a dump of remaining goroutines.
		select {
		case <-time.After(time.Minute):
			panic(errors.New("falcon did not shut down within one minute of interrupt"))
		case sig := <-sigCh:
			panic(fmt.Errorf("received signal %v; forcing quit", sig))
		}
	}()

	return rootCmd.ExecuteContext(ctx)
}

// lineBreakCommand returns a new instance of the lineBreakCommand every time to avoid
// data races in concurrent tests exercising commands.
func lineBreakCommand() *cobra.Command {
	return &cobra.Command{Run: func(*cobra.Command, []string) {}}
}

// withUsage wraps a PositionalArgs to display usage only when the PositionalArgs
// variant is violated.
func withUsage(inner cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := inner(cmd, args); err != nil {
			cmd.Root().SilenceUsage = false
			cmd.SilenceUsage = false
			return err
		}

		return nil
	}
}
