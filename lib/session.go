package lib

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

// GetSession アクセスキーとリージョンを取得し, session.Session を返す
//
// # アクセスキー取得方法
//
// 以下の優先順で取得する
//
// 1. アクセスキーID
//     a. コマンドライン引数 (aws-access-key-id)
//     b. 環境変数 (AWS_ACCESS_KEY_ID)
// 2. awsコンフィグ
//     a. コマンドライン引数 (profile)
//     b. 環境変数 (AWS_PROFILE_NAME)
//     c. デフォルト値 (default)
//
// # リージョン取得方法
//
// 以下の優先順で取得する
//
// 1. コマンドライン引数 (region)
// 2. 環境変数 (AWS_REGION)
//
func GetSession(args []string) (*session.Session, []string) {
	return nil, nil
}
