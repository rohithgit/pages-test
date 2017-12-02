package logging

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

// This is the struct use to expose logrus logger
type SpectreLog struct {
	Logger *logrus.Logger
}

// Function NewSpectreLog creates a new insatnce of struct SpectreLog
func NewSpectreLog(logger *logrus.Logger) *SpectreLog {
	return &SpectreLog{
		Logger: logger,
	}
}

// This is a global instance to set the logger
var Log = NewSpectreLog(logrus.New())

/* init() functions initializes the logger with default values */

func init() {
	// Log as logstash instead of the default ASCII formatter.
	Log.Logger.Formatter = new(logrus.TextFormatter)

	// Output to stdout, could also be a file.
	Log.Logger.Out = os.Stdout

	// Only log the warning severity or above.
	Log.Logger.Level = logrus.InfoLevel
}

/* LoggingInit functions initializes the logger with values provided by user */

func (log *SpectreLog) LoggingInit(input_format string, input_logoutput io.Writer, input_level string) error {
	// Checking the input log format and setting the formatter accordingly
	switch input_format {
	case "Text":
		Log.Logger.Formatter = new(logrus.TextFormatter)
		logrus.SetFormatter(Log.Logger.Formatter)
	case "Json":
		Log.Logger.Formatter = new(logrus.JSONFormatter)
		logrus.SetFormatter(Log.Logger.Formatter)
	default:
		return errors.New("Incorrect log format provided")
	}

	// If input log output provided is nil we set default value as os.Stdout
	if input_logoutput == nil {
		input_logoutput = os.Stdout
	}

	// Set log output as per the user input provided, could also be a file.
	Log.Logger.Out = input_logoutput

	// If input log level provided is empty string we set default value as Warn
	if input_level == "" {
		input_level = "Warn"
	}

	// Checking the input log level and setting the level accordingly
	switch input_level {
	case "Debug":
		Log.Logger.Level = logrus.DebugLevel
	case "Warn", "Warning":
		Log.Logger.Level = logrus.WarnLevel
	case "Info":
		Log.Logger.Level = logrus.InfoLevel
	case "Error":
		Log.Logger.Level = logrus.ErrorLevel
	case "Fatal":
		Log.Logger.Level = logrus.FatalLevel
	case "Panic":
		Log.Logger.Level = logrus.PanicLevel
	default:
		return errors.New("Incorrect log level provided")
	}

	return nil
}

/* LoggingWithFields function takes input as :
   - level (type string, expected values :
          Debug, Info, Warn, Error, Fatal, Panic),
   - logmsg (type string, expected value will contain following format
     microservice_name tenant_id transaction_id log_message).
 This function sets the log entry and prints it to the defined log output */

func (log *SpectreLog) GenerateLog(level string, microservice_name string, logmsg interface{}) (*logrus.Entry, error) {
	var msg string
	// Check if proper inputs were provided
	if level == "" || logmsg == nil || microservice_name == "" {
		return nil, errors.New("Not enough input provided")
	}
	// Get the current time in  proper format
	current_time := time.Now().Format("2006-01-02 15:04:05")
	// Set time and logmsg in proper format
	msg = fmt.Sprintf(
		"%s %s:%s %s",
		current_time,
		microservice_name,
		fileInfo(2),
		logmsg,
	)

	// Based on the log level print the logs
	switch level {
	case "Debug":
		result := Log.SetEntry(logrus.DebugLevel, logmsg)
		result.Debug(msg)
		return result, nil
	case "Info":
		result := Log.SetEntry(logrus.InfoLevel, logmsg)
		result.Info(msg)
		return result, nil
	case "Warn", "Warning":
		result := Log.SetEntry(logrus.WarnLevel, logmsg)
		result.Warn(msg)
		return result, nil
	case "Error":
		result := Log.SetEntry(logrus.ErrorLevel, logmsg)
		result.Error(msg)
		return result, nil
	case "Fatal":
		result := Log.SetEntry(logrus.FatalLevel, logmsg)
		result.Fatal(msg)
		return result, nil
	case "Panic":
		result := Log.SetEntry(logrus.PanicLevel, logmsg)
		result.Panic(msg)
		return result, nil
	default:
		return nil, errors.New("Incorrect Log Level provided")
	}

	return nil, nil
}

/* SetEntry function is used to set the log entry based on the data provided by the user */

func (log *SpectreLog) SetEntry(level logrus.Level, logmsg interface{}) *logrus.Entry {
	entry := logrus.NewEntry(Log.Logger)
	entry.Level = level
	entry.Message = logmsg.(string)
	return entry
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
