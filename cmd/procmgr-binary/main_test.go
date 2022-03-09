package main

import (
	"testing"

	"github.com/paketo-buildpacks/php-start/procmgr"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestProcmgr(t *testing.T) {
	spec.Run(t, "Procmgr", testProcmgr, spec.Report(report.Terminal{}))
}

func testProcmgr(t *testing.T, _ spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	it("should run a proc", func() {
		err := runProcs(procmgr.Procs{
			Processes: map[string]procmgr.Proc{
				"proc1": {
					Command: "echo",
					Args:    []string{"'Hello World!"},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
	})

	it("should fail when running a proc that doesn't exist", func() {
		err := runProcs(procmgr.Procs{
			Processes: map[string]procmgr.Proc{
				"proc1": {
					Command: "idontexist",
					Args:    []string{},
				},
			},
		})
		Expect(err).To(HaveOccurred())
	})

	it("should run two procs", func() {
		err := runProcs(procmgr.Procs{
			Processes: map[string]procmgr.Proc{
				"proc1": {
					Command: "echo",
					Args:    []string{"'Hello World!"},
				},
				"proc2": {
					Command: "echo",
					Args:    []string{"'Good-bye World!"},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
	})

	it("should fail if proc exits non-zero", func() {
		err := runProcs(procmgr.Procs{
			Processes: map[string]procmgr.Proc{
				"proc1": {
					Command: "false",
					Args:    []string{""},
				},
			},
		})
		Expect(err).To(HaveOccurred())
	})

	it("should run two procs, where one is shorter", func() {
		err := runProcs(procmgr.Procs{
			Processes: map[string]procmgr.Proc{
				"sleep0.25": {
					Command: "sleep",
					Args:    []string{"0.25"},
				},
				"sleep1": {
					Command: "sleep",
					Args:    []string{"1"},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
	})
}
