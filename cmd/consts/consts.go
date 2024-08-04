package consts

var (
	DB_FILE      = "2fa-cli.db"
	SESSION_FILE = ".2fa_app.ini"

	TABLE_SYSTEM_NAME   = "system"
	TABLE_SYSTEM_STRUCT = `
    CREATE TABLE IF NOT EXISTS system (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        password VARCHAR(256) NOT NULL,
		salt VARCHAR(64) NOT NULL,
        created DATETIME NOT NULL
    );
    `
	TABLE_ACCOUNT_NAME   = "account"
	TABLE_ACCOUNT_STRUCT = `
	CREATE TABLE IF NOT EXISTS account (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        plat VARCHAR(64) NOT NULL,
		account VARCHAR(128) NOT NULL,
		secret VARCHAR(256) NOT NULL,
		issuer VARCHAR(64) NULL,
        created DATETIME NOT NULL
    );
	`
)
