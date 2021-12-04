package types

type ResultSet interface{}

type Account struct {
	Username string
	Password string
}

type HostAddress struct {
	Host string
	Port int
}
