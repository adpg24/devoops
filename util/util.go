package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/go-ini/ini"
	"golang.design/x/clipboard"
)

const (
	awsConfigDir  = ".aws"
	awsConfigFile = "credentials"
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

func Insert[T any](array []T, element T, i int) []T {
	return append(array[:i], append([]T{element}, array[i:]...)...)
}

func HandleErr(err error, msg string, args ...any) {
	if err != nil {
		log.Fatalf(msg, args...)
	}
}

func AddProfileSection(saveTo string, iniFile *ini.File, sectionName string, keys map[string]string) error {
	sec := iniFile.Section(sectionName)
	for k, v := range keys {
		_, err := sec.NewKey(k, v)
		if err != nil {
			return err
		}
	}

	err := iniFile.SaveTo(saveTo)
	if err != nil {
		return err
	}
	return nil
}

func GetAwsConfigDir() string {
	home, err := os.UserHomeDir()
	HandleErr(err, "Failed to retrieve home dir: %v", err)
	return path.Join(home, awsConfigDir, awsConfigFile)
}
