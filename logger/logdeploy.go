package logger

//LogDeploy doc
//@Struct LogDeploy @Summary logger config informat
type LogDeploy struct {
	LogPath  string `xml:"log-path" yaml:"log-path" json:"log-path"`
	LogSize  int    `xml:"log-size" yaml:"log-size" json:"log-size"`
	LogLevel uint32 `xml:"log-level" yaml:"log-level" json:"log-level"`
}

//NewDefault doc
//@Method NewDefault @Summary create default value Logger deploy informat
func NewDefault() *LogDeploy {
	return &LogDeploy{LogPath: "", LogSize: 512, LogLevel: DEBUGLEVEL}
}
