package phpstart_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	phpstart "github.com/paketo-buildpacks/php-start"
	"github.com/paketo-buildpacks/php-start/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		layersDir  string
		workingDir string
		cnbDir     string

		buffer  *bytes.Buffer
		procMgr *fakes.ProcMgr

		build packit.BuildFunc
	)

	it.Before(func() {
		var err error
		layersDir, err = os.MkdirTemp("", "layers")
		Expect(err).NotTo(HaveOccurred())

		cnbDir, err = os.MkdirTemp("", "cnb")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		buffer = bytes.NewBuffer(nil)
		logEmitter := scribe.NewEmitter(buffer)

		procMgr = &fakes.ProcMgr{}
		// Expect(os.Setenv("PHPRC", "some-php-dist-path")).To(Succeed())
		build = phpstart.Build(procMgr, logEmitter)
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
		Expect(os.RemoveAll(cnbDir)).To(Succeed())
		Expect(os.RemoveAll(workingDir)).To(Succeed())
		// Expect(os.Unsetenv("PHPRC")).To(Succeed())
	})

	it("returns a result that sets a PHP web server and FPM start command", func() {
		result, err := build(packit.BuildContext{
			WorkingDir: workingDir,
			CNBPath:    cnbDir,
			Stack:      "some-stack",
			BuildpackInfo: packit.BuildpackInfo{
				Name:    "Some Buildpack",
				Version: "some-version",
			},
			Plan: packit.BuildpackPlan{
				Entries: []packit.BuildpackPlanEntry{},
			},
			Layers: packit.Layers{Path: layersDir},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Launch.Processes[0]).To(Equal(packit.Process{
			Type: "web",
			// Command: "progmgr",
			// Args:    "some-procfile",
			Default: true,
			Direct:  true,
		}))

		// Assert logs show both things being run
		// Add some REadProcs function and make sure it has:

		// "php-fpm -y $PHP_FPM_PATH -c $PHPRC"
		// "httpd -f $PHP_HTTPD_PATH -k start -DFOREGROUND"
		// procs, err := procmgr.ReadProcs(procFile)
		// Expect(err).ToNot(HaveOccurred())
		// Expect(procs.Processes).To(ContainElement(phpFpmProc))
		// Expect(procs.Processes).To(ContainElement(httpdProc))
	})

	context("failure cases", func() {})
}
