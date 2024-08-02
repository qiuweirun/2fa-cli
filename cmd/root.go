/*
Copyright © 2024 NAME HERE <10231021@qq.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	id      int64
	plat    string
	account string
	secret  string
	issuer  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "2fa",
	Version: "dev@0.1.1",
	Short:   "Two-factor authentication (2FA) verfity application",
	Long: `A CLI application to show your Two-factor authentication (2FA) codes, help developers to verification.
use Time-based one-time password (TOTP) algorithm to generates a one-time password (OTP), Authentication code automatically refreshed every second.
Your data is all encrypted storage in local disk, and without any internet connection!!!
examples and usage that application. For example:

0. Init your application:
	$ 2fa init --pwd=your-possword

1. Display the all 2FA code:
	$ 2fa show
	
2. Add your 2FA account to application:
	$ 2fa add --plat=GitHub --account=qiuweirun --secret=Z7OV*********** --issuer=`,
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.2fa.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
