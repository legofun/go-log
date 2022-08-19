package golog

import (
	"bytes"
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LevelFlag log level flags
const (
	LevelFlagAccess      = "ACCESS"
	LevelFlagError       = "ERROR"
	LevelFlagTransaction = "TRANSACTION"
	LevelFlagAbnormal    = "ABNORMAL"
	LevelFlagCommon      = "COMMON"
	LevelFlagDebug       = "DEBUG"
)

//message levels.
const (
	ACCESS      = iota // Access: access log
	ERROR              // Error: error conditions
	TRANSACTION        // Transaction: transaction messages
	ABNORMAL           // Abnormal: abnormal condition
	COMMON             // Common: common messages
	DEBUG              // Debug: debug-level messages
)

const (
	// default size or min size for record channel
	recordChannelSizeDefault = uint(4096)
	// default time layout
	defaultLayout = "2006/01/02 15:04:05"
	// timestamp with zone info
	timestampLayout = "2006-01-02T15:04:05.000+0800"
)

// LevelFlags level Flags set
var (
	LevelFlags = []string{
		LevelFlagAccess,
		LevelFlagTransaction,
		LevelFlagError,
		LevelFlagAbnormal,
		LevelFlagCommon,
		LevelFlagDebug,
	}
	DefaultLayout = defaultLayout
)

// default logger
var (
	loggerDefault     *Logger
	recordPool        *sync.Pool
	recordChannelSize = recordChannelSizeDefault // log chan size
)

// Record log record
type Record struct {
	level int
	time  string
	file  string
	msg   string
}

func (r *Record) String() string {
	return fmt.Sprintf("%s [%s] <%s> %s\n", r.time, LevelFlags[r.level], r.file, r.msg)
}

// Writer record writer
type Writer interface {
	Init() error
	Write(*Record) error
}

// Flusher record flusher
type Flusher interface {
	Flush() error
}

// Rotater record rotater
type Rotater interface {
	Rotate() error
	SetPathPattern(string) error
}

// Logger logger define
type Logger struct {
	writers         []Writer
	records         chan *Record
	recordsChanSize uint
	lastTime        int64
	lastTimeStr     string

	flushTimer  time.Duration // timer to flush logger record to chan
	rotateTimer time.Duration // timer to rotate logger record for writer

	c chan bool

	layout       string
	level        int
	fullPath     bool // show full path, default only show file:line_number
	withFuncName bool // show caller func name
	lock         sync.RWMutex
}

// NewLogger create the logger
func NewLogger() *Logger {
	if loggerDefault != nil {
		return loggerDefault
	}
	records := make(chan *Record, recordChannelSize)

	return newLoggerWithRecords(records)
}

// newLoggerWithRecords is useful for go test
func newLoggerWithRecords(records chan *Record) *Logger {
	l := new(Logger)
	l.writers = make([]Writer, 0, 1) // normal least has console writer
	if l.recordsChanSize == 0 {
		recordChannelSize = recordChannelSizeDefault
	}

	l.records = records
	l.c = make(chan bool, 1)
	l.level = DEBUG
	l.layout = DefaultLayout

	go bootstrapLogWriter(l)

	return l
}

// Register register writer
// the writer should be register once for writers by kind
func (l *Logger) Register(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}

	l.writers = append(l.writers, w)
}

// Close close logger
func (l *Logger) Close() {
	close(l.records)
	<-l.c

	for _, w := range l.writers {
		if f, ok := w.(Flusher); ok {
			if err := f.Flush(); err != nil {
				log.Println(err)
			}
		}
	}
}

// SetLayout set the logger time layout
func (l *Logger) SetLayout(layout string) {
	l.layout = layout
}

// SetLevel set the logger level
func (l *Logger) SetLevel(lvl int) {
	l.level = lvl
}

// WithFullPath set the logger with full path
func (l *Logger) WithFullPath(show bool) {
	l.fullPath = show
}

// WithFuncName set the logger with func name
func (l *Logger) WithFuncName(show bool) {
	l.withFuncName = show
}

// Debug level
func (l *Logger) Debug(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}

// Common level
func (l *Logger) Common(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(COMMON, fmt, args...)
}

// Abnormal level
func (l *Logger) Abnormal(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ABNORMAL, fmt, args...)
}

// Transaction level
func (l *Logger) Transaction(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(TRANSACTION, fmt, args...)
}

