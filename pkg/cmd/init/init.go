package init

import (
	"github.com/DanielLiu1123/gencoder/pkg/util"
	"log"
	"os"
	"path/filepath"

	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/spf13/cobra"
)

type initOptions struct {
	output string
}

func NewCmdInit(globalOptions *model.GlobalOptions) *cobra.Command {
	opt := &initOptions{}

	c := &cobra.Command{
		Use:   "init",
		Short: "Init basic configuration for gencoder",
		Example: `  # Init basic configuration for gencoder
  $ gencoder init

  # Init basic configuration in a specific directory
  $ gencoder init -o myproject`,
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, opt, globalOptions)
		},
	}

	c.Flags().StringVarP(&opt.output, "output", "o", "", "Output directory, default to current directory")

	return c
}

func run(_ *cobra.Command, _ []string, opt *initOptions, _ *model.GlobalOptions) {
	initGencoderYaml(opt)
	initTemplates(opt)

	log.Println("Init success! Please modify the gencoder.yaml and templates to fit your project needs.")
	log.Println()
	log.Println("Thank you for using gencoder!")
}

func initGencoderYaml(opt *initOptions) {
	gencoderYaml := `templates: templates
databases:
  - dsn: 'mysql://root:root@localhost:3306/testdb'
    tables:
      - name: 'user'
        properties:
          package: 'com.example'
`
	writeFileIfNotExists(filepath.Join(opt.output, "gencoder.yaml"), []byte(gencoderYaml))
}

func initTemplates(opt *initOptions) {
	entityJava := `/**
 * @gencoder.generated: src/main/java/{{_replaceAll properties.package '.' '/'}}/{{_pascalCase table.name}}.java
 */

package {{properties.package}};

/**
 * @gencoder.block.start: table
 * <p> table: {{table.name}}
 * <p> comment: {{table.comment}}
 * <p> indexes:
     {{#each table.indexes}}
 *   <p> {{name}}: ({{#each columns}}{{name}}{{#unless @last}}, {{/unless}}{{/each}})
     {{/each}}
 */
public record {{_pascalCase table.name}} (
    {{#each table.columns}}
    /**
     * {{comment}}
     */
    {{> 'java_type.partial.hbs' columnType=type}} {{_camelCase name}}{{#unless @last}},{{/unless}}
    {{/each}}

    // NOTE: you can't make changes in the block, it will be overwritten by generating again

    // @gencoder.block.end: table
) {

    // TIP: you can make changes outside the block, it will not be overwritten by generating again
    public void hello() {
        System.out.println("Hello, World!");
    }
}
`
	javaTypePartial := `{{~#if (_match 'varchar\(\d+\)|char|tinytext|text|mediumtext|longtext' columnType)}}String
{{~else if (_match 'bigint' columnType)}}Long
{{~else if (_match 'int|integer|mediumint' columnType)}}Integer
{{~else if (_match 'smallint' columnType)}}Short
{{~else if (_match 'tinyint' columnType)}}Byte
{{~else if (_match 'bit|bool|boolean' columnType)}}Boolean
{{~else if (_match 'decimal' columnType)}}java.math.BigDecimal
{{~else if (_match 'float' columnType)}}Double
{{~else if (_match 'datetime' columnType)}}java.time.LocalDateTime
{{~else if (_match 'date' columnType)}}java.time.LocalDate
{{~else if (_match 'time' columnType)}}java.time.LocalTime
{{~else if (_match 'timestamp' columnType)}}java.time.LocalDateTime
{{~else if (_match 'varbinary' columnType)}}byte[]
{{~else if (_match 'enum.*' columnType)}}String
{{~else}}Object
{{~/if}}`

	writeFileIfNotExists(filepath.Join(opt.output, "templates", "entity.java.hbs"), []byte(entityJava))
	writeFileIfNotExists(filepath.Join(opt.output, "templates", "java_type.partial.hbs"), []byte(javaTypePartial))
}

func writeFileIfNotExists(filename string, data []byte) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		e := util.WriteFile(filename, data)
		if e != nil {
			log.Fatal(e)
		}
	}
}
