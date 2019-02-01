package example

import (
	"github.com/eosspark/eos-go/common/container/allocator"
	"log"
	"os"
	"syscall"
)


type Item struct {
	id   uint32
	name [4]byte
	//p *byte
}

const _TestMemorySize = 1024 * 1024 * 1024 * 4
const _TestMemoryFilePath = "/tmp/data/mmap.bin"

var defaultAlloc = allocator.NewDefaultAllocator(func() []byte {
	f, err := os.OpenFile(_TestMemoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	//
	if nil != err {
		log.Fatalln(err)
	}
	//
	//	// extend file
	if _, err := f.WriteAt([]byte{0}, _TestMemorySize); nil != err {
		log.Fatalln("extend error: ", err)
	}
	data, err := syscall.Mmap(int(f.Fd()), 0, _TestMemorySize, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_COPY)

	if nil != err {
		log.Fatalln(err)
	}

	if err := f.Close(); nil != err {
		log.Fatalln(err)
	}

	return data
}, nil)

//go:generate go install "github.com/eosspark/eos-go/common/container/allocator/"
//go:generate go install "github.com/eosspark/eos-go/common/container/offsetptr/"
//go:generate go install "github.com/eosspark/eos-go/common/container/list/..."
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/list" ExampleList(Item,defaultAlloc)
//go:generate go build .

