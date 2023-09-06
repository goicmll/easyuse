package easyuse

import (
	"os"
	"os/signal"
	"syscall"
)

func IsDir(path string) bool {
	fInfo, err := os.Stat(path)
	if err == nil {
		return fInfo.IsDir()
	} else {
		// os.IsExist(nil) 永远等于false， os.IsNotExist(nil)也是
		return os.IsExist(err) && fInfo.IsDir()
	}

}

func IsRegular(path string) bool {
	fInfo, err := os.Stat(path)
	if err == nil {
		return fInfo.Mode().IsRegular()
	} else {
		// os.IsExist(nil) 永远等于false, os.IsNotExist(nil)也是
		return os.IsExist(err) && fInfo.Mode().IsRegular()
	}
}

var _ Hook = (*hook)(nil)

// Hook a graceful shutdown hook, default with signals of SIGINT and SIGTERM
type Hook interface {
	// WithSignals add more signals into hook
	WithSignals(signals ...syscall.Signal) Hook

	// Close register shutdown handles
	Close(funcs ...func())
}

type hook struct {
	ctx chan os.Signal
}

// NewHook create a Hook instance
func NewHook() Hook {
	hook := &hook{
		ctx: make(chan os.Signal, 1),
	}

	return hook.WithSignals(syscall.SIGINT, syscall.SIGTERM)
}

func (h *hook) WithSignals(signals ...syscall.Signal) Hook {
	for _, s := range signals {
		signal.Notify(h.ctx, s)
	}

	return h
}

func (h *hook) Close(funcs ...func()) {
	<-h.ctx
	signal.Stop(h.ctx)

	for _, f := range funcs {
		f()
	}
}
