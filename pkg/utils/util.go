package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/errors"
	"golang.design/x/clipboard"
	"io"
	"os"
	"strings"
	"syscall"
	"time"
)

func PrintJson(data []byte) error {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, data, "", "\t")
	if err != nil {
		return errors.Wrap(err, "Indent")
	}

	fmt.Println(string(prettyJSON.Bytes()))
	return nil
}

func SaveResource(fileName string, bin []byte) error {
	if _, err := os.Stat(fileName); err != nil {
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			return errors.Wrap(err, "OpenFile")
		}
		defer file.Close()

		if _, err := file.Write(bin); err != nil {
			return errors.Wrap(err, "Write")
		}
	}

	return nil
}

func CopyToClipboard(text string) error {
	if err := clipboard.Init(); err != nil {
		return err
	}
	clipboard.Write(clipboard.FmtText, []byte(text))

	return nil
}

func ShowSplash(text any) {
	if drv, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
		w := drv.CreateSplashWindow()
		w.SetIcon(theme.ContentCopyIcon())
		w.SetContent(widget.NewRichTextWithText(fmt.Sprintf("%v", text)))
		w.Show()
		go func() {
			time.Sleep(1 * time.Second)
			w.Close()
		}()
	}
}

func LinesFromFile(path string) ([]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "OpenFile")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "ReadAll")
	}

	size := bytes.Count(data, []byte{'\n'})

	lines := make([]string, 0, size)
	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		lines = append(lines, strings.TrimSpace(line))
	}

	return lines, nil
}

func EnableWindowsConsoleColors() {
	// Try to make ANSI work
	handle := syscall.Handle(os.Stdout.Fd())
	kernel32DLL := syscall.NewLazyDLL("kernel32.dll")
	setConsoleModeProc := kernel32DLL.NewProc("SetConsoleMode")
	// If it fails, fallback to no colors
	setConsoleModeProc.Call(uintptr(handle), 0x0001|0x0002|0x0004)
}

func Reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func ContainsString(a []string, x string) bool {
	for i := 0; i < len(a); i++ {
		if a[i] == x {
			return true
		}
	}

	return false
}
