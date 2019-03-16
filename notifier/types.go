package notifier

import "github.com/blueskan/gopheart/provider"

type Notifier interface {
	GetThreshold() int
	GetName() string
	Notify(statistics provider.Statistics) error
}
