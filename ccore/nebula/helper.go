package nebula

import (
	"strconv"
	"strings"

	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
)

func VersionHelper(address string, port int, username string, password string) (Version, error) {
	host := strings.Join([]string{address, strconv.Itoa(port)}, ":")
	for _, version := range []Version{
		V3_0,
		V2_6,
		V2_5,
	} {
		c, err := NewClient(ConnectionInfo{
			GraphEndpoints: []string{host},
			GraphAccount: Account{
				Username: username,
				Password: password,
			},
		}, WithVersion(version))
		if err == nil {
			if c.Graph().Open() == nil {
				return version, c.Graph().Close()
			}
		}
	}
	return "", errors.ErrUnsupportedVersion
}
