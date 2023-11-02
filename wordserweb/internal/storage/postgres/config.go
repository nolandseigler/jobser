package postgres

type Config struct {
	Protocol     string `json:"protocol" form:"protocol" query:"protocol" mapstructure:"POSTGRES_PROTOCOL"`
	Username     string `json:"username" form:"username" query:"username" mapstructure:"POSTGRES_USERNAME"`
	Password     string `json:"password" form:"password" query:"password" mapstructure:"POSTGRES_PASSWORD"`
	Hostname     string `json:"hostname" form:"hostname" query:"hostname" mapstructure:"POSTGRES_HOSTNAME"`
	Port         int    `json:"port" form:"port" query:"port" mapstructure:"POSTGRES_PORT"`
	DatabaseName string `json:"databaseName" form:"databaseName" query:"databaseName" mapstructure:"POSTGRES_DATABASE_NAME"`
}
