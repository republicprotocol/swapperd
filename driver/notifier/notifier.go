package notifier

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/fsnotify.v1"
)

type notifier struct {
	logger logrus.FieldLogger
}

func New(logger logrus.FieldLogger) *notifier {
	return &notifier{
		logger: logger,
	}
}

func (notifier *notifier) Watch(done <-chan struct{}, files ...string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		notifier.logger.Fatal(err)
	}
	defer watcher.Close()

	for _, file := range files {
		if err := watcher.Add(file); err != nil {
			notifier.logger.Fatal(err)
		}
	}

	for {
		select {
		case <-done:
			return
		case _, ok := <-watcher.Events:
			if !ok {
				return
			}
			os.Exit(1)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			notifier.logger.Fatal(err)
		}
	}
}
