package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kemper0530/go-handson-lambda/models/entity"

	// authクラス
	authcnfg "github.com/kemper0530/go-handson-lambda/config"
	// constクラス
	cnst "github.com/kemper0530/go-handson-lambda/common"
	// DBアクセス用モジュール
	db "github.com/kemper0530/go-handson-lambda/models/db"
	// httpアクセス用モジュール
	rest "github.com/kemper0530/go-handson-lambda/models/rest"
)

// メールバッチ処理
func FetchMailSendSelect() {
	resultProduct := db.FetchMailSendSelect()
	if resultProduct.Result == cnst.ONE {
		fmt.Println("Run AmazonMail SES!")
	} else {
		fmt.Println("Not Run AmazonMail SES!")
	}
}

// Goアプリのステータスを返却する
func ActuaterHealth(c *gin.Context) {
	resultProducts := db.ActuaterHealth()
	c.JSON(http.StatusOK, resultProducts)
}

// FetchAllMembers は メンバー情報を取得する
func FetchAllMembers(c *gin.Context) {
	resultStatus, errMsg := authcnfg.AuthFirebase(c, cnst.Auth)
	if resultStatus == cnst.JsonStatusNG {
		c.JSON(http.StatusBadRequest, errMsg)
	} else {
		resultProducts := db.FetchAllMembers()
		// URLへのアクセスに対してJSONを返す
		c.JSON(http.StatusOK, resultProducts)
	}
}

// work情報を取得する
func FetchAllWorker(c *gin.Context) {
	resultStatus, errMsg := authcnfg.AuthFirebase(c, cnst.Auth)
	if resultStatus == cnst.JsonStatusNG {
		c.JSON(http.StatusBadRequest, errMsg)
	} else {
		resultProducts := db.FetchAllWorker()
		// URLへのアクセスに対してJSONを返す
		c.JSON(http.StatusOK, resultProducts)
	}
}

// FetchLoginInfo は 指定したIDのパスワードを取得する
func FetchLoginInfo(c *gin.Context) {
	resultStatus, errMsg := authcnfg.AuthFirebase(c, cnst.Auth)
	if resultStatus == cnst.JsonStatusNG {
		c.JSON(http.StatusBadRequest, errMsg)
	} else {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if len(username) == cnst.ZERO || len(password) == cnst.ZERO {
			log.Printf(username)
			log.Printf(password)
			log.Panic("Error nothing URL parameter!!")
		}

		resultProduct := db.FindLoginID(username, password)

		// URLへのアクセスに対してJSONを返す
		c.JSON(http.StatusOK, resultProduct)
	}
}

// クレジットカード情報を登録する
func FetchCreditInfoRegist(c *gin.Context) {
	resultStatus, errMsg := authcnfg.AuthFirebase(c, cnst.Auth)
	if resultStatus == cnst.JsonStatusNG {
		c.JSON(http.StatusBadRequest, errMsg)
	} else {
		cardnumber := c.PostForm("cardnumber")
		cardname := c.PostForm("cardname")
		cardmonth := c.PostForm("cardmonth")
		cardyear := c.PostForm("cardyear")
		cardcvv := c.PostForm("cardcvv")

		if len(cardnumber) == cnst.ZERO {
			log.Panic("Error nothing URL parameter!!")
		}

		cardmonthInt, _ := strconv.Atoi(cardmonth)
		cardyearInt, _ := strconv.Atoi(cardyear)

		resultProduct := db.AddCardInfo(cardnumber, cardname, cardmonthInt, cardyearInt, cardcvv)

		// URLへのアクセスに対してJSONを返す
		c.JSON(http.StatusOK, resultProduct)
	}
}

// お問合せ内容を登録する
func FetchSendMailRegist(c *gin.Context) {
	resultStatus, errMsg := authcnfg.AuthFirebase(c, cnst.Auth)
	if resultStatus == cnst.JsonStatusNG {
		c.JSON(http.StatusBadRequest, errMsg)
	} else {
		to_email := c.PostForm("to_email")
		name := c.PostForm("name")
		text := c.PostForm("text")
		from_email := c.PostForm("from_email")
		personal_name := c.PostForm("personal_name")

		if len(to_email) == cnst.ZERO &&
			len(name) == cnst.ZERO &&
			len(text) == cnst.ZERO {
			log.Panic("Error nothing URL parameter!!")
		}

		// 顧客向けのメール情報を設定する
		Mail_send_rslt := db.SetMailSendRslt()
		Mail_send_inf := db.SetMailSendInf2C(to_email, name, text, from_email, personal_name, Mail_send_rslt.Msg_id, cnst.ONE)
		resultProduct := db.SetMailRegist(&Mail_send_inf, &Mail_send_rslt)

		// メールバッチ処理を直接コールする
		FetchMailSendSelect()

		// 管理者向けのメール情報を設定する
		Mail_send_rslt = db.SetMailSendRslt()
		Mail_send_inf = db.SetMailSendInf2Y(to_email, name, text, from_email, personal_name, Mail_send_rslt.Msg_id, cnst.ONE)
		resultProduct = db.SetMailRegist(&Mail_send_inf, &Mail_send_rslt)

		// URLへのアクセスに対してJSONを返す
		c.JSON(http.StatusOK, resultProduct)

		// メールバッチ処理を直接コールする
		FetchMailSendSelect()
	}
}

// profile情報を取得する
func FetchProfileInfo(c *gin.Context) {
	resultStatus, errMsg := authcnfg.AuthFirebase(c, cnst.Auth)
	if resultStatus == cnst.JsonStatusNG {
		c.JSON(http.StatusBadRequest, errMsg)
	} else {
		resultProducts := db.FetchProfileInfo()
		// URLへのアクセスに対してJSONを返す
		c.JSON(http.StatusOK, resultProducts)
	}
}

