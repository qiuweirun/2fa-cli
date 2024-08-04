/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/qiuweirun/2fa/cmd/consts"
	"github.com/qiuweirun/2fa/cmd/utils"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add your 2fa account",
	Long: `Add your 2fa account, your data will be encrypted storage in local disk and likely contains examples
and usage command. For example:

$ 2fa add --plat=GitHub --account=qiuweirun --secret=Z7OV*********** --issuer=`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查参数是否完整了
		if len(plat) <= 0 {
			log.Fatal("args `plat` miss!")
		}

		if len(account) <= 0 {
			log.Fatal("args `account` miss!")
		}

		if len(secret) <= 0 {
			log.Fatal("args `secret` miss!")
		}

		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Fatal("Connect DB Err. " + err.Error())
		}
		defer db.Close()

		// plat + account看看存在了没有
		var id int64
		stmt, _ := db.Prepare("select id from " + consts.TABLE_ACCOUNT_NAME + " where plat=? and account=?")
		row := stmt.QueryRow(plat, account)
		row.Scan(&id)
		if id > 0 {
			log.Fatal("plat:" + plat + ", account:" + account + " is already exist!")
		}

		now := time.Now()
		created := now.Format("2006-01-02 15:04:05")
		encryptPassword := utils.AesEncryptGCM(secret, pwd)
		stmt, _ = db.Prepare("INSERT INTO " + consts.TABLE_ACCOUNT_NAME + "(plat, account, secret, issuer, created) values(?,?,?,?,?)")
		res, err := stmt.Exec(plat, account, encryptPassword, issuer, created)
		if err != nil {
			log.Fatal("plat:"+plat+", account:"+account+" insert data err!", err)
		}

		id, err = res.LastInsertId()
		if err != nil {
			log.Fatal("plat:"+plat+", account:"+account+" insert data err!", err)
		}
		fmt.Println("Add successful!(" + fmt.Sprint(id) + ")")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&plat, "plat", "p", "", "your account platform name，example: GitHub……")
	addCmd.Flags().StringVarP(&account, "account", "a", "", "your 2fa account name")
	addCmd.Flags().StringVarP(&secret, "secret", "s", "", "your 2fa account secret")
	addCmd.Flags().StringVar(&issuer, "issuer", "", "your 2fa account issuer tag")
}
