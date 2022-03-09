package fakes

import (
	"sync"

	"github.com/paketo-buildpacks/php-start/procmgr"
)

type ProcMgr struct {
	AddCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Name string
			Proc procmgr.Proc
		}
		Stub func(string, procmgr.Proc)
	}
	AppendOrUpdateProcsCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Name string
			Proc procmgr.Proc
		}
		Returns struct {
			Error error
		}
		Stub func(string, procmgr.Proc) error
	}
	WriteProcsCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Path string
		}
		Returns struct {
			Error error
		}
		Stub func(string) error
	}
}

func (f *ProcMgr) Add(param1 string, param2 procmgr.Proc) {
	f.AddCall.mutex.Lock()
	defer f.AddCall.mutex.Unlock()
	f.AddCall.CallCount++
	f.AddCall.Receives.Name = param1
	f.AddCall.Receives.Proc = param2
	if f.AddCall.Stub != nil {
		f.AddCall.Stub(param1, param2)
	}
}
func (f *ProcMgr) AppendOrUpdateProcs(param1 string, param2 procmgr.Proc) error {
	f.AppendOrUpdateProcsCall.mutex.Lock()
	defer f.AppendOrUpdateProcsCall.mutex.Unlock()
	f.AppendOrUpdateProcsCall.CallCount++
	f.AppendOrUpdateProcsCall.Receives.Name = param1
	f.AppendOrUpdateProcsCall.Receives.Proc = param2
	if f.AppendOrUpdateProcsCall.Stub != nil {
		return f.AppendOrUpdateProcsCall.Stub(param1, param2)
	}
	return f.AppendOrUpdateProcsCall.Returns.Error
}
func (f *ProcMgr) WriteProcs(param1 string) error {
	f.WriteProcsCall.mutex.Lock()
	defer f.WriteProcsCall.mutex.Unlock()
	f.WriteProcsCall.CallCount++
	f.WriteProcsCall.Receives.Path = param1
	if f.WriteProcsCall.Stub != nil {
		return f.WriteProcsCall.Stub(param1)
	}
	return f.WriteProcsCall.Returns.Error
}
