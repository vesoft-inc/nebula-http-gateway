package clientv2_5_1

import (
	"errors"

	nebula "github.com/vesoft-inc/nebula-go/v2"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/client/types"
)

type Session struct {
	session *nebula.Session
}

func (s *Session) Release() {
	s.session.Release()
}

func (s *Session) Execute(stmt string) (types.ResultSet, error) {
	resultSet, err := s.session.Execute(stmt)
	if err != nil {
		return nil, err
	}
	return resultSet, nil
}

func (s *Session) ExecuteJson(stmt string) ([]byte, error) {
	return nil, errors.New("execute json not support")
}
