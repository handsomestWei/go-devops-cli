package conf

type VersionConf struct {
	Version string `json:"version"`
	Meta    string `json:"meta"`
}

var AppVersionConf VersionConf
