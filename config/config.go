type Config struct{
	DBHost		string	`env:MYSQL_HOST" envDefault:"127.0.0.1"`
	DBUser		int		`env:MYSQL_USER" envDefault:"mydb"`
	DBPassword	string	`env:MYSQL_PASSWORD" envDefault:"mydb"`
	DBName		string	`env:MYSQL_DATABASE" envDefault:"mydb"`
}