package dbmodel

// DataSource is database setting.
type DataSource struct {
	Driver   string
	Version  string
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Options  map[string]string
}

// NewDataSource returns new DataSource with initialized given arguments.
func NewDataSource(driver string, version string, host string, port int, user string, password string, database string, options map[string]string) DataSource {
	return DataSource{
		Driver:   driver,
		Version:  version,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
		Options:  options,
	}
}

// InitDataSource returns new DataSource that has initialized Options.
func InitDataSource() DataSource {
	return DataSource{Options: map[string]string{}}
}
