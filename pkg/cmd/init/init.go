package init

import (
	"log"
	"os"
	"path/filepath"

	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/spf13/cobra"
)

type initOptions struct {
}

func NewCmdInit(globalOptions *model.GlobalOptions) *cobra.Command {
	opt := &initOptions{}

	return &cobra.Command{
		Use:   "init",
		Short: "Init basic configuration for gencoder",
		Example: `  # Init basic configuration for gencoder
  $ gencoder init`,
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, opt, globalOptions)
		},
	}
}

func run(_ *cobra.Command, _ []string, _ *initOptions, _ *model.GlobalOptions) {
	initGencoderYaml()
	initTemplatesDir()
	initTemplates()

	log.Println("Init success! Please modify the gencoder.yaml and templates to fit your project needs.")
	log.Println()
	log.Println("Thank you for using gencoder!")
}

func initGencoderYaml() {
	gencoderYaml := `templates: templates
databases:
  - dsn: 'mysql://root:root@localhost:3306/testdb'
    tables:
      - name: 'user'
        properties:
          package: 'com.example'
`
	writeFileIfNotExists("gencoder.yaml", []byte(gencoderYaml))
}

func initTemplatesDir() {
	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		err := os.Mkdir("templates", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func initTemplates() {
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
	writeFileIfNotExists(filepath.Join("templates", "entity.java.hbs"), []byte(entityJava))

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
	writeFileIfNotExists(filepath.Join("templates", "java_type.partial.hbs"), []byte(javaTypePartial))
}

func writeFileIfNotExists(filename string, data []byte) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := os.WriteFile(filename, data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
