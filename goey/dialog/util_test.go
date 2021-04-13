package dialog

import (
	"flag"
	"time"
)

var (
	asyncWait = 1000 * time.Millisecond
)

type durationValue time.Duration

func (tv durationValue) String() string {
	return time.Duration(tv).String()
}

func (tv *durationValue) Set(s string) error {
	value, err := time.ParseDuration(s)
	if err == nil {
		*(*time.Duration)(tv) = value
	}
	return err
}

func init() {
	flag.Var((*durationValue)(&asyncWait), "async-wait", "Set delay before typing keys")
}

func asyncKeyEnter() {
	asyncTypeKeys("\n", asyncWait)
}

func asyncKeyEscape() {
	asyncTypeKeys("\x1b", asyncWait)
}
