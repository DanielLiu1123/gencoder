package jsruntime

import "github.com/dop251/goja"

var vm *goja.Runtime

func GetVM() *goja.Runtime {
	if vm == nil {
		vm = goja.New()
	}
	return vm
}
