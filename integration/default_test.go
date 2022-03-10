package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testDefault(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
		source string
		name   string
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("when the buildpack is run with pack build", func() {
		var (
			image     occam.Image
			container occam.Container
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())

			source, err = occam.Source(filepath.Join("testdata", "default_app"))
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		context("HTTPD and FPM", func() {
			it("successfully starts a PHP app with HTTPD and FPM", func() {
				var (
					logs fmt.Stringer
					err  error
				)

				image, logs, err = pack.WithNoColor().Build.
					WithPullPolicy("never").
					WithBuildpacks(
						phpDistBuildpack,
						phpFpmBuildpack,
						httpdBuildpack,
						phpHttpdBuildpack,
						buildpack,
					).
					WithEnv(map[string]string{
						"BP_LOG_LEVEL":  "DEBUG",
						"BP_PHP_SERVER": "httpd",
					}).
					Execute(name, source)
				Expect(err).ToNot(HaveOccurred(), logs.String)

				Expect(logs).To(ContainLines(
					MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, buildpackInfo.Buildpack.Name)),
					"  Getting the layer associated with the server start command",
					"    /layers/paketo-buildpacks_php-start/php-start",
					"",
					"  Determining start commands to include:",
					"    HTTPD: httpd -f /layers/paketo-buildpacks_php-httpd/php-httpd-config/httpd.conf -k start -DFOREGROUND",
					"    FPM: php-fpm -y /layers/paketo-buildpacks_php-fpm/php-fpm-config/base.conf -c /layers/paketo-buildpacks_php-dist/php/etc",
					"    Writing process file to /layers/paketo-buildpacks_php-start/php-start/procs.yml",
					"",
					"  Copying procmgr-binary into /layers/paketo-buildpacks_php-start/php-start/procmgr-binary",
					"",
					"  Assigning launch processes:",
					"    web (default): /layers/paketo-buildpacks_php-start/php-start/procmgr-binary /layers/paketo-buildpacks_php-start/php-start/procs.yml",
				))

				container, err = docker.Container.Run.
					WithEnv(map[string]string{"PORT": "8080"}).
					WithPublish("8080").
					WithPublishAll().
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(container).Should(Serve(ContainSubstring("Hello World!")).OnPort(8080))
			})

		})
	})
}
