package db

import (
	//b64 "encoding/base64"

	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	// constクラス

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	cnst "github.com/kemper0530/go-handson-lambda/common"

	// エンティティ(データベースのテーブルの行に対応)
	entity "github.com/kemper0530/go-handson-lambda/models/entity"
)

// Actuator
func ActuaterHealth() entity.ActuatorRslt {
	testmember := []entity.Testmember{}
	actuatorRslt := entity.ActuatorRslt{}

	db := open()
	db.Order("id asc").Find(&testmember)
	actuatorRslt.Db_Status = "UP"
	close(db)

	actuatorRslt.App_Status = "UP"
	// 日本時間へ変換
	jst, _ := time.LoadLocation("Asia/Tokyo")
	actuatorRslt.Time = time.Now().In(jst)
	actuatorRslt.Host, _ = os.Hostname()

	log.Println(actuatorRslt)
	return actuatorRslt
}

// FindAllMembersはメンバー全件取得する
func FetchAllMembers() []entity.Testmember {
	testmember := []entity.Testmember{}

	db := open()
	db.Order("id asc").Find(&testmember)
	close(db)
	return testmember
}

// ログイン情報を取得する
func FindLoginID(username string, password string) entity.Rslt {
	login_info := []entity.Login_info{}
	Rslt := entity.Rslt{}

	db := open()

	// select
	db.First(&login_info, "username=?", username)

	if len(login_info) == cnst.ONE {
		// verify
		errLogin := verify(login_info[0].Password, password)

		if errLogin == nil {
			fmt.Println("Login success!")
			// ログイン成功
			Rslt.Responce = cnst.JsonStatusOK
			Rslt.Result = cnst.ONE
			Rslt.Name = login_info[0].Name
			Rslt.Id = login_info[0].Id
		} else {
			fmt.Println("Login error: ", errLogin)
			// ログイン失敗
			Rslt.Responce = cnst.JsonStatusOK
			Rslt.Result = cnst.ZERO
		}
	} else {
		fmt.Println("Login error no data: ")
		// ログイン失敗
		Rslt.Responce = cnst.JsonStatusOK
		Rslt.Result = cnst.ZERO
	}

	close(db)

	return Rslt
}

// verify
func verify(hash, s string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
}

// work全件取得する
func FetchAllWorker() []entity.Work {
	work := []entity.Work{}

	db := open()
	db.Order("work_id asc").Find(&work)
	close(db)

	return work
}

// カード情報を登録する
func AddCardInfo(cardnumber string, cardname string, cardmonth int, cardyear int, cardcvv string) entity.Rslt {
	crdcardinfo := []entity.Crdcardinfo{}
	Rslt := entity.Rslt{}

	// ハッシュ値の生成　セキュリティコードはbcryptで暗号化して登録
	hashCardcvv, err := bcrypt.GenerateFromPassword([]byte(cardcvv), bcrypt.DefaultCost)
	if err != nil {
		log.Panic("Error bcrypt.GenerateFromPassword!")
		Rslt.Responce = cnst.JsonStatusNG
		Rslt.Result = cnst.ZERO
		return Rslt
	}

	db := open()

	// select
	db.First(&crdcardinfo, "cardnumber=?", cardnumber)

	if len(crdcardinfo) == cnst.ONE {
		// 登録失敗
		Rslt.Responce = cnst.JsonStatusOK
		Rslt.Result = cnst.TWO
	} else {
		var crdcardinfoIns = entity.Crdcardinfo{
			Cardnumber: cardnumber,
			Cardname:   cardname,
			Cardmonth:  cardmonth,
			Cardyear:   cardyear,
			Cardcvv:    string(hashCardcvv),
		}
		// insert
		db.Create(&crdcardinfoIns)
		Rslt.Responce = cnst.JsonStatusOK
		Rslt.Result = cnst.ONE
	}

	close(db)

	return Rslt
}

// メール送信結果テーブルを設定する
func SetMailSendRslt() entity.Mail_send_rslt {

	mail_send_rslt := []entity.Mail_send_rslt{}

	db := open()

	// 送信連番の取得
	count := cnst.ZERO
	sendno := cnst.ZERO

	db.Find(&mail_send_rslt).Count(&count)
	sendno = count + 1

	// テナントIDの定義
	tnntid := cnst.TNNTID

	// msg_idの生成
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		// 登録失敗
		return mail_send_rslt[0]
	}
	msgid := u.String()

	// 日本時間へ変換
	jst, _ := time.LoadLocation("Asia/Tokyo")
	_time := time.Now().In(jst)

	// insert メール送信結果情報(顧客用)
	mail_send_rsltIns := entity.Mail_send_rslt{
		Send_no:         sendno,
		Msg_id:          msgid,
		Tnnt_id:         tnntid,
		Target_sys_type: strconv.Itoa(cnst.ONE),
		Status:          strconv.Itoa(cnst.ZERO),
		Server_id:       cnst.SERVID,
		Priority:        cnst.ONE,
		Send_reg_at:     _time,
		Queue_remove:    strconv.Itoa(cnst.ZERO),
		Updated_at:      _time,
	}

	close(db)

	return mail_send_rsltIns

}

