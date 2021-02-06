package common

import "firebase.google.com/go/auth"

const (
	// 返却時の正常レスポンス
	HttpStatusOK int = 200

	// JSON返却時の正常レスポンス
	JsonStatusOK int = 200

	// JSON返却時のNGレスポンス
	JsonStatusNG int = 500

	// テナントID
	TNNTID string = "9999"

	// サーバID
	SERVID string = "kemper0530.com"

	// ZERO
	ZERO int = 0
	// ONE
	ONE int = 1
	// TWO
	TWO int = 2
	// THREE
	THREE int = 3
	// FOUR
	FOUR int = 4
	// QueueID
	QUEUEID string = "Amazon SES"
	// Token用の文字列
	Rs6Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	Rs6LetterIdxBits = 6
	Rs6LetterIdxMask = 1<<Rs6LetterIdxBits - 1
	Rs6LetterIdxMax  = 63 / Rs6LetterIdxBits

	BOUNCE  = "Bounce"
)

var (
	Auth *auth.Client
)
