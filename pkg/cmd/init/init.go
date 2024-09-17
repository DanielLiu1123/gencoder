package init

import (
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/spf13/cobra"
	"log"
	"os"
)

type initOptions struct {
}

func NewCmdInit(globalOptions *model.GlobalOptions) *cobra.Command {

	opt := &initOptions{}

	c := &cobra.Command{
		Use:   "init",
		Short: "Init basic configuration for gencoder",
		PreRun: func(cmd *cobra.Command, args []string) {
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, opt, globalOptions)
		},
	}

	return c
}

func run(_ *cobra.Command, _ []string, _ *initOptions, _ *model.GlobalOptions) {

	// init gencoder.yaml
	gencoderYaml := `templatesDir: templates
databases:
  - name: mysql
    dsn: 'mysql://root:root@localhost:3306/testdb'
    tables:
      - name: 'user'
        properties:
          package: mysqlpackage
        ignoreColumns:
          - deleted_at
  - name: postgres
    dsn: 'postgres://root:root@localhost:5432/testdb?sslmode=disable'
    properties:
      package: defaultpackage
    tables:
      - name: 'user'
        properties:
          package: postgrespackage
        ignoreColumns:
          - deleted_at`

	if _, err := os.Stat("gencoder.yaml"); err != nil {
		err := os.WriteFile("gencoder.yaml", []byte(gencoderYaml), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// init templates dir
	if _, err := os.Stat("templates"); err != nil {
		err = os.Mkdir("templates", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// init templates
	entityJava := `/**
 * @gencoder.generated: src/main/java/{{replaceAll properties.package '\.' '/'}}/{{pascalCase table.name}}.java
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
public record {{pascalCase table.name}} (
    {{#each table.columns}}
    /**
     * {{comment}}
     */
    {{> 'java_type.partial.hbs' columnType=type}} {{camelCase name}}{{#unless @last}},{{/unless}}
    {{/each}}
    // @gencoder.block.end: table
) {}
`

	if _, err := os.Stat("templates/entity.java.hbs"); err != nil {
		err = os.WriteFile("templates/entity.java.hbs", []byte(entityJava), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	javaTypePartial := `{{~#if (match 'varchar\(\d+\)|char|tinytext|text|mediumtext|longtext' columnType)}}String
{{~else if (match 'bigint' columnType)}}Long
{{~else if (match 'int|integer|mediumint' columnType)}}Integer
{{~else if (match 'smallint' columnType)}}Short
{{~else if (match 'tinyint' columnType)}}Byte
{{~else if (match 'bit|bool|boolean' columnType)}}Boolean
{{~else if (match 'decimal' columnType)}}java.math.BigDecimal
{{~else if (match 'float' columnType)}}Double
{{~else if (match 'datetime' columnType)}}java.time.LocalDateTime
{{~else if (match 'date' columnType)}}java.time.LocalDate
{{~else if (match 'time' columnType)}}java.time.LocalTime
{{~else if (match 'timestamp' columnType)}}java.time.LocalDateTime
{{~else if (match 'varbinary' columnType)}}byte[]
{{~else if (match 'enum.*' columnType)}}String
{{~else}}Object
{{~/if}}`

	if _, err := os.Stat("templates/java_type.partial.hbs"); err != nil {
		err = os.WriteFile("templates/java_type.partial.hbs", []byte(javaTypePartial), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Init success! Please modify the gencoder.yaml and templates to fit your project needs.")
	log.Println()
	log.Println("Thank you for using gencoder!")
}
