package logger

type LogDeploy struct {
	LogPath  string `json:"log-path"`
	LogSize  int    `json:"log-size"`
	LogLevel int    `json:"log-level"`
}

func NewDefault() *LogDeploy {
	return &LogDeploy{LogPath: "", LogSize: 256, LogLevel: PANICLEVEL}
}
