package factory

import (
	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/notifier"
)

type NotifierFactory interface {
	CreateNotifier(name, notifierType string, config config.NotifierService) notifier.Notifier
}
