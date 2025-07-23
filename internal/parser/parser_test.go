package parser_test

import (
	"github.com/Powerisinschool/gomarkdown/internal/ast"
	"github.com/Powerisinschool/gomarkdown/internal/parser"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("Parse", func(t *testing.T) {
		test_input := "Hello, world!"
		p := parser.New([]byte(test_input))
		doc := p.Parse()
		if doc.Children[0].(*ast.Paragraph).Children[0].(*ast.Text).Value != test_input {
			t.Errorf("expected '%s', got '%s'", test_input, doc.Children[0].(*ast.Paragraph).Children[0].(*ast.Text).Value)
		}
	})
}

func TestParseEmpty(t *testing.T) {
	t.Run("ParseEmpty", func(t *testing.T) {
		test_input := ""
		p := parser.New([]byte(test_input))
		doc := p.Parse()
		if doc.Children != nil {
			t.Errorf("expected nil, got '%s'", doc.Children[0].(*ast.Paragraph).Children[0].(*ast.Text).Value)
		}
	})
}

func TestParseHeading(t *testing.T) {
	t.Run("ParseHeading", func(t *testing.T) {
		test_input := "# Hello, world!"
		p := parser.New([]byte(test_input))
		doc := p.Parse()
		if doc.Children[0].(*ast.Heading).Children[0].(*ast.Text).Value != test_input[2:] {
			t.Errorf("expected '%s', got '%s'", test_input[2:], doc.Children[0].(*ast.Heading).Children[0].(*ast.Text).Value)
		}
	})
}

//func TestParseStrong(t *testing.T) {
//	t.Run("ParseStrong", func(t *testing.T) {
//		test_input := "**Hello, world!**"
//		p := parser.New([]byte(test_input))
//		doc := p.Parse()
//		if doc.Children[0].(*ast.Paragraph).Children[0].(*ast.Strong).Children[0].(*ast.Text).Value != test_input {
//			t.Errorf("expected '%s', got '%s'", test_input, doc.Children[0].(*ast.Paragraph).Children[0].(*ast.Strong).Children[0].(*ast.Text).Value)
//		}
//	})
//}
