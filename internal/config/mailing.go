package config

import (
	"context"
	"fmt"

	"github.com/mrz1836/postmark"
)

func SendMail(emailString,HTMLBody,subject, recipient string) error{
	if GetEnvType() == "testing" {
		fmt.Println("testing mail")
		return nil
	} else if (GetEnvType() == "testing-mail") {
		client := postmark.NewClient(GetSendMailServerToken(),GetSendMailAcctToken())

	email := postmark.Email{
		From:       GetHomeMail(),
		To:         recipient,
		Subject:    subject,
		HTMLBody:   HTMLBody,
		TextBody:   emailString,
		Tag:        "pw-reset",
		TrackOpens: true,
	}

	_, err := client.SendEmail(context.Background(), email)
	if err != nil {
		return err
	}
	return nil
	}

	//add actual server mail
	return nil
	
}
