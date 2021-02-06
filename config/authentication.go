package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"encoding/base64"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	cnst "github.com/kemper0530/go-handson-lambda/common"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// FireBaseの設定ファイル読込
func SetUpFirebase() (auth *auth.Client, err error) {
	// 変数初期化
	var opt option.ClientOption

	// FirebaseのSDKを使用するためのkeyを読み込み
	if os.Getenv("GO_ENV") == "production" {
		sEnc := os.Getenv("FIREBASE_KEYFILE_JSON")
		sDec, _ := base64.StdEncoding.DecodeString(sEnc)
		opt = option.WithCredentialsJSON(sDec)
	}else{
		opt = option.WithCredentialsFile(GetFireBasePath())
	}

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	auth, errAuth := app.Auth(context.Background())
	return auth, errAuth
}

// JWT検証
func AuthFirebase(c *gin.Context, auth *auth.Client) (result int, errMsg string) {

	// クライアントから送られてきた JWT 取得
	authHeader := c.GetHeader("Authorization")
	log.Printf(authHeader)
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)
	log.Printf(idToken)

	// JWT の検証
	u := ""
	if os.Getenv("MODE") == "TEST" {
		if os.Getenv("TEST_JWT") == idToken {
			u = fmt.Sprintf("[TEST] Success JWT Verify :%v", idToken)
			log.Printf("[TEST] success JWT Verify :%v", idToken)
		} else {
			u = fmt.Sprintf("[TEST] error JWT Verify :%v", idToken)
			log.Printf("[TEST] error JWT Verify :%v", idToken)
		}
	} else {
		log.Println(auth)
		// ID トークンの取り消しを検出
		_, err := auth.VerifyIDTokenAndCheckRevoked(context.Background(), idToken)
		if err != nil {
        if err.Error() == "ID token has been revoked" {
					// Token is revoked. Inform the user to reauthenticate or signOut() the user.
					log.Printf("Token is revoked.")
					return cnst.JsonStatusNG, u
        } else {
					// Token is invalid
					log.Printf("Token is invalid")
					return cnst.JsonStatusNG, u
        }
		}
		_, err = auth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			u = fmt.Sprintf("error verifying ID token: %v\n", err)
			log.Printf("error verifying ID token: %v\n", err)
			return cnst.JsonStatusNG, u
		}
	}
	return cnst.JsonStatusOK, u
}

// firebase json path
func GetFireBasePath() string {
	// 環境変数の読込
	firebasejsonpath := os.Getenv("FIREBASE_PATH")
	return firebasejsonpath
}
