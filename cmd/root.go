package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	falcon "github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/store"
)

const (
	defaultCoinType = 60
	appName         = "falcon"

	EnvPassphrase = "passphrase"
	DbPath        = "db_path"
)

var defaultHome = filepath.Join(os.Getenv("HOME"), ".falcon")

// NewRootCmd returns the root command for falcon.
func NewRootCmd() *cobra.Command {
	// RootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use: appName,
		Short: fmt.Sprintf(
			"%s relays tss tunnel messages from BandChain to destination chains/smart contracts",
			appName,
		),
		Long: fmt.Sprintf(`This application has:
   1. Configuration Management: Handles the configuration of the program.
   2. Key Management: Supports managing multiple keys across multiple chains.
   3. Transaction Execution: Enables executing transactions on destination chains.
   4. Query Functionality: Facilitates querying data from both source and destination chains.

   NOTE: Most of the commands have aliases that make typing them much quicker 
         (i.e. '%s tx', '%s q', etc...)`,
			appName, appName,
		),
	}

	registerCommonFlags(rootCmd, defaultHome)

	ac := &AppCreator{}

	// Register subcommands
	rootCmd.AddCommand(
		ConfigCmd(ac.NewApp, defaultHome),
		ChainsCmd(ac.NewApp, defaultHome),
		KeysCmd(ac.NewApp, defaultHome),

		lineBreakCommand(),
		TransactionCmd(ac.NewApp, defaultHome),
		QueryCmd(ac.NewApp, defaultHome),
		StartCmd(ac.NewApp, defaultHome),

		lineBreakCommand(),
		VersionCmd(defaultHome),
	)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd()
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
			panic(errors.New("program did not shut down within one minute of interrupt"))
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

// createApp creates a new application instance.
func createApp(
	cmd *cobra.Command,
	appCreator falcon.AppCreator,
	defaultHome string,
) (falcon.Application, error) {
	// bind flags to the Context's Viper so we can get flags.
	vp := viper.New()
	if err := vp.BindEnv(EnvPassphrase); err != nil {
		return nil, err
	}

	if err := vp.BindEnv(DbPath); err != nil {
		return nil, err
	}

	if err := vp.BindPFlags(cmd.Flags()); err != nil {
		return nil, err
	}

	home := vp.GetString(FlagHome)
	if home == "" {
		home = defaultHome
	}

	store, err := store.NewFileSystem(home)
	if err != nil {
		return nil, err
	}

	return appCreator(store, vp)
}

// syncLog syncs the log to the specific output at the end of the program.
func syncLog(log logger.Logger) {
	if err := log.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) && !errors.Is(err, syscall.EINVAL) {
		fmt.Fprintf(os.Stderr, "failed to sync logs: %v\n", err)
	}
}
