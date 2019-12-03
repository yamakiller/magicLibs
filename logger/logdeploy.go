package logger

//LogDeploy doc
//@Struct LogDeploy @Summary logger config informat
type LogDeploy struct {
	LogPath  string `json:"log-path"`
	LogSize  int    `json:"log-size"`
	LogLevel uint32 `json:"log-level"`
}

//NewDefault doc
//@Method NewDefault @Summary create default value Logger deploy informat
func NewDefault() *LogDeploy {
	return &LogDeploy{LogPath: "", LogSize: 512, LogLevel: DEBUGLEVEL}
}
