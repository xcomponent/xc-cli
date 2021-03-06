// Code generated by counterfeiter. DO NOT EDIT.
package servicesfake

import (
	"io"
	"os"
	"sync"

	"github.com/xcomponent/xc-cli/services"
)

type FakeIoService struct {
	CopyStub        func(dst io.Writer, src io.Reader) (written int64, err error)
	copyMutex       sync.RWMutex
	copyArgsForCall []struct {
		dst io.Writer
		src io.Reader
	}
	copyReturns struct {
		result1 int64
		result2 error
	}
	copyReturnsOnCall map[int]struct {
		result1 int64
		result2 error
	}
	ReadDirStub        func(dirname string) ([]os.FileInfo, error)
	readDirMutex       sync.RWMutex
	readDirArgsForCall []struct {
		dirname string
	}
	readDirReturns struct {
		result1 []os.FileInfo
		result2 error
	}
	readDirReturnsOnCall map[int]struct {
		result1 []os.FileInfo
		result2 error
	}
	ReadFileStub        func(filename string) ([]byte, error)
	readFileMutex       sync.RWMutex
	readFileArgsForCall []struct {
		filename string
	}
	readFileReturns struct {
		result1 []byte
		result2 error
	}
	readFileReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	TempDirStub        func(dir, prefix string) (name string, err error)
	tempDirMutex       sync.RWMutex
	tempDirArgsForCall []struct {
		dir    string
		prefix string
	}
	tempDirReturns struct {
		result1 string
		result2 error
	}
	tempDirReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	WriteFileStub        func(filename string, data []byte, perm os.FileMode) error
	writeFileMutex       sync.RWMutex
	writeFileArgsForCall []struct {
		filename string
		data     []byte
		perm     os.FileMode
	}
	writeFileReturns struct {
		result1 error
	}
	writeFileReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeIoService) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	fake.copyMutex.Lock()
	ret, specificReturn := fake.copyReturnsOnCall[len(fake.copyArgsForCall)]
	fake.copyArgsForCall = append(fake.copyArgsForCall, struct {
		dst io.Writer
		src io.Reader
	}{dst, src})
	fake.recordInvocation("Copy", []interface{}{dst, src})
	fake.copyMutex.Unlock()
	if fake.CopyStub != nil {
		return fake.CopyStub(dst, src)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.copyReturns.result1, fake.copyReturns.result2
}

func (fake *FakeIoService) CopyCallCount() int {
	fake.copyMutex.RLock()
	defer fake.copyMutex.RUnlock()
	return len(fake.copyArgsForCall)
}

func (fake *FakeIoService) CopyArgsForCall(i int) (io.Writer, io.Reader) {
	fake.copyMutex.RLock()
	defer fake.copyMutex.RUnlock()
	return fake.copyArgsForCall[i].dst, fake.copyArgsForCall[i].src
}

func (fake *FakeIoService) CopyReturns(result1 int64, result2 error) {
	fake.CopyStub = nil
	fake.copyReturns = struct {
		result1 int64
		result2 error
	}{result1, result2}
}

func (fake *FakeIoService) CopyReturnsOnCall(i int, result1 int64, result2 error) {
	fake.CopyStub = nil
	if fake.copyReturnsOnCall == nil {
		fake.copyReturnsOnCall = make(map[int]struct {
			result1 int64
			result2 error
		})
	}
	fake.copyReturnsOnCall[i] = struct {
		result1 int64
		result2 error
	}{result1, result2}
}

