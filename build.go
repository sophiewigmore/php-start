package phpstart

import (
	"fmt"
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		// TODO: add logging
		logger.Process("START BUILDPACK")

		var serverStartCmd string
		httpdConfPath := os.Getenv("PHP_HTTPD_PATH")
		if httpdConfPath != "" {
			serverStartCmd = fmt.Sprintf("httpd -f  %s -k start -DFOREGROUND", httpdConfPath)
		}
		logger.Subprocess(serverStartCmd)

		fpmConfPath := os.Getenv("PHP_FPM_PATH")
		if fpmConfPath != "" {
			phprcPath, ok := os.LookupEnv("PHPRC")
			if !ok {
				// return packit.BuildResult{}, errors.New("failed searching for HTTPD configuration path")
				panic("no PHPRC path set")
			}
			fpmStartCmd := fmt.Sprintf("php-fpm -y %s -c %s", fpmConfPath, phprcPath)
			logger.Subprocess(fpmStartCmd)
		}

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: []packit.Process{
					{
						Type: "web",
						// Command: procmgr,
						Default: true,
					},
				},
			},
		}, nil
	}
}
