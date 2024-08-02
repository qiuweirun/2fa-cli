/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"database/sql"

	"github.com/qiuweirun/2fa/cmd/consts"
	"github.com/qiuweirun/2fa/cmd/utils"
	"github.com/spf13/cobra"

	_ "github.com/mattn/go-sqlite3"
)

var (
	pwd string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init your 2fa app",
	Long: `Init your 2fa app, initial password and DB. For example:

2fa init --pwd=your-possword`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查文件是否存在
		if utils.CheckFileExist(consts.DB_PATH) {
			log.Fatal("Your 2fa app is already initialized.")
		}

		// 检查参数是否完整
		if len(pwd) <= 3 {
			log.Fatal("Your password should be greater than 3 characters.")
		}

		db, err := sql.Open("sqlite3", consts.DB_PATH)
		if err != nil {
			log.Fatal("Connect DB Err. " + err.Error())
		}
		defer db.Close()
		_, err = db.Exec(consts.TABLE_SYSTEM_STRUCT)
		if err != nil {
			os.Remove(consts.DB_PATH)
			log.Fatal("Create Table " + consts.TABLE_SYSTEM_NAME + " Err. " + err.Error())
		}

		now := time.Now()
		uuid := utils.CreateUUID()
		password := utils.GetMd5(pwd + uuid)
		created := now.Format("2006-01-02 15:04:05")

		stmt, err := db.Prepare("INSERT INTO " + consts.TABLE_SYSTEM_NAME + "(password, salt, created) values(?,?,?)")
		if err != nil {
			os.Remove(consts.DB_PATH)
			log.Fatal("Insert Table " + consts.TABLE_SYSTEM_NAME + " Err. " + err.Error())
		}
		res, err := stmt.Exec(password, uuid, created)
		if err != nil {
			os.Remove(consts.DB_PATH)
			log.Fatal("Insert Table " + consts.TABLE_SYSTEM_NAME + " Err. " + err.Error())
		}

		id, err := res.LastInsertId()
		if err != nil {
			os.Remove(consts.DB_PATH)
			log.Fatal("Init System Err. " + err.Error())
		}

		// 创建账号数据表
		_, err = db.Exec(consts.TABLE_ACCOUNT_STRUCT)
		if err != nil {
			os.Remove(consts.DB_PATH)
			log.Fatal("Create Table " + consts.TABLE_ACCOUNT_NAME + " Err. " + err.Error())
		}

		fmt.Println("Init successful, please keep your password!(" + fmt.Sprint(id) + ")")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&pwd, "pwd", "", "Input your possword.")
}