func (fake *FakeIoService) ReadDir(dirname string) ([]os.FileInfo, error) {
	fake.readDirMutex.Lock()
	ret, specificReturn := fake.readDirReturnsOnCall[len(fake.readDirArgsForCall)]
	fake.readDirArgsForCall = append(fake.readDirArgsForCall, struct {
		dirname string
	}{dirname})
	fake.recordInvocation("ReadDir", []interface{}{dirname})
	fake.readDirMutex.Unlock()
	if fake.ReadDirStub != nil {
		return fake.ReadDirStub(dirname)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.readDirReturns.result1, fake.readDirReturns.result2
}

func (fake *FakeIoService) ReadDirCallCount() int {
	fake.readDirMutex.RLock()
	defer fake.readDirMutex.RUnlock()
	return len(fake.readDirArgsForCall)
}

func (fake *FakeIoService) ReadDirArgsForCall(i int) string {
	fake.readDirMutex.RLock()
	defer fake.readDirMutex.RUnlock()
	return fake.readDirArgsForCall[i].dirname
}

func (fake *FakeIoService) ReadDirReturns(result1 []os.FileInfo, result2 error) {
	fake.ReadDirStub = nil
	fake.readDirReturns = struct {
		result1 []os.FileInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeIoService) ReadDirReturnsOnCall(i int, result1 []os.FileInfo, result2 error) {
	fake.ReadDirStub = nil
	if fake.readDirReturnsOnCall == nil {
		fake.readDirReturnsOnCall = make(map[int]struct {
			result1 []os.FileInfo
			result2 error
		})
	}
	fake.readDirReturnsOnCall[i] = struct {
		result1 []os.FileInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeIoService) ReadFile(filename string) ([]byte, error) {
	fake.readFileMutex.Lock()
	ret, specificReturn := fake.readFileReturnsOnCall[len(fake.readFileArgsForCall)]
	fake.readFileArgsForCall = append(fake.readFileArgsForCall, struct {
		filename string
	}{filename})
	fake.recordInvocation("ReadFile", []interface{}{filename})
	fake.readFileMutex.Unlock()
	if fake.ReadFileStub != nil {
		return fake.ReadFileStub(filename)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.readFileReturns.result1, fake.readFileReturns.result2
}

func (fake *FakeIoService) ReadFileCallCount() int {
	fake.readFileMutex.RLock()
	defer fake.readFileMutex.RUnlock()
	return len(fake.readFileArgsForCall)
}

func (fake *FakeIoService) ReadFileArgsForCall(i int) string {
	fake.readFileMutex.RLock()
	defer fake.readFileMutex.RUnlock()
	return fake.readFileArgsForCall[i].filename
}

func (fake *FakeIoService) ReadFileReturns(result1 []byte, result2 error) {
	fake.ReadFileStub = nil
	fake.readFileReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeIoService) ReadFileReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.ReadFileStub = nil
	if fake.readFileReturnsOnCall == nil {
		fake.readFileReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.readFileReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeIoService) TempDir(dir string, prefix string) (name string, err error) {
	fake.tempDirMutex.Lock()
	ret, specificReturn := fake.tempDirReturnsOnCall[len(fake.tempDirArgsForCall)]
	fake.tempDirArgsForCall = append(fake.tempDirArgsForCall, struct {
		dir    string
		prefix string
	}{dir, prefix})
	fake.recordInvocation("TempDir", []interface{}{dir, prefix})
	fake.tempDirMutex.Unlock()
	if fake.TempDirStub != nil {
		return fake.TempDirStub(dir, prefix)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.tempDirReturns.result1, fake.tempDirReturns.result2
}

func (fake *FakeIoService) TempDirCallCount() int {
	fake.tempDirMutex.RLock()
	defer fake.tempDirMutex.RUnlock()
	return len(fake.tempDirArgsForCall)
}

func (fake *FakeIoService) TempDirArgsForCall(i int) (string, string) {
	fake.tempDirMutex.RLock()
	defer fake.tempDirMutex.RUnlock()
	return fake.tempDirArgsForCall[i].dir, fake.tempDirArgsForCall[i].prefix
}

func (fake *FakeIoService) TempDirReturns(result1 string, result2 error) {
	fake.TempDirStub = nil
	fake.tempDirReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeIoService) TempDirReturnsOnCall(i int, result1 string, result2 error) {
	fake.TempDirStub = nil
	if fake.tempDirReturnsOnCall == nil {
		fake.tempDirReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.tempDirReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeIoService) WriteFile(filename string, data []byte, perm os.FileMode) error {
	var dataCopy []byte
	if data != nil {
		dataCopy = make([]byte, len(data))
		copy(dataCopy, data)
	}
	fake.writeFileMutex.Lock()
	ret, specificReturn := fake.writeFileReturnsOnCall[len(fake.writeFileArgsForCall)]
	fake.writeFileArgsForCall = append(fake.writeFileArgsForCall, struct {
		filename string
		data     []byte
		perm     os.FileMode
	}{filename, dataCopy, perm})
	fake.recordInvocation("WriteFile", []interface{}{filename, dataCopy, perm})
	fake.writeFileMutex.Unlock()
	if fake.WriteFileStub != nil {
		return fake.WriteFileStub(filename, data, perm)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.writeFileReturns.result1
}

func (fake *FakeIoService) WriteFileCallCount() int {
	fake.writeFileMutex.RLock()
	defer fake.writeFileMutex.RUnlock()
	return len(fake.writeFileArgsForCall)
}

func (fake *FakeIoService) WriteFileArgsForCall(i int) (string, []byte, os.FileMode) {
	fake.writeFileMutex.RLock()
	defer fake.writeFileMutex.RUnlock()
	return fake.writeFileArgsForCall[i].filename, fake.writeFileArgsForCall[i].data, fake.writeFileArgsForCall[i].perm
}

func (fake *FakeIoService) WriteFileReturns(result1 error) {
	fake.WriteFileStub = nil
	fake.writeFileReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeIoService) WriteFileReturnsOnCall(i int, result1 error) {
	fake.WriteFileStub = nil
	if fake.writeFileReturnsOnCall == nil {
		fake.writeFileReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.writeFileReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeIoService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.copyMutex.RLock()
	defer fake.copyMutex.RUnlock()
	fake.readDirMutex.RLock()
	defer fake.readDirMutex.RUnlock()
	fake.readFileMutex.RLock()
	defer fake.readFileMutex.RUnlock()
	fake.tempDirMutex.RLock()
	defer fake.tempDirMutex.RUnlock()
	fake.writeFileMutex.RLock()
	defer fake.writeFileMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeIoService) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ services.IoService = new(FakeIoService)
