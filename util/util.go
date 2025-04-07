package util

import (
	"context"
	"fmt"
	"time"

	"golang.design/x/clipboard"
)

func CopyToClipboard(content string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	select {
	case <-clipboard.Write(clipboard.FmtText, []byte(content)):
		return fmt.Errorf("A new value was written to the clipboard, the value produced is lost")
	case <-ctx.Done():
		return nil
	}
}
