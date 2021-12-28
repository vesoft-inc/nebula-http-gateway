package nebula

import (
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
	"strconv"
	"strings"
)

func VersionHelper(address string, port int, username string, password string) (Version, error) {
	host := strings.Join([]string{address, strconv.Itoa(port)}, ":")
	for _, version := range []Version{
		V2_5,
		V2_6,
		V3_0,
	} {
		c, err := NewClient(ConnectionInfo{
			GraphEndpoints: []string{host},
			GraphAccount: Account{
				Username: username,
				Password: password,
			},
		}, WithVersion(version))
		if err == nil {
			err = c.Graph().Open()

			if err == nil {
				return version, c.Graph().Close()
			}
		}
	}
	return "", errors.ErrVersionEstimateFailed
}
