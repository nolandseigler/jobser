package auth

type Config struct {
	PubKeyPath     string `json:"protocol" form:"protocol" query:"protocol" mapstructure:"POSTGRES_PROTOCOL"`
	PrivKeyPath     string `json:"username" form:"username" query:"username" mapstructure:"POSTGRES_USERNAME"`
}

