package nebula

import (
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

type (
	Version           = types.Version
	AuthResponse      = types.AuthResponse
	ExecutionResponse = types.ExecutionResponse
)

var (
	versionAuto = Version("auto")

	Version2_5 = types.Version2_5
	Version2_6 = types.Version2_6
	Version3_0 = types.Version3_0

	supportedVersions = []Version{
		Version3_0,
		Version2_6,
		Version2_5,
	}
)
