/*
Copyright © 2024 NAME HERE <10231021@qq.com>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/user"
	"syscall"

	"github.com/qiuweirun/2fa/cmd/consts"
	"github.com/qiuweirun/2fa/cmd/setting"
	"github.com/qiuweirun/2fa/cmd/utils"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	id      int64
	plat    string
	account string
	secret  string
	issuer  string
)

var (
	pwd    string
	salt   string
	dbFile = utils.SessionPath() + string(os.PathSeparator) + consts.DB_FILE
	Conf   *setting.Conf
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "2fa",
	Version: "dev@0.1.1",
	Short:   "Two-factor authentication (2FA) verify application",
	Long: `A CLI application to show your Two-factor authentication (2FA) codes, help developers to verification.
use Time-based one-time password (TOTP) algorithm to generates a one-time password (OTP), Authentication code automatically refreshed every second.
Your data is all encrypted storage in local disk, and without any internet connection!!!
For example:

0. Init your application:
	$ 2fa init --pwd=your-possword

1. Display the all 2FA code:
	$ 2fa show
	
2. Add your 2FA account to application:
	$ 2fa add --plat=GitHub --account=qiuweirun --secret=Z7OV*********** --issuer=`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !utils.CheckFileExist(dbFile) {
			log.Fatal("You should run init commond first!")
		}

		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Fatal("Connect DB Err. " + err.Error())
		}
		defer db.Close()

		// single user record
		row := db.QueryRow("select password,salt from " + consts.TABLE_SYSTEM_NAME + " where id = 1")
		err = row.Scan(&pwd, &salt)
		if err != nil || len(pwd) <= 0 || len(salt) <= 0 {
			log.Fatal("You should run init commond first!", err)
		}

		systemUser, _ := user.Current()
		// check user session
		Conf = setting.NewConf()
		if !Conf.IsVaildSession(pwd) {
			fmt.Printf("Hi '\x1b[1m%v\x1b[0m'! Please login.\n", systemUser.Username)
			loggedIn := false
			for i := 0; i < 3; i++ {
				fmt.Print("Your Password: ")
				password, err := terminal.ReadPassword(int(syscall.Stdin))
				if err == nil && utils.GetMd5(string(password)+salt) == pwd {
					loggedIn = true
					break
				} else {
					fmt.Println("\nPassword incorrect!")
				}
			}

			defaultLifeTime := 72 // unit hours
			if loggedIn {
				fmt.Print("\nLogin expiration(default " + fmt.Sprint(defaultLifeTime) + " hours):")
				var input rune
				_, err := fmt.Fscanf(os.Stdin, "%c", &input)
				if err != nil {
					log.Fatalf("read input err:", err)
					os.Exit(0)
				}
				fmt.Println()
				// todo~~~
				if input != '\n' && input != '\r' {
					defaultLifeTime = int(input)
				}

				// set session
				if !Conf.SetSession(defaultLifeTime, pwd) {
					fmt.Println("System err~")
					os.Exit(0)
				}
			} else {
				fmt.Println("Please try again!")
				os.Exit(0)
			}
		}
	},
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
