package config

// import (
// 	"fmt"

// 	"github.com/resend/resend-go/v2"
// )

// func SendHtmlMail(html string, recipient []string) error{
// 	resend_key := GetResendSecret()
//     client := resend.NewClient(resend_key)
// 	params := &resend.SendEmailRequest{
//         From:    "fme <onboarding@resend.dev>",
//         To:      recipient,
//         Html:    html, // Pass the rendered HTML content here
//         Subject: "OTP Retrieval",
//         Cc:      []string{"cc@example.com"},
//         Bcc:     []string{"bcc@example.com"},
//         ReplyTo: "replyto@example.com",
//     }
// 	sent, err := client.Emails.Send(params)
//     if err != nil {
//         return err
//     }
//     fmt.Println(sent.Id)
//     return nil

// }
