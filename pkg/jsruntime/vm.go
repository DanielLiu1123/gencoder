package jsruntime

import (
	"github.com/dop251/goja"
	"log"
)

var vm *goja.Runtime

func GetVM() *goja.Runtime {
	if vm == nil {
		vm = goja.New()

		_, err := vm.RunString(HandlebarsJS)
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
			function registerHelper(name, fn) {
				Handlebars.registerHelper(name, fn);
			};
			function registerPartial(name, template) {
				Handlebars.registerPartial(name, template);
			};
		`)
		if err != nil {
			log.Fatalf("Error initializing Handlebars.js: %v", err)
		}
	}

	return vm
}
