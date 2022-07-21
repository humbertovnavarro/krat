package db

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const DAY = time.Hour * 24
const LOG_CLEANUP_INTERVAL = DAY * 7

var queue = make([]logrus.Entry, 0)
var lastClean = time.Now()

type SqliteHook struct{}

func NewSqliteHook() *SqliteHook {
	return &SqliteHook{}
}

func (hook *SqliteHook) Log(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read entry, %v", err)
		return err
	}
	_, err = DB.Exec(`INSERT INTO Logs (logLevel, logEntry, createdAt) VALUES (?, ?, ?)`, entry.Level, line, time.Now().UnixMilli())
	if err != nil {
		fmt.Printf("unable to log to sqlite db, %s \n", err)
		return err
	}
	return nil
}

func (hook *SqliteHook) Fire(entry *logrus.Entry) error {
	if DB == nil {
		queue = append(queue, *entry)
		return nil
	}
	queue = append(queue, *entry)
	if len(queue) > 0 {
		queueEntry := queue[0]
		hook.Log(&queueEntry)
		queue = queue[1:]
	}
	if lastClean.UnixMilli() < time.Now().UnixMilli()-LOG_CLEANUP_INTERVAL.Milliseconds() {
		CleanupLogs()
		lastClean = time.Now()
	}
	return nil
}

func (hook *SqliteHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
