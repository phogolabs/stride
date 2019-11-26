package cmd

import (
	"fmt"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/stride/torrent"
)

func get(ctx *cli.Context, key string) (string, error) {
	// get the spec async
	task, err := torrent.GetAsync(ctx.String(key))
	if err != nil {
		return "", err
	}

	// make the task available within the app
	ctx.Metadata["task"] = task

	// wait the download to finish
	if err = task.Wait(); err != nil {
		// if there are some error stop
		return "", err
	}

	path := fmt.Sprintf("%v", task.Data())
	return path, nil
}
