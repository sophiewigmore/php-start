package phpstart

import (
	"github.com/paketo-buildpacks/packit/v2"
)

type BuildPlanMetadata struct {
	Launch bool
}

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {

		// // only pass detection if $BP_PHP_SERVER is set to httpd
		// server := os.Getenv("BP_PHP_SERVER")
		// if server != "httpd" {
		// 	return packit.DetectResult{}, packit.Fail
		// }

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "php",
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "php-fpm",
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "httpd",
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "httpd-conf",
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
				},
			},
		}, nil
	}
}
