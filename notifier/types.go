package notifier

import "github.com/blueskan/gopheart/provider"

type Notifier interface {
	Notify(statistics provider.Statistics) error
}
