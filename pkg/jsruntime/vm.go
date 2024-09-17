package jsruntime

import (
	"github.com/DanielLiu1123/gencoder/pkg/jsruntime/gen"
	"github.com/dop251/goja"
	"log"
	"sync"
)

var vmFunc = sync.OnceValue(func() *goja.Runtime {
	return initVM()
})

// GetVM returns the shared JS runtime
func GetVM() *goja.Runtime {
	return vmFunc()
}

func initVM() *goja.Runtime {

	vm := goja.New()

	// Load Handlebars.js
	_, err := vm.RunString(gen.HandlebarsJS)
	if err != nil {
		log.Fatalf("Error loading Handlebars.js: %v", err)
	}

	_, err = vm.RunString(`
			function compile(template) {
				return Handlebars.compile(template);
			};
			function render(template, context) {
				return template(context);
			};
			function registerPartial(name, template) {
				Handlebars.registerPartial(name, template);
			};
		`)
	if err != nil {
		log.Fatalf("Error initializing Handlebars.js: %v", err)
	}

	_, err = vm.RunString(gen.HelperJS)
	if err != nil {
		log.Fatalf("Error registering helper functions: %v", err)
	}

	return vm
}
