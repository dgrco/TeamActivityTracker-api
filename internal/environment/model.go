package environment

type Environment struct {
	Env				 string		// "dev" | "prod"
	DatabaseURL      string
	JWTSecret        string
	Port             string
	CookieSecureMode bool
	WebAppServerURL	 string
}
