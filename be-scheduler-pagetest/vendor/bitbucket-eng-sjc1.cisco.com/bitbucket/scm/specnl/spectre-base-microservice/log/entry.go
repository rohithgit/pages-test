package log // spectre

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"    // spectre
	"runtime" // spectre
	"strings" // spectre
	"time"
)

// Defines the key when adding errors using WithError.
var ErrorKey = "error"

// An entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Debug, Info,
// Warn, Error, Fatal or Panic is called on it. These objects can be reused and
// passed around as much as you wish to avoid field duplication.
type Entry struct {
	Logger *Logger

	// Contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Debug, Info, Warn, Error, Fatal or Panic
	Level Level

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// Stack depth - used when getting the file and line number
	StackDepth int // spectre

	// Filename of source file
	Filename string // spectre
	// Line number in source file
	Lineno int // spectre
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		// Default is three fields, give a little extra room
		Data:       make(Fields, 5),
		StackDepth: logger.StackDepth, // spectre
	}
}

func (entry *Entry) incrStackDepth(val int) { // spectre
	entry.StackDepth += val // spectre
} // spectre

func (entry *Entry) decrStackDepth(val int) { // spectre
	entry.StackDepth -= val // spectre
} // spectre

func (entry *Entry) resetStackDepth() { // spectre
	entry.StackDepth = STACKDEPTH // spectre
} // spectre

func (entry *Entry) getLineInfo() { // spectre
	_, file, line, ok := runtime.Caller(entry.StackDepth) // spectre
	if !ok {                                              // spectre
		file = "???" // spectre
		line = 0     // spectre
	} // spectre
	base := path.Base(file)            // spectre
	dir := path.Dir(file)              // spectre
	folders := strings.Split(dir, "/") // spectre
	parent := ""                       // spectre
	if folders != nil {                // spectre
		parent = folders[len(folders)-1] + "/" // spectre
	} // spectre
	file = parent + base  // spectre
	entry.Filename = file // spectre
	entry.Lineno = line   // spectre
}

// Returns a reader for the entry, which is a proxy to the formatter.
func (entry *Entry) Reader() (*bytes.Buffer, error) {
	serialized, err := entry.Logger.Formatter.Format(entry)
	return bytes.NewBuffer(serialized), err
}

// Returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) {
	reader, err := entry.Reader()
	if err != nil {
		return "", err
	}

	return reader.String(), err
}

// Add an error as single field (using the key defined in ErrorKey) to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField(ErrorKey, err)
}

// Add a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Fields{key: value})
}

// Add a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	data := Fields{}
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: data, StackDepth: entry.Logger.StackDepth - 2} // spectre
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func (entry Entry) log(level Level, msg string) {
	entry.Time = time.Now()
	entry.Level = level
	entry.Message = msg
	entry.getLineInfo() // spectre

	if err := entry.Logger.Hooks.Fire(level, &entry); err != nil {
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
		entry.Logger.mu.Unlock()
	}

	reader, err := entry.Reader()
	if err != nil {
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		entry.Logger.mu.Unlock()
	}

	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()

	_, err = io.Copy(entry.Logger.Out, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel {
		panic(&entry)
	}
}

func (entry *Entry) Debug(args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.log(DebugLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Print(args ...interface{}) {
	entry.Info(args...)
}

func (entry *Entry) Info(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Warn(args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Warning(args ...interface{}) {
	entry.Warn(args...)
}

func (entry *Entry) Error(args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.log(ErrorLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Fatal(args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.log(FatalLevel, fmt.Sprint(args...))
	}
	os.Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.log(PanicLevel, fmt.Sprint(args...))
	}
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.incrStackDepth(1) // spectre
		entry.Debug(fmt.Sprintf(format, args...))
		entry.resetStackDepth() // spectre
	}
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.incrStackDepth(1) // spectre
		entry.Info(fmt.Sprintf(format, args...))
		entry.resetStackDepth() // spectre
	}
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.incrStackDepth(1) // spectre
	entry.Infof(format, args...)
	entry.resetStackDepth() // spectre
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.incrStackDepth(1) // spectre
		entry.Warn(fmt.Sprintf(format, args...))
		entry.resetStackDepth() // spectre
	}
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.incrStackDepth(1) // spectre
	entry.Warnf(format, args...)
	entry.resetStackDepth() // spectre
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.incrStackDepth(1) // spectre
		entry.Error(fmt.Sprintf(format, args...))
		entry.resetStackDepth() // spectre
	}
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.incrStackDepth(1) // spectre
		entry.Fatal(fmt.Sprintf(format, args...))
		entry.resetStackDepth() // spectre
	}
	os.Exit(1)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.incrStackDepth(1) // spectre
		entry.Panic(fmt.Sprintf(format, args...))
		entry.resetStackDepth() // spectre
	}
}

// Entry Println family functions

func (entry *Entry) Debugln(args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.incrStackDepth(1) // spectre
		entry.Debug(entry.sprintlnn(args...))
		entry.resetStackDepth() // spectre
	}
}

func (entry *Entry) Infoln(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.incrStackDepth(1) // spectre
		entry.Info(entry.sprintlnn(args...))
		entry.resetStackDepth() // spectre
	}
}

func (entry *Entry) Println(args ...interface{}) {
	entry.incrStackDepth(1) // spectre
	entry.Infoln(args...)
	entry.resetStackDepth() // spectre
}

func (entry *Entry) Warnln(args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.incrStackDepth(1) // spectre
		entry.Warn(entry.sprintlnn(args...))
		entry.resetStackDepth() // spectre
	}
}

func (entry *Entry) Warningln(args ...interface{}) {
	entry.incrStackDepth(1) // spectre
	entry.Warnln(args...)
	entry.resetStackDepth() // spectre
}

func (entry *Entry) Errorln(args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.incrStackDepth(1) // spectre
		entry.Error(entry.sprintlnn(args...))
		entry.resetStackDepth() // spectre
	}
}

func (entry *Entry) Fatalln(args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.incrStackDepth(1) // spectre
		entry.Fatal(entry.sprintlnn(args...))
		entry.resetStackDepth() // spectre
	}
	os.Exit(1)
}

func (entry *Entry) Panicln(args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.incrStackDepth(1) // spectre
		entry.Panic(entry.sprintlnn(args...))
		entry.resetStackDepth() // spectre
	}
}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func (entry *Entry) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
