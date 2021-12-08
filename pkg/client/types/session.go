package types

type Session interface {
	Release()
	Execute(stmt string) (ResultSet, error)
	ExecuteJson(stmt string) ([]byte, error)
}
