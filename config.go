package aurora

// Config server.
type Config struct {
	// Port for address server, if it is empty by default it will be 80.
	Port string

	// Engine eg.: "mysql" or "postgres".
	EngineDB string

	// Host when is running the database Engine.
	HostDB string

	// Port of database Engine server.
	PortDB string

	// User of database, eg.: "root".
	UserDB string

	// Password of User database.
	PasswordDB string

	// Name of SQL database.
	NameDB string
}
