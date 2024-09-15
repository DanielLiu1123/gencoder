package handlebars

import (
	"github.com/DanielLiu1123/gencoder/pkg/jsruntime"
	"github.com/dop251/goja"
	"log"
)

// Compile compiles a Handlebars template
func Compile(template string) goja.Value {
	vm := jsruntime.GetVM()

	compileFunc, ok := goja.AssertFunction(vm.Get("compile"))
	if !ok {
		log.Fatal("Error getting 'compile' function")
	}

	compiledTemplate, err := compileFunc(goja.Undefined(), vm.ToValue(template))
	if err != nil {
		log.Fatalf("Error compiling template: %v", err)
	}

	return compiledTemplate
}

// Render renders a Handlebars template with the given context
func Render(template goja.Value, context map[string]interface{}) string {
	vm := jsruntime.GetVM()

	renderFunc, ok := goja.AssertFunction(vm.Get("render"))
	if !ok {
		log.Fatal("Error getting 'render' function")
	}

	result, err := renderFunc(goja.Undefined(), template, vm.ToValue(context))
	if err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}

	return result.String()
}

// RegisterPartial registers a Handlebars partial
func RegisterPartial(name string, template string) {
	vm := jsruntime.GetVM()

	registerPartialFunc, ok := goja.AssertFunction(vm.Get("registerPartial"))
	if !ok {
		log.Fatal("Error getting 'registerPartial' function")
	}

	_, err := registerPartialFunc(goja.Undefined(), vm.ToValue(name), vm.ToValue(template))
	if err != nil {
		log.Fatalf("Error registering partial: %v", err)
	}
}
