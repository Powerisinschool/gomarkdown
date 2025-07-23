package parser

import (
	"reflect"
	"testing"

	"github.com/Powerisinschool/gomarkdown/internal/ast"
)

func TestInlineParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []ast.InlineNode
	}{
		{
			name:  "Simple text",
			input: "hello world",
			want:  []ast.InlineNode{&ast.Text{Value: "hello world"}},
		},
		{
			name:  "Emphasis",
			input: "*hello*",
			want:  []ast.InlineNode{&ast.Emphasis{Children: []ast.InlineNode{&ast.Text{Value: "hello"}}}},
		},
		{
			name:  "Strong",
			input: "**hello**",
			want:  []ast.InlineNode{&ast.Strong{Children: []ast.InlineNode{&ast.Text{Value: "hello"}}}},
		},
		{
			name:  "Code",
			input: "`hello`",
			want:  []ast.InlineNode{&ast.Code{Value: "hello"}},
		},
		{
			name:  "Mixed emphasis and text",
			input: "hello *world*",
			want: []ast.InlineNode{
				&ast.Text{Value: "hello "},
				&ast.Emphasis{Children: []ast.InlineNode{&ast.Text{Value: "world"}}},
			},
		},
		{
			name:  "Unmatched asterisk",
			input: "hello *world",
			want:  []ast.InlineNode{&ast.Text{Value: "hello *world"}},
		},
		{
			name:  "Nested strong and emphasis",
			input: "**hello *world*!**",
			want: []ast.InlineNode{
				&ast.Strong{Children: []ast.InlineNode{
					&ast.Text{Value: "hello "},
					&ast.Emphasis{Children: []ast.InlineNode{&ast.Text{Value: "world"}}},
					&ast.Text{Value: "!"},
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got := p.inlineParse(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inlineParse() = %v, want %v", got, tt.want)
			}
		})
	}
}
