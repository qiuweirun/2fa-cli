/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/qiuweirun/2fa/cmd/utils"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout right now!",
	Long: `logout right now! You don't need to wait for the expiration times. For example:

	$ 2fa logout`,
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.CheckFileExist(dbFile) {
			log.Fatal("You should run init commond first!")
		}

		if Conf.Clear() {
			fmt.Println("logout success!")
		} else {
			fmt.Println("logout failure!")
		}
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