// Error level
func (l *Logger) Error(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, fmt, args...)
}

// Access level
func (l *Logger) Access(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ACCESS, fmt, args...)
}

func (l *Logger) deliverRecordToWriter(level int, f string, args ...interface{}) {
	var msg string
	var fi bytes.Buffer

	if level > l.level {
		return
	}

	msg = f
	sz := len(args)
	if sz != 0 {
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
		} else {
			msg += strings.Repeat("%v", len(args))
		}
	}
	msg = fmt.Sprintf(msg, args...)

	// source code, file and line num
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		fileName := path.Base(file)
		if l.fullPath {
			fileName = file
		}
		fi.WriteString(fmt.Sprintf("%s:%d", fileName, line))

		if l.withFuncName {
			funcName := runtime.FuncForPC(pc).Name()
			funcName = path.Base(funcName)
			fi.WriteString(fmt.Sprintf(" %s", funcName))
		}
	}

	// format time
	now := time.Now()
	l.lock.Lock() // avoid data race
	if now.Unix() != l.lastTime {
		l.lastTime = now.Unix()
		l.lastTimeStr = now.Format(l.layout)
	}
	lastTimeStr := l.lastTimeStr
	l.lock.Unlock()

	r := recordPool.Get().(*Record)
	r.msg = msg
	r.file = fi.String()
	r.time = lastTimeStr
	r.level = level

	l.records <- r
}

func bootstrapLogWriter(logger *Logger) {
	var (
		r  *Record
		ok bool
	)

	if r, ok = <-logger.records; !ok {
		logger.c <- true
		return
	}

	for _, w := range logger.writers {
		if err := w.Write(r); err != nil {
			log.Printf("%v\n", err)
		}
	}

	flushTimer := time.NewTimer(logger.flushTimer)
	rotateTimer := time.NewTimer(logger.rotateTimer)

	for {
		select {
		case r, ok = <-logger.records:
			if !ok {
				logger.c <- true
				return
			}

			for _, w := range logger.writers {
				if err := w.Write(r); err != nil {
					log.Printf("%v\n", err)
				}
			}

			recordPool.Put(r)

		case <-flushTimer.C:
			for _, w := range logger.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						log.Printf("%v\n", err)
					}
				}
			}
			flushTimer.Reset(logger.flushTimer)

		case <-rotateTimer.C:
			for _, w := range logger.writers {
				if r, ok := w.(Rotater); ok {
					if err := r.Rotate(); err != nil {
						log.Printf("%v\n", err)
					}
				}
			}
			rotateTimer.Reset(logger.rotateTimer)
		}
	}
}

func init() {
	loggerDefault = NewLogger()
	loggerDefault.flushTimer = time.Millisecond * 500
	loggerDefault.rotateTimer = time.Second * 10
	recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
}

// Register register writer
func Register(w Writer) {
	loggerDefault.Register(w)
}

// Close close logger
func Close() {
	loggerDefault.Close()
}

// SetLayout set the logger time layout, should call before logger real use
func SetLayout(layout string) {
	loggerDefault.layout = layout
}

// SetLevel set the logger level, should call before logger real use
func SetLevel(lvl int) {
	loggerDefault.level = lvl
}

// WithFullPath set the logger with full path, should call before logger real use
func WithFullPath(show bool) {
	loggerDefault.fullPath = show
}

// WithFuncName set the logger with func name, should call before logger real use
func WithFuncName(show bool) {
	loggerDefault.withFuncName = show
}

// Debug level
func Debug(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(DEBUG, fmt, args...)
}

// Common level
func Common(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(COMMON, fmt, args...)
}

// Abnormal level
func Abnormal(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(ABNORMAL, fmt, args...)
}

// Transaction level
func Transaction(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(TRANSACTION, fmt, args...)
}

// Error level
func Error(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(ERROR, fmt, args...)
}

// Access level
func Access(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(ACCESS, fmt, args...)
}

// The method is put here, so it's easy to test
func getLevelDefault(flag string, defaultFlag int, writer string) int {
	for i, f := range LevelFlags {
		if strings.TrimSpace(strings.ToUpper(flag)) == f {
			return i
		}
	}
	log.Printf("[golog] no matching level for writer(%v, flag:%v), use default level(%d, flag:%v)", writer, flag, defaultFlag, LevelFlags[defaultFlag])
	return defaultFlag
}
