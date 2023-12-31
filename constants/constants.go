package constants

var (
	// 5 MB
	FiberBodyLimit = 5 * 1024 * 1024

	// Key Redis
	// token - google id - phone number - identity
	KeyVerificationSession = "verification#%s#%s#%s#%s"
	KeyOTP                 = "otp#%s"
	KeyOTPIncrement        = "otp#%s#increment"
	KeyUserPoint           = "user_point#%s"
	KeyBindGoogle          = "bind_google#%s#%s"

	// Key Locals
	KeyLocalsUser = "user"

	// Regex
	RegexPhoneNumber = `^(\+62|62)?[\s-]?(0)?8[1-9]{1}\d{1}[\s-]?\d{4}[\s-]?\d{2,5}$`

	// Anything
	WhatsappNumberSuffix      = "@s.whatsapp.net"
	OTPLength                 = 6
	DefaultPageLimit          = 10
	DefaultPageLimitStatistic = 6
)
