package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	phpstart "github.com/paketo-buildpacks/php-start"
)

func main() {
	logEmitter := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))

	packit.Run(
		phpstart.Detect(),
		phpstart.Build(logEmitter),
	)
}
