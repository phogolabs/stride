package torrent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-getter"
	"github.com/phogolabs/stride/async"
)

// GetAsync downloads a given file
func GetAsync(path string) (*async.Task, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var (
		epoch = fmt.Sprintf("%v", time.Now().Unix())
		dst   = filepath.Join(os.TempDir(), epoch)
		file  = filepath.Join(dst, filepath.Base(path))
	)

	fn := func(ctx context.Context) error {
		client := &getter.Client{
			Ctx:     ctx,
			Pwd:     pwd,
			Src:     path,
			Dst:     dst,
			Mode:    getter.ClientModeAny,
			Options: []getter.ClientOption{},
		}

		return client.Get()
	}

	task := async.NewTask(fn, file)
	task.Run()

	return task, nil
}
