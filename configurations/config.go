package configurations

type Configuration struct {
	ConnectionString string
	TokenIssuer      string
	TokenAudience    string
	TokenSigningKey  string
}

var appConfig *Configuration

//to be set at startup
func SetConfiguration(config *Configuration) {
	appConfig = config
}

//returns TokenIssuer and TokenAudience
func GetTokenIssuerAndAudience() (string, string) {
	return appConfig.TokenIssuer, appConfig.TokenAudience
}

//returns the token signing key (base64 encoded)
func GetTokenSigningKey() string {
	return appConfig.TokenSigningKey
}

//returns the db connection string
func GetConnectionString() string {
	return appConfig.ConnectionString
}