// FetchRegistAcount は アカウントの登録を実施する
func FetchRegistAccount(c *gin.Context) {
	resultStatus, errMsg := authcnfg.AuthFirebase(c, cnst.Auth)
	if resultStatus == cnst.JsonStatusNG {
		c.JSON(http.StatusBadRequest, errMsg)
	} else {
		username := c.PostForm("email")
		password := c.PostForm("password")
		name := c.PostForm("name")

		if len(username) == cnst.ZERO || len(password) == cnst.ZERO {
			log.Panic("Error nothing URL parameter!!")
		}

		resultProduct := db.RegistLoginID(username, password, name)

		// URLへのアクセスに対してJSONを返す
		c.JSON(http.StatusOK, resultProduct)
	}
}

// FetchRegistAcountMail は 送信先へのメール情報を登録する
func FetchRegistAccountMail(c *gin.Context) {
	resultStatus, errMsg := authcnfg.AuthFirebase(c, cnst.Auth)
	if resultStatus == cnst.JsonStatusNG {
		c.JSON(http.StatusBadRequest, errMsg)
	} else {
		to_email := c.PostForm("to_email")
		name := c.PostForm("name")
		text := c.PostForm("text")
		from_email := c.PostForm("from_email")
		personal_name := c.PostForm("personal_name")

		if len(to_email) == cnst.ZERO {
			log.Panic("Error nothing URL parameter!!")
		}

		// 顧客向けのメール情報を設定する
		Mail_send_rslt := db.SetMailSendRslt()
		Mail_send_inf := db.SetMailSendInf2C(to_email, name, text, from_email, personal_name, Mail_send_rslt.Msg_id, cnst.TWO)
		resultProduct := db.SetMailRegist(&Mail_send_inf, &Mail_send_rslt)

		// メールバッチ処理を直接コールする
		FetchMailSendSelect()

		// URLへのアクセスに対してJSONを返す
		c.JSON(http.StatusOK, resultProduct)
	}
}

// FetchMailAdrInfo は 指定したIDのメールアドレスと氏名を取得する
func FetchMailAdrInfo(c *gin.Context) {
	resultStatus, errMsg := authcnfg.AuthFirebase(c, cnst.Auth)
	if resultStatus == cnst.JsonStatusNG {
		c.JSON(http.StatusBadRequest, errMsg)
	} else {
		id := c.PostForm("id")
		if len(id) == cnst.ZERO {
			log.Panic("Error nothing URL parameter!!")
		}

		id_int, _ := strconv.Atoi(id)
		resultProduct := db.FetchMailAdrInfo(id_int)

		// URLへのアクセスに対してJSONを返す
		c.JSON(http.StatusOK, resultProduct)
	}
}

// FetchSignUpAccountMail は メールリンク認証を実装する。
func FetchSignUpAccountMail(c *gin.Context) {
	token := c.Query("token")

	if len(token) == cnst.ZERO {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{})
	} else {
		authFlg := db.FetchSignUpAccountMail(token)

		if authFlg == cnst.ONE {
			c.HTML(http.StatusOK, "success.tmpl", gin.H{})
		} else {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{})
		}
	}
}

// NEWSAPIの記事を取得し、フロントへ返却する
func FetchNewsInfo(c *gin.Context) {
	resultStatus, errMsg := authcnfg.AuthFirebase(c, cnst.Auth)
	if resultStatus == cnst.JsonStatusNG {
		c.JSON(http.StatusBadRequest, errMsg)
	} else {
		category := c.PostForm("category")
		if len(category) == cnst.ZERO {
			log.Panic("Error nothing URL parameter!!")
		}

		url := os.Getenv("NEWS_URL")
		apikey := os.Getenv("NEWS_APIKEY")

		fmt.Println("【URL】:" + url + "&apikey=" + apikey + "&category=" + category)

		// httpリクエストを実施する
		status, resp := rest.FetchRequest(url+"&apikey="+apikey+"&category="+category, "GET")
		// 200以外の場合は即時返却する
		if status != 200 {
			c.JSON(status, nil)
		}
		// 返却する構造体を定義
		var na entity.NewsAPI

		// Jsonデコード
		if err := json.Unmarshal(resp, &na); err != nil {
			log.Fatal(err)
		}
		// URLへのアクセスに対してJSONを返す
		c.JSON(status, na)
	}
}

// FetchRegistAccessLog は アクセスログの登録を実施する
func FetchRegistAccessLog(c *gin.Context) {
	resultStatus, errMsg := authcnfg.AuthFirebase(c, cnst.Auth)
	if resultStatus == cnst.JsonStatusNG {
		c.JSON(http.StatusBadRequest, errMsg)
	} else {
		// アクセスログをDBへ登録する
		resultProduct := db.RegistAccessLog(c)
		// URLへのアクセスに対してJSONを返す
		c.JSON(http.StatusOK, resultProduct)
	}
}

// Lambdaからリクエストされた内容を登録する
func FetchRegistBounce(c *gin.Context) {

	messageId := c.PostForm("messageId")

	if len(messageId) == cnst.ZERO {
		// URLへのアクセスに対してJSONを返す
		c.JSON(http.StatusOK, gin.H{
						"Responce": cnst.JsonStatusOK,
						"Result": cnst.ONE,
        })
	}else{
		resultProduct := db.RegistBounce(c)
		// URLへのアクセスに対してJSONを返す
		c.JSON(http.StatusOK, resultProduct)
	}
}
