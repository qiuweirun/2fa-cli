/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
doc: https://github.com/pquerna/otp/tree/master https://github.com/xlzd/gotp
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/qiuweirun/2fa/cmd/consts"
	"github.com/qiuweirun/2fa/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/xlzd/gotp"
)

type totp struct {
	plat    string
	account string
	secret  string
	issuer  string
	code    string
	delay   int64
}

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show your 2fa one-time passwords",
	Long: `Show your 2fa one-time passwords. For example:

1. Display the all 2FA code:
	$ 2fa show`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Fatal("Connect DB Err. " + err.Error())
		}
		defer db.Close()

		stmt, _ := db.Prepare("select plat, account, secret, issuer from " + consts.TABLE_ACCOUNT_NAME + " where 1")
		rows, err := stmt.Query()
		if err != nil {
			log.Fatal("Select Table " + consts.TABLE_SYSTEM_NAME + " Err. " + err.Error())
		}
		defer rows.Close()
		list := make([]totp, 0)
		for rows.Next() {
			var plat, account, secret, issuer string
			err = rows.Scan(&plat, &account, &secret, &issuer)
			if err != nil {
				log.Fatal("Select Table " + consts.TABLE_SYSTEM_NAME + " Err. " + err.Error())
			}
			decryptPwd := utils.AesDecryptGCM(secret, pwd)
			list = append(list, totp{
				plat:    plat,
				account: account,
				secret:  decryptPwd,
				issuer:  issuer,
			})
		}
		err = rows.Err()
		if err != nil {
			log.Fatal("Select Table " + consts.TABLE_SYSTEM_NAME + " Err. " + err.Error())
		}

		if len(list) <= 0 {
			fmt.Println("No anything! please add your account at first!")
			return
		}

		fmt.Println("Press 'Ctrl'+'C' to quit: ")
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		// output right now!
		renderData(toCalculTOTPCode(list), true)
		// after 1s to refresh
		timer := time.NewTimer(1 * time.Second)
		for {
			select {
			case <-c:
				fmt.Println("Bye-Bye! Auto-Expire on: \033[1m" + Conf.GetSessionExpireTime() + "\033[0m")
				return
			case <-timer.C:
				renderData(toCalculTOTPCode(list), false)
			}
			timer.Reset(1 * time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().StringVarP(&account, "account", "a", "", "Specific one account to show one-time passwords.")
	showCmd.Flags().Int64VarP(&id, "id", "i", 0, "Specific one to show one-time passwords.")
}

// toCalculTOTPCode 计算所有账号的code和剩余有效时间
func toCalculTOTPCode(list []totp) []totp {
	for k, v := range list {
		t := gotp.NewDefaultTOTP(v.secret)
		t.ProvisioningUri(v.account, v.issuer)
		code, expireTime := t.NowWithExpiration()
		list[k].code = code
		list[k].delay = expireTime - time.Now().Unix()
	}
	return list
}

// renderTable
func renderData(list []totp, first bool) {
	if !first {
		counter := len(list) + 1
		fmt.Print("\033[" + fmt.Sprint(counter) + "A")
	}

	fmt.Printf("\033[30m\033[47m%-*s | %-*s | %-*s | %-*s | %-*s\033[0m\n", 10, "No#", 10, "Plat", 40, "Account", 10, "Code", 9, "remain(S)")
	for k, v := range list {
		color := ""
		end := "\033[0m"
		if v.delay <= 2 {
			color = "\033[41m\033[37m"
			end = "\033[7m"
		} else if v.delay <= 5 {
			color = "\033[41m\033[37m"
		} else if v.delay <= 10 {
			color = "\033[47m\033[31m"
		}
		fmt.Printf(color+"%-*s | %-*s | %-*s | %-*s | %-*s"+end+"\n", 10, fmt.Sprint(k+1), 10, v.plat, 40, v.account, 10, v.code, 9, fmt.Sprint(v.delay))
	}
}
