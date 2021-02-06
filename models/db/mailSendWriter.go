package db

import (
	"errors"
	"log"
	"reflect"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/jinzhu/gorm"

	// constクラス
	cnst "github.com/kemper0530/go-handson-lambda/common"

	// エンティティ(データベースのテーブルの行に対応)
	entity "github.com/kemper0530/go-handson-lambda/models/entity"
)

func FetchMailSendSelect() entity.Rslt {
	mail_send_rslt := []entity.Mail_send_rslt{}
	Rslt := entity.Rslt{}
	var count int

	db := open()
	db.Where("mail_send_rslt.status = ?", cnst.ZERO).Find(&mail_send_rslt).Count(&count)
	if reflect.DeepEqual(count, cnst.ZERO) {
		Rslt.Responce = cnst.JsonStatusOK
		Rslt.Result = cnst.ONE
		close(db)
		return Rslt
	}

	rslt := SendEmailPrepare(db)
	if rslt {
		Rslt.Responce = cnst.JsonStatusOK
		Rslt.Result = cnst.ONE
		log.Println("mail sending success")
	} else {
		Rslt.Responce = cnst.JsonStatusNG
		Rslt.Result = cnst.ZERO
		log.Println("mail sending error!!")
	}

	close(db)

	return Rslt
}

func SendEmailPrepare(db *gorm.DB) bool {
	mail_send_inf := []entity.Mail_send_inf{}

	db.Where("rslt.status = ?", cnst.ZERO).Joins("inner join mail_send_rslt rslt on rslt.msg_id = mail_send_inf.msg_id").Order("rslt.send_no asc").Limit(1).Find(&mail_send_inf)

	if reflect.DeepEqual(len(mail_send_inf), cnst.ZERO) {
		return false
	}
	msgid := mail_send_inf[0].Msg_id
	from := mail_send_inf[0].From_email
	to := mail_send_inf[0].To_email
	subject := mail_send_inf[0].Subject
	body := mail_send_inf[0].Body

	// メールを送信する
	messageIdSES, err := sendSESEmail(from, to, subject, body)
	if err != nil {
		UpdateMailSendRslt(msgid, db, cnst.FOUR, messageIdSES)
		return false
	}

	// メール送信結果情報のステータスを送信済に変更する
	err = UpdateMailSendRslt(msgid, db, cnst.TWO, messageIdSES)
	if err != nil {
		log.Println(err)
		return false
	}

	// メール送信情報から対象レコードを削除する
	err = DeleteMailSendInf(msgid, db)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func sendSESEmail(from string, to string, title string, body string) (*string, error) {
	// 環境変数ファイルの読込
	// AWS_REGION := awsEnvload("AWS_REGION")
	// AWS_ACCESS_KEY_ID := awsEnvload("AWS_ACCESS_KEY_ID")
	// AWS_SECRET_KEY := awsEnvload("AWS_SECRET_KEY")
	AWS_REGION := os.Getenv("AWS_SES_REGION")
	AWS_ACCESS_KEY_ID := os.Getenv("AWS_SES_ACCESS_KEY_ID")
	AWS_SECRET_KEY := os.Getenv("AWS_SES_SECRET_KEY")

	//AWS-SESの設定情報を格納する
	awsSession := session.New(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_KEY, ""),
	})

	// メール送信情報を設定する
	svc := ses.New(awsSession)
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(title),
			},
		},
		Source: aws.String(from),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return result.MessageId, errors.New(err.Error())
	}

	log.Println("Email Sent to address: " + to)
	log.Println(result)

	return result.MessageId, nil
}

func UpdateMailSendRslt(msgid string, db *gorm.DB, status int, msgidSES *string) error {
	mail_send_rslt := []entity.Mail_send_rslt{}

	// 送信済に変更する
	rslt := db.Model(&mail_send_rslt).Where("msg_id = ?", msgid).Updates(map[string]interface{}{"status": status, "queue_id": cnst.QUEUEID, "msg_id_ses": msgidSES})
	if rslt == nil {
		return errors.New("ステータスの更新に失敗しました")
	}

	return nil
}

func DeleteMailSendInf(msgid string, db *gorm.DB) error {
	mail_send_inf := []entity.Mail_send_inf{}

	// 送信済のレコードを削除する
	rslt := db.Where("msg_id = ?", msgid).Delete(&mail_send_inf)
	if rslt == nil {
		return errors.New("メール送信情報テーブルの削除に失敗しました")
	}

	return nil
}
