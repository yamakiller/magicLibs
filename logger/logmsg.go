package logger

//LogMessage desc
//@Method LogMessage desc: log message
//@Member (uint32) log level
//@Member (string) log prefix informat
//@Member (string) log message
type LogMessage struct {
	_level   uint32
	_prefix  string
	_message string
}
