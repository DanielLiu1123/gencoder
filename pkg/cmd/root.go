package cmd

import (
	"github.com/DanielLiu1123/gencoder/pkg/cmd/generate"
	"github.com/DanielLiu1123/gencoder/pkg/cmd/introspect"
	"github.com/DanielLiu1123/gencoder/pkg/jsruntime"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/spf13/cobra"
	"log"
)

func NewCmdRoot() *cobra.Command {

	registerHelperFunctions()

	opt := &model.GlobalOptions{}

	c := &cobra.Command{
		Use:   "gencoder <command> [flags]",
		Short: "gencoder short",
		Long:  `gencoder longlonglong`,
	}

	c.Flags().StringVarP(&opt.Config, "config", "f", "gencoder.yaml", "Config file to use")

	c.AddCommand(generate.NewCmdGenerate(opt))
	c.AddCommand(introspect.NewCmdIntrospect(opt))

	return c
}

func registerHelperFunctions() {

	vm := jsruntime.GetVM()

	_, err := vm.RunString(`
		Handlebars.registerHelper('replaceAll', function(target, old, newV) {
		    return target.replace(new RegExp(old, 'g'), newV);
		});

		Handlebars.registerHelper('match', function(pattern, target) {
			return new RegExp(pattern).test(target);
		});

		Handlebars.registerHelper('ne', function(left, right) {
		    return left !== right;
		});
		
		Handlebars.registerHelper('snakeCase', function(s) {
		    return s.replace(/([a-z])([A-Z])/g, '$1_$2').toLowerCase();
		});
		
		Handlebars.registerHelper('camelCase', function(s) {
			if (!s) {
				return s;
			}
		    return s.replace(/([-_][a-z])/ig, function($1) {
		        return $1.toUpperCase()
		            .replace('-', '')
		            .replace('_', '');
		    });
		});
		
		Handlebars.registerHelper('pascalCase', function(s) {
		    return s.replace(/(\w)(\w*)/g, function($0, $1, $2) {
		        return $1.toUpperCase() + $2.toLowerCase();
		    });
		});
		
		Handlebars.registerHelper('upperFirst', function(s) {
		    return s.charAt(0).toUpperCase() + s.slice(1);
		});
		
		Handlebars.registerHelper('lowerFirst', function(s) {
		    return s.charAt(0).toLowerCase() + s.slice(1);
		});
		
		Handlebars.registerHelper('uppercase', function(s) {
		    return s.toUpperCase();
		});
		
		Handlebars.registerHelper('lowercase', function(s) {
		    return s.toLowerCase();
		});
		
		Handlebars.registerHelper('trim', function(s) {
		    return s.trim();
		});
		
		Handlebars.registerHelper('removePrefix', function(s, prefix) {
		    return s.startsWith(prefix) ? s.slice(prefix.length) : s;
		});
		
		Handlebars.registerHelper('removeSuffix', function(s, suffix) {
		    return s.endsWith(suffix) ? s.slice(0, -suffix.length) : s;
		});
	`)
	if err != nil {
		log.Fatalf("Error registering helper functions: %v", err)
	}
}