// メール送信情報テーブルを設定
func SetMailSendInf2C(to_email string, name string, text string, from_email string, personal_name string, msgid string, id int) entity.Mail_send_inf {

	mst_ssmlknr := []entity.Mst_ssmlknr{}
	tmpuserinfo := []entity.Tmpuserinfo{}

	db := open()

	// 送信管理マスタの取得
	db.Where("id = ?", id).First(&mst_ssmlknr)
	subject := mst_ssmlknr[0].Subject
	body := mst_ssmlknr[0].Body

	// 仮登録の場合
	if id == cnst.TWO {
		// tokenの取得
		db.Where("email = ?", to_email).First(&tmpuserinfo)
		token := tmpuserinfo[0].Token
		// URLの生成
		path := os.Getenv("SIGN_UP_PATH")
		query := path + "?token=" + token

		// 文字列の置き換え　$1　→　登録名、$2 -> URL
		body = strings.Replace(body, "$1", name, -1)
		body = strings.Replace(body, "$2", query, -1)
	} else {
		body = strings.Replace(body, "$1", name, -1)
	}

	// insert メール送信情報(顧客用)
	mail_send_infIns := entity.Mail_send_inf{
		Msg_id:        msgid,
		From_email:    from_email,
		To_email:      to_email,
		Subject:       subject,
		Body:          body,
		Personal_name: personal_name,
	}

	close(db)

	return mail_send_infIns
}

// メール送信情報テーブルを設定
func SetMailSendInf2Y(to_email string, name string, text string, from_email string, personal_name string, msgid string, id int) entity.Mail_send_inf {

	mst_ssmlknr := []entity.Mst_ssmlknr{}

	db := open()

	// 送信管理マスタの取得
	db.Where("id = ?", id).First(&mst_ssmlknr)
	replytitle := mst_ssmlknr[0].Replytitle
	toreply := mst_ssmlknr[0].Toreply

	// insert メール送信情報(送信者用)
	mail_send_infIns := entity.Mail_send_inf{
		Msg_id:        msgid,
		From_email:    from_email,
		To_email:      toreply,
		Subject:       replytitle,
		Body:          text,
		Personal_name: personal_name,
	}

	close(db)

	return mail_send_infIns
}

// お問合せ内容をメール送信情報テーブル、結果テーブルへ登録する
func SetMailRegist(sendInf *entity.Mail_send_inf, sendRslt *entity.Mail_send_rslt) entity.Rslt {

	Rslt := entity.Rslt{}

	db := open()

	// insert
	db.Create(&sendInf)
	db.Create(&sendRslt)

	Rslt.Responce = cnst.JsonStatusOK
	Rslt.Result = cnst.ONE

	close(db)

	return Rslt
}

