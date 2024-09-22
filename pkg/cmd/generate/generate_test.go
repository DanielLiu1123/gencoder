package generate

import (
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
