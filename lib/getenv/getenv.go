package getenv

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// 優先順位
// 環境変数 < ~/.pnzr < .pnzr < コマンドライン引数

func init() {
	godotenv.Load("~/.pnzr")
	godotenv.Load(".pnzr")
}

func convertStringToBoolean(s string) bool {
	s = strings.ToLower(s)
	switch s {
	case "true", "t", "1":
		return true
	default:
		return false
	}
}

func Bool(key string, def ...bool) bool {
	var d bool
	if len(def) != 0 {
		d = def[0]
	} else {
		d = false
	}
	v := os.Getenv(key)
	if v == "" {
		return d
	}
	return convertStringToBoolean(v)
}

func String(key string, def ...string) string {
	var d string
	if len(def) != 0 {
		d = def[0]
	}
	v := os.Getenv(key)
	if v == "" {
		return d
	}
	return v
}