// お問合せ内容をメール送信情報テーブル、結果テーブルへ登録する
func SendMailRegist(to_email string, name string, text string, from_email string, personal_name string) entity.Rslt {

	mail_send_rslt := []entity.Mail_send_rslt{}
	mst_ssmlknr := []entity.Mst_ssmlknr{}
	Rslt := entity.Rslt{}

	db := open()

	// 送信連番の取得
	count := cnst.ZERO
	sendno1 := cnst.ZERO
	db.Find(&mail_send_rslt).Count(&count)
	sendno1 = count + 1
	sendno2 := sendno1 + 1

	// テナントIDの定義
	tnntid := cnst.TNNTID

	// msg_id1の生成
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		// 登録失敗
		Rslt.Responce = cnst.JsonStatusOK
		Rslt.Result = cnst.TWO
		return Rslt
	}
	msgid1 := u.String()

	// msg_id2の生成
	u2, err2 := uuid.NewRandom()
	if err2 != nil {
		fmt.Println(err2)
		// 登録失敗
		Rslt.Responce = cnst.JsonStatusOK
		Rslt.Result = cnst.TWO
		return Rslt
	}
	msgid2 := u2.String()

	// 送信管理マスタの取得
	db.First(&mst_ssmlknr)
	subject := mst_ssmlknr[0].Subject
	body := mst_ssmlknr[0].Body
	replytitle := mst_ssmlknr[0].Replytitle
	toreply := mst_ssmlknr[0].Toreply

	// 文字列の置き換え　$1　→　登録名
	_body := strings.Replace(body, "$1", name, -1)

	// 日本時間へ変換
	jst, _ := time.LoadLocation("Asia/Tokyo")
	_time := time.Now().In(jst)

	// insert メール送信情報(顧客用)
	var mail_send_infIns = entity.Mail_send_inf{
		Msg_id:        msgid1,
		From_email:    from_email,
		To_email:      to_email,
		Subject:       subject,
		Body:          _body,
		Personal_name: personal_name,
	}

	// insert メール送信結果情報(顧客用)
	var mail_send_rsltIns = entity.Mail_send_rslt{
		Send_no:         sendno1,
		Msg_id:          msgid1,
		Tnnt_id:         tnntid,
		Target_sys_type: strconv.Itoa(cnst.ONE),
		Status:          strconv.Itoa(cnst.ZERO),
		Server_id:       cnst.SERVID,
		Priority:        cnst.ONE,
		Send_reg_at:     _time,
		Queue_remove:    strconv.Itoa(cnst.ZERO),
		Updated_at:      _time,
	}

	// insert
	db.Create(&mail_send_infIns)
	db.Create(&mail_send_rsltIns)

	// insert メール送信情報(送信者用)
	mail_send_infIns = entity.Mail_send_inf{
		Msg_id:        msgid2,
		From_email:    from_email,
		To_email:      toreply,
		Subject:       replytitle,
		Body:          text,
		Personal_name: personal_name,
	}

	// insert メール送信結果情報(送信者用)
	mail_send_rsltIns = entity.Mail_send_rslt{
		Send_no:         sendno2,
		Msg_id:          msgid2,
		Tnnt_id:         tnntid,
		Target_sys_type: strconv.Itoa(cnst.ONE),
		Status:          strconv.Itoa(cnst.ZERO),
		Server_id:       cnst.SERVID,
		Priority:        cnst.ONE,
		Send_reg_at:     _time,
		Queue_remove:    strconv.Itoa(cnst.ZERO),
		Updated_at:      _time,
	}

	// insert
	db.Create(&mail_send_infIns)
	db.Create(&mail_send_rsltIns)

	Rslt.Responce = cnst.JsonStatusOK
	Rslt.Result = cnst.ONE

	close(db)

	return Rslt
}

// profile全件取得する
func FetchProfileInfo() []entity.Profile {
	profile := []entity.Profile{}

	db := open()
	db.Order("id asc").Find(&profile)
	close(db)

	return profile
}

// 仮アカウント情報を登録する
func RegistLoginID(email string, password string, name string) entity.Rslt {
	login_info := []entity.Login_info{}
	Rslt := entity.Rslt{}

	db := open()

	// 登録情報の確認
	db.First(&login_info, "username=?", email)

	if len(login_info) == cnst.ONE {
		Rslt.Responce = cnst.JsonStatusOK
		Rslt.Result = cnst.ZERO
		return Rslt
	}

	// ハッシュ値の生成　パスワードはbcryptで暗号化して登録
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic("Error bcrypt.GenerateFromPassword!")
		Rslt.Responce = cnst.JsonStatusNG
		Rslt.Result = cnst.ZERO
		return Rslt
	}

	// Tokenの生成
	token := RandString6(36)

	// 期限日の設定
	expired := time.Now()
	expired = expired.Add(time.Duration(24) * time.Hour)

	// insert Tmpuser情報
	var tmpuserinfoIns = entity.Tmpuserinfo{
		Email:      email,
		Password:   string(hashPassword),
		Name:       name,
		Token:      token,
		Expired:    expired,
		Updated_at: time.Now(),
	}

	// insert
	db.Create(&tmpuserinfoIns)

	Rslt.Responce = cnst.JsonStatusOK
	Rslt.Result = cnst.ONE

	close(db)

	return Rslt
}

// work全件取得する
func FetchMailAdrInfo(id int) entity.Rslt {
	login_info := []entity.Login_info{}
	Rslt := entity.Rslt{}

	db := open()
	db.First(&login_info, "id=?", id)

	Rslt.Id = login_info[0].Id
	Rslt.Email = login_info[0].Username
	Rslt.Name = login_info[0].Name
	Rslt.Responce = cnst.JsonStatusOK
	Rslt.Result = cnst.ONE

	close(db)

	return Rslt
}

func RandString6(n int) string {
	var randSrc = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	cache, remain := randSrc.Int63(), cnst.Rs6LetterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), cnst.Rs6LetterIdxMax
		}
		idx := int(cache & cnst.Rs6LetterIdxMask)
		if idx < len(cnst.Rs6Letters) {
			b[i] = cnst.Rs6Letters[idx]
			i--
		}
		cache >>= cnst.Rs6LetterIdxBits
		remain--
	}
	return string(b)
}

