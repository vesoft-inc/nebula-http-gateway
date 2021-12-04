package clientv2_6_0

import (
	nebula "github.com/vesoft-inc/nebula-go/v2"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/types"
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
	resultAsBytes, err := s.session.ExecuteJson(stmt)
	if err != nil {
		return nil, err
	}
	return resultAsBytes, nil
}
