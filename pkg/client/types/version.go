package types

type Version int

const (
	VersionNotSupport = iota
	Version_2_0_0_ga
	Version_2_5_0
	Version_2_5_1
	Version_2_6_0
)

var versionMap = map[Version]string{
	VersionNotSupport: "VersionNotSupport",
	Version_2_0_0_ga:  "v2.0.0-ga",
	Version_2_5_0:     "v2.5.0",
	Version_2_5_1:     "v2.5.1",
	Version_2_6_0:     "v2.6.0",
}

var versionRevMap = map[string]Version{
	"VersionNotSupport": VersionNotSupport,
	"v2.0.0-ga":         Version_2_0_0_ga,
	"v2.5.0":            Version_2_5_0,
	"v2.5.1":            Version_2_5_1,
	"v2.6.0":            Version_2_6_0,
}

func NewVersion(version string) Version {
	if v, ok := versionRevMap[version]; ok {
		return v
	}
	return VersionNotSupport
}

func (version Version) String() string {
	if v, ok := versionMap[version]; ok {
		return v
	}
	return "VersionNotSupport"
}
