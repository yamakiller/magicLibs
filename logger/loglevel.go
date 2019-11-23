package logger

const (
	// PANICLEVEL Crash log level
	PANICLEVEL uint32 = iota
	// FATALLEVEL Critical error log level
	FATALLEVEL
	// ERRORLEVEL Error log level
	ERRORLEVEL
	// WARNLEVEL  Warning log level
	WARNLEVEL
	// INFOLEVEL  General information log level
	INFOLEVEL
	// DEBUGLEVEL Debug log level
	DEBUGLEVEL
	// TRACELEVEL Trace log level
	TRACELEVEL
)
