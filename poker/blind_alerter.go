package poker

import (
	"fmt"
	"io"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int, to io.Writer)
}

// an example implementation of Alerter
type BlindAlerterFunc func(duration time.Duration, amount int, to io.Writer)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	a(duration, amount, to)
}

func Alerter(duration time.Duration, amount int, to io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(to, "Blind is now %d\n", amount)
	})
}
