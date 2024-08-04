package setting

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/qiuweirun/2fa/cmd/consts"
	"github.com/qiuweirun/2fa/cmd/utils"
	"gopkg.in/ini.v1"
)

type Conf struct {
	User       string
	LifeTime   int // hours
	VerifyTime time.Time
	Token      string
	pwd        string
}

var (
	sessionFile = utils.SessionPath() + string(os.PathSeparator) + consts.SESSION_FILE
	timeFormat  = "2006-01-02 15:04:05"
)

func NewConf() *Conf {
	return &Conf{}
}

// IsVaildSession
func (u *Conf) IsVaildSession(pwd string) bool {
	if !utils.CheckFileExist(sessionFile) {
		return false
	}

	cfg, err := ini.Load(sessionFile)
	if err != nil {
		fmt.Println("Error loading INI file:", err)
		return false
	}

	section, err := cfg.GetSection("SESSION")
	if err != nil {
		fmt.Println("Error getting section:", err)
		return false
	}

	u.User = section.Key("user").Value()
	u.LifeTime, err = section.Key("life_time").Int()
	if err != nil {
		return false
	}
	u.VerifyTime, err = section.Key("verify_time").TimeFormat(timeFormat)
	u.Token = section.Key("_token_").Value()
	if err != nil {
		fmt.Println("Error read the section:", err)
		return false
	}

	systemUser, _ := user.Current()
	expireTime := u.VerifyTime.Add(time.Duration(u.LifeTime * 60 * 60))
	var sb strings.Builder
	sb.WriteString(u.User)
	sb.WriteString(fmt.Sprintf("%d", u.LifeTime))
	sb.WriteString(u.VerifyTime.Format(timeFormat))
	sb.WriteString(pwd)
	hash := utils.GetMd5(sb.String())
	if systemUser.Username != u.User || u.LifeTime <= 0 || time.Now().After(expireTime) || hash != u.Token {
		return false
	}
	return true
}

// SetSession
func (u *Conf) SetSession(life_time int, pwd string) bool {
	cfg := ini.Empty()

	// 创建一个新节
	section, err := cfg.NewSection("SESSION")
	if err != nil {
		fmt.Println("Error creating section:", err)
		return false
	}

	now := time.Now().Format(timeFormat)
	systemUser, _ := user.Current()
	section.Key("user").SetValue(systemUser.Username)
	section.Key("life_time").SetValue(fmt.Sprintf("%d", life_time))
	section.Key("verify_time").SetValue(now)
	var sb strings.Builder
	sb.WriteString(u.User)
	sb.WriteString(fmt.Sprintf("%d", u.LifeTime))
	sb.WriteString(now)
	sb.WriteString(pwd)
	hash := utils.GetMd5(sb.String())
	section.Key("_token_").SetValue(hash)
	err = cfg.SaveTo(sessionFile)
	if err != nil {
		fmt.Println("Error saving session:", err)
		return false
	}
	return true
}

// GetSessionExpireTime
func (u *Conf) GetSessionExpireTime() time.Time {
	return u.VerifyTime.Add(time.Duration(u.LifeTime) * time.Hour)
}
