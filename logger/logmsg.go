package logger

//LogMessage desc
//@method LogMessage desc: log message
//@member (uint32) log level
//@member (string) log prefix informat
//@member (string) log message
type LogMessage struct {
	_level   uint32
	_prefix  string
	_message string
}
