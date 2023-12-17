package constants

const (
	TemplateAuthOtp = "Your {{ .Name }} verification code is: *{{ .Code }}*.\nDo not share this code with anyone."
	TemplateHistory = "Hi, {{ .Name }}.\nThis is a information history of your transaction.\n\nDevice      : *{{ .Device }}*\nSuccess   : *{{ .Success }}*\nFailed       : *{{ .Failed }}*\n\nThank you for using our service.\nHave a nice day."
)
