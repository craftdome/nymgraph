package utils

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/errors"
	"golang.design/x/clipboard"
	"os"
	"time"
)

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
