package spinner

import (
	"fmt"
	"io"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	DEFAULT_FRAME_RATE = time.Millisecond * 150
)

var DefaultCharst = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

type Spinner struct {
	Title     string
	Charset   []string
	Framerate time.Duration
	Output    io.Writer
	NoTty     bool
	runChan   chan struct{}
	stopOnce  sync.Once
}

func NewSpinner(title string) *Spinner {
	sp := &Spinner{
		Title:     title,
		Charset:   DefaultCharst,
		Framerate: DEFAULT_FRAME_RATE,
		runChan:   make(chan struct{}),
	}

	if !terminal.IsTerminal(syscall.Stdout) {
		sp.NoTty = true
	}

	return sp
}

func StartNew(title string) *Spinner {
	return NewSpinner(title).Start()
}

func (sp *Spinner) Start() *Spinner {
	go sp.writer()
	return sp
}

func (sp *Spinner) Stop() {
	sp.stopOnce.Do(func() {
		close(sp.runChan)
		sp.clearline()
	})
}

func (sp *Spinner) writer() {
	sp.animate()
	for {
		select {
		case <-sp.runChan:
			return
		default:
			sp.animate()
		}
	}
}

func (sp *Spinner) animate() {
	var out string
	for i := 0; i < len(sp.Charset); i++ {
		out = sp.Charset[i] + " " + sp.Title
		switch {
		case sp.Output != nil:
			fmt.Fprint(sp.Output, out)
		case !sp.NoTty:
			fmt.Print(out)
		}
		time.Sleep(sp.Framerate)
		sp.clearline()
	}
}

// https://en.wikipedia.org/wiki/ANSI_escape_code
func (sp *Spinner) clearline() {
	fmt.Printf("\033[2K") // erase the entire line
	fmt.Println()
	fmt.Printf("\033[1A") // move the cursor up 1 line
}
