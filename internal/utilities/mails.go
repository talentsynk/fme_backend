package utilities

import (
	"fme_backend/internal/config"
	"fmt"
)

func SendOtpEmail(email, otp string) {
	emailString := "Your OTP to reset your password is " + otp
	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Reset Your Password</title>
		<style>
			.container {
				font-family: Arial, sans-serif;
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
				border: 1px solid #ddd;
				border-radius: 5px;
				background-color: #f9f9f9;
			}
			.header {
				text-align: center;
				margin-bottom: 20px;
			}
			.otp {
				font-size: 24px;
				font-weight: bold;
				color: #333;
				text-align: center;
				margin: 20px 0;
			}
			.footer {
				text-align: center;
				margin-top: 30px;
				font-size: 12px;
				color: #777;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h2>Reset Your Password</h2>
			</div>
			<p>Hello,</p>
			<p>You requested to reset your password. Use the OTP below to reset it:</p>
			<div class="otp">%s</div>
			<p>If you didn't request this, please ignore this email.</p>
			<p>Thank you,<br>NASIC Team</p>
			<div class="footer">
				<p>© 2024 NASIC. All rights reserved.</p>
			</div>
		</div>
	</body>
	</html>`, otp)

	subject := "Reset your password"
	_ = config.SendMail(emailString, htmlBody, subject, email) // TODO: log this error 
}