// 押下したURLを検証する
func FetchSignUpAccountMail(token string) int {
	tmpuserinfo := []entity.Tmpuserinfo{}

	db := open()
	db.First(&tmpuserinfo, "token=?", token)

	// 値の取得ができなかった場合
	if len(tmpuserinfo) == cnst.ZERO {
		return cnst.ZERO
	}

	tn := time.Now()

	// 期限日を超過している場合
	if !tmpuserinfo[0].Expired.After(tn) {
		return cnst.ZERO
	}

	email := tmpuserinfo[0].Email
	password := tmpuserinfo[0].Password
	name := tmpuserinfo[0].Name

	// insert ログイン情報
	var login_infoIns = entity.Login_info{
		Username:   email,
		Password:   password,
		Name:       name,
		Updated_at: tn,
	}

	// ログイン情報にInsert
	db.Create(&login_infoIns)

	// 仮ログインテーブルから削除
	db.Where("email = ?", email).Delete(entity.Tmpuserinfo{})

	close(db)

	return cnst.ONE
}

// アクセスログを登録する
func RegistAccessLog(c *gin.Context) entity.Rslt {
	Rslt := entity.Rslt{}

	// 日本時間へ変換
	jst, _ := time.LoadLocation("Asia/Tokyo")
	_time := time.Now().In(jst)

	db := open()

	// insert アクセスログ
	var access_logsIns = entity.Access_logs{
		User_id:        c.PostForm("user_id"),
		Event_id:       c.PostForm("event_id"),
		Access_ip:      c.PostForm("access_ip"),
		City:           c.PostForm("city"),
		Region:         c.PostForm("region"),
		Region_code:    c.PostForm("region_code"),
		Country_name:   c.PostForm("country_name"),
		Country_code:   c.PostForm("country_code"),
		Continent_name: c.PostForm("continent_name"),
		Continent_code: c.PostForm("continent_code"),
		Latitude:       c.PostForm("latitude"),
		Longitude:      c.PostForm("longitude"),
		Postal:         c.PostForm("postal"),
		Calling_code:   c.PostForm("calling_code"),
		Created_at:     _time,
	}

	// アクセスログにInsert
	db.Create(&access_logsIns)

	close(db)

	Rslt.Responce = cnst.JsonStatusOK
	Rslt.Result = cnst.ONE

	return Rslt
}

// Lambdaからリクエストされた内容を登録する
func RegistBounce(c *gin.Context) entity.Rslt {

	mail_send_rslt := []entity.Mail_send_rslt{}
	Rslt := entity.Rslt{}

	messageId := c.PostForm("messageId")

	db := open()
	db.First(&mail_send_rslt, "msg_id_ses=?", messageId)

	// 値の取得ができなかった場合
	if len(mail_send_rslt) == cnst.ZERO {
		Rslt.Responce = cnst.JsonStatusNG
		Rslt.Result = cnst.ONE
		return Rslt
	}

	if( c.PostForm("notificationType") == cnst.BOUNCE ) {

		// 日本時間へ変換
		jst, _ := time.LoadLocation("Asia/Tokyo")
		_time := time.Now().In(jst)

		// @の出現位置取得
		i_atmrk := strings.Index(c.PostForm("source"), "@")

		// UNIXタイムスタンプの取得
		i_timestamp := int(time.Now().Unix())

		// bounce_mail_detail へ insert
		var bounce_mail_detail = entity.Bounce_mail_detail{
			Send_no:         mail_send_rslt[0].Send_no,
			Msg_id_ses:      c.PostForm("messageId"),
			Addresser:       c.PostForm("source"),
			Timestamp:       i_timestamp,
			Smtp_command:    string([]rune(c.PostForm("diagnosticCode"))[:10]),
			Recipient:       c.PostForm("recipients"),
			Destination:     c.PostForm("destination"),
			Reply_code:      c.PostForm("notificationType"),
			Rhost:           c.PostForm("sourceIp"),
			Diagnostic_type: c.PostForm("bounceType"),
			Action:          c.PostForm("action"),
			Diagnostic_code: c.PostForm("status"),
			Reason:          c.PostForm("bounceType"),
			Delivery_status: c.PostForm("notificationType"),
			Sender_domain:   c.PostForm("source")[i_atmrk+1:],
		}

		// bounce_mail_detailにInsertする
		db.Create(&bounce_mail_detail)

		// mail_send_rsltをUpdateする
		db.Model(&mail_send_rslt).Where("msg_id_ses = ?", c.PostForm("messageId")).Update("status", "9")

		// blacklistにInsertする
		var blklst_send_mailadr = entity.Blklst_send_mailadr{
			Email:       c.PostForm("recipients"),
			Updated_at:  _time,
		}

		db.Create(&blklst_send_mailadr)

	}

	Rslt.Responce = cnst.JsonStatusOK
	Rslt.Result = cnst.ONE
	return Rslt
}
