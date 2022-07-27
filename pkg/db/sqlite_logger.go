package db

import (
	"fmt"
	"os"
	"time"

	"github.com/humbertovnavarro/krat/pkg/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const LOG_CLEANUP_INTERVAL = time.Minute

var queue = make([]logrus.Entry, 0)
var lastClean = time.Now()

type SqliteHook struct {
	db *gorm.DB
}

func NewSqliteHook(db *gorm.DB) *SqliteHook {
	return &SqliteHook{
		db,
	}
}

func (hook *SqliteHook) Log(entry *logrus.Entry) error {
	line, err := entry.String()
	hook.db.Create(&models.Log{
		Level: uint32(entry.Level),
		Text:  line,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read entry, %v", err)
		return err
	}
	if err != nil {
		fmt.Printf("unable to log to sqlite db, %s \n", err)
		return err
	}
	return nil
}

func (hook *SqliteHook) Fire(entry *logrus.Entry) error {
	if hook.db == nil {
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
		lastClean = time.Now()
		hook.Cleanup()
	}
	return nil
}

func (hook *SqliteHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *SqliteHook) Cleanup() {
	hook.db.Delete(models.Log{}, "CreatedAt < ?", time.Now().UnixMilli()-LOG_CLEANUP_INTERVAL.Milliseconds())
}
