package config

import (
	"github.com/hpcloud/tail"
	"logagent/tailhandler"
)

func InitTail() error {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}

	file, err := tail.TailFile(Cfg.Collect.LogFilePath, config)
	if err != nil {
		return err
	}

	tailhandler.TailHandler = file

	return nil
}
