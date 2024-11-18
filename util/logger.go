package util

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	Prog                = "todoister"
	DefaultLogFile      = "out.log"
	DefaultLogFileCount = 3
	DefaultLogFileSize  = 128 * 1024
)

var (
	tty     = isatty()
	systemd = isSystemdService()
)

func isatty() bool {
	_, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	return err == nil
}

func isSystemdService() bool {
	_, ok := os.LookupEnv("SYSTEMD_EXEC_PID")
	return ok
}

type FileLogWriter struct {
	logfile      string
	logFileSize  int64
	logFileCount int
	currentFile  *os.File
	currentSize  int64
}

func (w *FileLogWriter) Write(p []byte) (n int, err error) {
	if w.currentSize+int64(len(p)) > w.logFileSize {
		err := w.rotate()
		if err != nil {
			return 0, err
		}
	}

	n, err = w.currentFile.Write(p)
	w.currentSize += int64(n)
	return n, err
}

func (w *FileLogWriter) rotate() error {
	if w.currentFile != nil {
		err := w.currentFile.Close()
		if err != nil {
			return err
		}
	}

	// Rotate old log files.
	for i := w.logFileCount - 1; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", w.logfile, i)
		newPath := fmt.Sprintf("%s.%d", w.logfile, i+1)
		if _, err := os.Stat(oldPath); err == nil {
			err := os.Rename(oldPath, newPath)
			if err != nil {
				return err
			}
		}
	}

	// Rename current log file.
	if _, err := os.Stat(w.logfile); err == nil {
		err := os.Rename(w.logfile, fmt.Sprintf("%s.1", w.logfile))
		if err != nil {
			return err
		}
	}

	// Create a new log file.
	file, err := os.Create(w.logfile)
	if err != nil {
		return err
	}

	w.currentFile = file
	w.currentSize = 0
	return nil
}

func NewFileLogWriter(logfile string, logFileSize int64, logFileCount int) (*FileLogWriter, error) {
	writer := &FileLogWriter{
		logfile:      logfile,
		logFileSize:  logFileSize,
		logFileCount: logFileCount,
	}
	// Resume writing to the current log file.
	if fileInfo, err := os.Stat(logfile); err == nil {
		writer.currentSize = fileInfo.Size()
		writer.currentFile, err = os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
	}
	return writer, nil
}

func InitLogger(logfile string) {
	// No need to initialize the logger if running in a terminal or as a systemd service.
	if tty || systemd {
		return
	}

	var logDir string

	// If logfile is explicitly set, use it.
	if logfile != "" {
		logDir = filepath.Dir(logfile)
	} else {
		// Otherwise, use the user cache directory $HOME/.cache.
		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			Die("Error getting user cache directory", err)
		}
		logDir = filepath.Join(userCacheDir, Prog)
		logfile = filepath.Join(logDir, DefaultLogFile)
	}

	// Create the log directory if it doesn't exist.
	if _, err := os.Stat(logDir); err != nil {
		if err = os.MkdirAll(logDir, 0755); err != nil {
			Die("Error creating log directory", err)
		}
	}
	writer, err := NewFileLogWriter(logfile, DefaultLogFileSize, DefaultLogFileCount)
	if err != nil {
		Die("Error creating log file writer", err)
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(writer, nil)))
}

func writeLogEntry(level slog.Level, text string, e error) {
	var formattedText string
	if e != nil {
		formattedText = fmt.Sprintf("%s: %v", text, e)
	} else {
		formattedText = text
	}
	if tty || systemd {
		if level == slog.LevelInfo {
			_, _ = fmt.Fprintln(os.Stdout, formattedText)
		} else {
			_, _ = fmt.Fprintln(os.Stderr, formattedText)
		}
	} else {
		switch level {
		case slog.LevelInfo:
			slog.Info(formattedText)
		case slog.LevelWarn:
			slog.Warn(formattedText)
		case slog.LevelError:
			slog.Error(formattedText)
		}
	}
}

func Info(text string) {
	writeLogEntry(slog.LevelInfo, text, nil)
}

func Warn(text string, e error) {
	writeLogEntry(slog.LevelWarn, text, e)
}

func Error(text string, e error) {
	writeLogEntry(slog.LevelError, text, e)
}

func Die(text string, e error) {
	writeLogEntry(slog.LevelError, text, e)
	os.Exit(1)
}
