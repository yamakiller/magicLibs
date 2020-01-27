package logger

//LogMessage doc
//@Summary log message
//@Method LogMessage
//@Member (uint32) log level
//@Member (string) log prefix informat
//@Member (string) log message
type LogMessage struct {
	Level   uint32
	Prefix  string
	Message string
}
