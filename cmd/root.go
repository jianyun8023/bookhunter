package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/log"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bookhunter",
	Short: "A downloader for downloading books from internet",
	Long: `You can use this command to download book from these websites

1. Self-hosted talebook websites
2. https://www.sanqiu.mobi
3. Telegram channel`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Download commands.
	rootCmd.AddCommand(talebookCmd)
	rootCmd.AddCommand(sanqiuCmd)
	rootCmd.AddCommand(tianlangCmd)
	rootCmd.AddCommand(telegramCmd)
	rootCmd.AddCommand(sobooksCmd)
	rootCmd.AddCommand(youyiduCmd)
	rootCmd.AddCommand(k12Cmd)

	// Tool commands.
	rootCmd.AddCommand(aliyunCmd)
	rootCmd.AddCommand(versionCmd)

	persistentFlags := rootCmd.PersistentFlags()

	// Common flags.
	persistentFlags.StringVarP(&flags.ConfigRoot, "config", "c", flags.ConfigRoot, "The config path for bookhunter")
	persistentFlags.StringVar(&flags.Proxy, "proxy", flags.Proxy, "The request proxy")
	persistentFlags.StringVarP(&flags.UserAgent, "user-agent", "a", flags.UserAgent, "The request user-agent")
	persistentFlags.BoolVar(&log.EnableDebug, "verbose", false, "Print all the logs for debugging")
}
