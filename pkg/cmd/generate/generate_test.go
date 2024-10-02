package generate

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"

	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestParseBlocks(t *testing.T) {
	type args struct {
		cfg     *model.Config
		content string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Test buildBlocks",
			args: args{
				cfg: &model.Config{
					BlockMarker: model.BlockMarker{
						Start: "gencoder block start:",
						End:   "gencoder block end:",
					},
				},
				content: `
out of block
gencoder block start: block1
block1
gencoder block end: block1

gencoder block start: block2
block2
gencoder block end: block2
out of block
`,
			},
			want: map[string]string{
				"block1": `gencoder block start: block1
block1
gencoder block end: block1`,
				"block2": `gencoder block start: block2
block2
gencoder block end: block2`,
			},
		},
		{
			name: "Test buildBlocks with no end marker",
			args: args{
				cfg: &model.Config{},
				content: `
out of block

@gencoder.block.start: block1
block1
@gencoder.block.end: block1

@gencoder.block.start: block2
block2
block2

`,
			},
			want: map[string]string{
				"block1": `@gencoder.block.start: block1
block1
@gencoder.block.end: block1`,
				"block2": `@gencoder.block.start: block2
block2
block2

`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseBlocks(tt.args.cfg, tt.args.content)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReplaceBlocks(t *testing.T) {
	tests := []struct {
		name       string
		cfg        *model.Config
		oldContent string
		newContent string
		wantResult string
	}{
		{
			name: "Replace existing blocks",
			cfg: &model.Config{
				BlockMarker: model.BlockMarker{
					Start: "gencoder block start:",
					End:   "gencoder block end:",
				},
			},
			oldContent: `
out of block
gencoder block start: block1
old content 1
gencoder block end: block1

gencoder block start: block2
old content 2
gencoder block end: block2
out of block
`,
			newContent: `
gencoder block start: block1
new content 1
gencoder block end: block1

gencoder block start: block2
new content 2
gencoder block end: block2
`,
			wantResult: `
out of block
gencoder block start: block1
new content 1
gencoder block end: block1

gencoder block start: block2
new content 2
gencoder block end: block2
out of block
`,
		},
		{
			name: "Keep non-existing blocks",
			cfg: &model.Config{
				BlockMarker: model.BlockMarker{
					Start: "gencoder block start:",
					End:   "gencoder block end:",
				},
			},
			oldContent: `
out of block
gencoder block start: block1
old content 1
gencoder block end: block1

gencoder block start: block2
old content 2
gencoder block end: block2
out of block
`,
			newContent: `
gencoder block start: block1
new content 1
gencoder block end: block1

gencoder block start: block3
new content 3
gencoder block end: block3
`,
			wantResult: `
out of block
gencoder block start: block1
new content 1
gencoder block end: block1

gencoder block start: block2
old content 2
gencoder block end: block2
out of block
`,
		},
		{
			name: "No end marker",
			cfg: &model.Config{
				BlockMarker: model.BlockMarker{
					Start: "gencoder block start:",
					End:   "gencoder block end:",
				},
			},
			oldContent: `
out of block
gencoder block start: block1
old content 1
gencoder block end: block1

gencoder block start: block2
old content 2
old content 2
`,
			newContent: `
gencoder block start: block1
new content 1
gencoder block end: block1

gencoder block start: block2
new content 2
`,
			wantResult: `
out of block
gencoder block start: block1
new content 1
gencoder block end: block1

gencoder block start: block2
new content 2
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := replaceBlocks(tt.cfg, tt.oldContent, tt.newContent)
			assert.Equal(t, tt.wantResult, gotResult)
		})
	}
}

func TestNewCmdGenerate_whenConfigIsSet_thenShouldUsingSpecificConfigFile(t *testing.T) {
	workDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)
	defer os.RemoveAll(workDir)

	_ = os.Chdir(workDir)

	createNewFile(filepath.Join(workDir, "config/gencoder.yaml"), []byte("templates: tpl"))
	createNewFile(filepath.Join(workDir, "tpl/test1.text.hbs"), []byte(`@gencoder.generated: test1.txt
Hello, {{properties.name}}!`))

	cmd := NewCmdGenerate(&model.GlobalOptions{})
	cmd.SetArgs([]string{"--config", "config/gencoder.yaml", "--properties", "name=World"})

	err = cmd.Execute()
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(workDir, "test1.txt"))
	assert.NoError(t, err)
	assert.Equal(t, `@gencoder.generated: test1.txt
Hello, World!`, string(content))
}

func TestNewCmdGenerate_whenImportHelperIsSet_thenShouldRegisterCustomHelpers(t *testing.T) {
	workDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)
	defer os.RemoveAll(workDir)

	_ = os.Chdir(workDir)

	createNewFile(filepath.Join(workDir, "gencoder.yaml"), []byte(`
templates: templates
`))
	createNewFile(filepath.Join(workDir, "templates/test1.text.hbs"), []byte(`@gencoder.generated: test1.txt
Hello, {{_toUpperCase properties.name}}!`))
	createNewFile(filepath.Join(workDir, "helpers.js"), []byte(`
Handlebars.registerHelper('_toUpperCase', function (target) {
	return target.toUpperCase();
});
`))

	cmd := NewCmdGenerate(&model.GlobalOptions{})
	cmd.SetArgs([]string{"--config", "gencoder.yaml", "--properties", "name=World", "--import-helper", "helpers.js"})

	err = cmd.Execute()
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(workDir, "test1.txt"))
	assert.NoError(t, err)
	assert.Equal(t, `@gencoder.generated: test1.txt
Hello, WORLD!`, string(content))
}

func TestNewCmdGenerate_whenUsingIncludeNonTpl_thenShouldGenerateNonTemplateFiles(t *testing.T) {

	// Create template directory
	tplDir, err := os.MkdirTemp("", "tpl")
	require.NoError(t, err)
	defer os.RemoveAll(tplDir)

	_ = os.Chdir(tplDir)

	// Create non-template files
	createNewFile(filepath.Join(tplDir, "templates/non-template1.txt"), []byte("This is a non-template file"))
	createNewFile(filepath.Join(tplDir, "templates/foo/non-template2.txt"), []byte("This is a non-template file"))
	// Create partial file
	createNewFile(filepath.Join(tplDir, "templates/header.txt.hbs"), []byte("This is a header"))
	// Create template file
	createNewFile(filepath.Join(tplDir, "templates/test1.text.hbs"), []byte(`@gencoder.generated: test1.txt
{{> header.txt.hbs}}

Hello, {{properties.name}}!`))
	createNewFile(filepath.Join(tplDir, "templates/test2.text.hbs"), []byte(`@gencoder.generated: foo/test2.txt
{{> header.txt.hbs}}

Hello, {{properties.name}}!`))

	// Create generated directory
	generatedDir, err := os.MkdirTemp("", "generated")
	require.NoError(t, err)
	defer os.RemoveAll(generatedDir)

	_ = os.Chdir(generatedDir)

	cmd := NewCmdGenerate(&model.GlobalOptions{})

	cmd.SetArgs([]string{"--templates", filepath.Join(tplDir, "templates"), "--include-non-tpl", "--properties", "name=World"})

	err = cmd.Execute()
	require.NoError(t, err)

	// Verify non-template files are generated
	_, err = os.Stat(filepath.Join(generatedDir, "non-template1.txt"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(generatedDir, "foo/non-template2.txt")) // keep directory structure
	assert.NoError(t, err)

	// Verify partial files do not exist
	_, err = os.Stat(filepath.Join(generatedDir, "header.txt.hbs"))
	assert.Error(t, err)

	// Verify template files are generated
	content, err := os.ReadFile(filepath.Join(generatedDir, "test1.txt"))
	assert.NoError(t, err)
	assert.Equal(t, `@gencoder.generated: test1.txt
This is a header
Hello, World!`, string(content))

	content, err = os.ReadFile(filepath.Join(generatedDir, "foo/test2.txt")) // keep directory structure
	assert.NoError(t, err)
	assert.Equal(t, `@gencoder.generated: foo/test2.txt
This is a header
Hello, World!`, string(content))
}

func TestNewCmdGenerate_whenTemplatesIsSet_thenShouldOverrideTemplatesInConfigFile(t *testing.T) {
	workDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)
	defer os.RemoveAll(workDir)

	_ = os.Chdir(workDir)

	createNewFile(filepath.Join(workDir, "gencoder.yaml"), []byte(`
templates: tpl
`))
	createNewFile(filepath.Join(workDir, "tpl/test1.text.hbs"), []byte(`@gencoder.generated: test1.txt
Hello, {{properties.name}}! -- from tpl`))
	createNewFile(filepath.Join(workDir, "templates/test1.text.hbs"), []byte(`@gencoder.generated: test1.txt
Hello, {{properties.name}}! -- from templates`))

	cmd := NewCmdGenerate(&model.GlobalOptions{})
	cmd.SetArgs([]string{"--config", "gencoder.yaml", "--templates", "templates", "--properties", "name=World"})

	err = cmd.Execute()
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(workDir, "test1.txt"))
	assert.NoError(t, err)
	assert.Equal(t, `@gencoder.generated: test1.txt
Hello, World! -- from templates`, string(content))
}

func TestNewCmdGenerate_whenPropertiesIsSet_thenShouldOverridePropertiesInConfigFile(t *testing.T) {
	workDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)
	defer os.RemoveAll(workDir)

	_ = os.Chdir(workDir)

	createNewFile(filepath.Join(workDir, "gencoder.yaml"), []byte(`
templates: templates
properties:
  name: John
  age: 20
  height: 180
`))
	createNewFile(filepath.Join(workDir, "templates/test1.text.hbs"), []byte(`@gencoder.generated: test1.txt
Hello, I'm {{properties.name}}, {{properties.age}} years old, and {{properties.height}}cm tall!`))

	cmd := NewCmdGenerate(&model.GlobalOptions{})
	cmd.SetArgs([]string{"--config", "gencoder.yaml", "--properties", "name=Daniel,height=170"})

	err = cmd.Execute()
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(workDir, "test1.txt"))
	assert.NoError(t, err)
	assert.Equal(t, `@gencoder.generated: test1.txt
Hello, I'm Daniel, 20 years old, and 170cm tall!`, string(content))
}

func TestNewCmdGenerate_whenOutputIsSet_thenShouldOverrideOutputInConfigFile(t *testing.T) {
	workDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)
	defer os.RemoveAll(workDir)

	_ = os.Chdir(workDir)

	createNewFile(filepath.Join(workDir, "gencoder.yaml"), []byte(`
templates: templates
output: output
`))
	createNewFile(filepath.Join(workDir, "templates/test1.text.hbs"), []byte(`@gencoder.generated: test1.txt
Hello, {{properties.name}}!`))

	cmd := NewCmdGenerate(&model.GlobalOptions{})
	cmd.SetArgs([]string{"--config", "gencoder.yaml", "--output", "output2", "--properties", "name=World"})

	err = cmd.Execute()
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(workDir, "output2/test1.txt"))
	assert.NoError(t, err)
	assert.Equal(t, `@gencoder.generated: test1.txt
Hello, World!`, string(content))
}
