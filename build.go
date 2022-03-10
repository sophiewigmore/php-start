package phpstart

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate faux --interface ProcMgr --output fakes/procmgr.go

// ProcMgr
type ProcMgr interface {
	Add(name string, proc Proc)
	WriteFile(path string) error
}

func Build(procs ProcMgr, logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		// TODO: add logging
		// TODO: add code comments
		// TODO: add failure case tests
		logger.Process("START BUILDPACK")

		layer, err := context.Layers.Get("php-start")
		if err != nil {
			panic(err)
		}
		layer, err = layer.Reset()
		if err != nil {
			panic(err)
		}
		layer.Launch = true

		httpdConfPath := os.Getenv("PHP_HTTPD_PATH")
		if httpdConfPath != "" {
			serverProc := NewProc("httpd", []string{"-f", httpdConfPath, "-k", "start", "-DFOREGROUND"})
			procs.Add("httpd", serverProc)
		}

		fpmConfPath := os.Getenv("PHP_FPM_PATH")
		if fpmConfPath != "" {
			phprcPath, ok := os.LookupEnv("PHPRC")
			if !ok {
				return packit.BuildResult{}, errors.New("failed to lookup $PHPRC path for FPM")
			}
			fpmProc := NewProc("php-fpm", []string{"-y", fpmConfPath, "-c", phprcPath})
			procs.Add("fpm", fpmProc)
		}

		err = procs.WriteFile(filepath.Join(layer.Path, "procs.yml"))
		if err != nil {
			panic(err)
		}
		err = fs.Copy(filepath.Join(context.CNBPath, "bin", "procmgr-binary"), filepath.Join(layer.Path, "procmgr-binary"))
		// err = copyExecutable()
		if err != nil {
			panic(err)
		}

		return packit.BuildResult{
			Layers: []packit.Layer{layer},
			Launch: packit.LaunchMetadata{
				Processes: []packit.Process{
					{
						Type:    "web",
						Command: filepath.Join(layer.Path, "procmgr-binary"),
						Args:    []string{filepath.Join(layer.Path, "procs.yml")},
						Default: true,
						Direct:  true,
					},
				},
			},
		}, nil
	}
}
