package generate

import (
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"reflect"
	"testing"
)

func Test_parseBlocks(t *testing.T) {
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
block1
gencoder block end: block1

gencoder block start: block2
block2
block2
gencoder block end: block2
out of block
`,
			},
			want: map[string]string{
				"block1": "block1\nblock1\n",
				"block2": "block2\nblock2\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBlocks(tt.args.cfg, tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseBlocks() = %v, want %v", got, tt.want)
			}
		})
	}
}
