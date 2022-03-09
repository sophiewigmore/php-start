package phpstart_test

import (
	"bytes"
	"os"
	"path/filepath"
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

		Expect(os.Mkdir(filepath.Join(cnbDir, "bin"), 0700)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(cnbDir, "bin", "procmgr-binary"), []byte{}, 0644)).To(Succeed())

		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		buffer = bytes.NewBuffer(nil)
		logEmitter := scribe.NewEmitter(buffer)

		procMgr = &fakes.ProcMgr{}
		build = phpstart.Build(procMgr, logEmitter)
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
		Expect(os.RemoveAll(cnbDir)).To(Succeed())
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("the PHP_HTTPD_PATH env var is set", func() {
		it.Before(func() {
			Expect(os.Setenv("PHP_HTTPD_PATH", "httpd-conf-path")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("PHP_HTTPD_PATH")).To(Succeed())
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
				Type:    "web",
				Command: filepath.Join(layersDir, "php-start", "procmgr-binary"),
				Args: []string{
					filepath.Join(layersDir, "php-start", "procs.yml"),
				},
				Default: true,
				Direct:  true,
			}))
			Expect(result.Layers[0].Name).To(Equal("php-start"))
			Expect(result.Layers[0].Path).To(Equal(filepath.Join(layersDir, "php-start")))
			Expect(result.Layers[0].Launch).To(BeTrue())
			Expect(result.Layers[0].Build).To(BeFalse())

			Expect(procMgr.AddCall.CallCount).To(Equal(1))
			Expect(procMgr.AddCall.Receives.Name).To(Equal("httpd"))
			Expect(procMgr.AddCall.Receives.Proc.Command).To(Equal("httpd"))
			Expect(procMgr.AddCall.Receives.Proc.Args).To(Equal([]string{
				"-f",
				"httpd-conf-path",
				"-k",
				"start",
				"-DFOREGROUND",
			}))

			Expect(procMgr.WriteFileCall.Receives.Path).To(Equal(filepath.Join(layersDir, "php-start", "procs.yml")))
			// Assert logs show both things being run
			// Add some REadProcs function and make sure it has:
		})
	})

	context("the PHP_FPM_PATH and PHPRC env vars are set", func() {
		it.Before(func() {
			Expect(os.Setenv("PHP_FPM_PATH", "fpm-conf-path")).To(Succeed())
			Expect(os.Setenv("PHPRC", "phprc-path")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("PHP_FPM_PATH")).To(Succeed())
			Expect(os.Unsetenv("PHPRC")).To(Succeed())
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
				Type:    "web",
				Command: filepath.Join(layersDir, "php-start", "procmgr-binary"),
				Args: []string{
					filepath.Join(layersDir, "php-start", "procs.yml"),
				},
				Default: true,
				Direct:  true,
			}))

			Expect(procMgr.AddCall.CallCount).To(Equal(1))

			Expect(procMgr.AddCall.Receives.Name).To(Equal("fpm"))
			Expect(procMgr.AddCall.Receives.Proc.Command).To(Equal("php-fpm"))
			Expect(procMgr.AddCall.Receives.Proc.Args).To(Equal([]string{
				"-y",
				"fpm-conf-path",
				"-c",
				"phprc-path",
			}))

			Expect(procMgr.WriteFileCall.Receives.Path).To(Equal(filepath.Join(layersDir, "php-start", "procs.yml")))
		})
	})

	context("failure cases", func() {
		context("the PHP_FPM_PATH is set but PHPRC is not", func() {
			it.Before(func() {
				Expect(os.Setenv("PHP_FPM_PATH", "fpm-conf-path")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("PHP_FPM_PATH")).To(Succeed())
			})

			it("returns an error, since the PHPRC is needed for FPM start command", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					CNBPath:    cnbDir,
					BuildpackInfo: packit.BuildpackInfo{
						Name:    "Some Buildpack",
						Version: "some-version",
					},
					Layers: packit.Layers{Path: layersDir},
				})
				Expect(err).To(MatchError(ContainSubstring("failed to lookup $PHPRC path for FPM")))
			})
		})

	})
}
