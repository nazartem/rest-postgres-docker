package logging

import "log"

type Logger struct {
	Info *log.Logger
	Err  *log.Logger
}
