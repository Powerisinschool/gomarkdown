package ast

// Node is the base interface for all AST nodes
type Node interface {
	// A marker method
	node()
}

// BlockNode represents block-level nodes like paragraphs and headings
type BlockNode interface {
	Node
	blockNode()
}

// InlineNode represents inline-level nodes like text and links
type InlineNode interface {
	Node
	inlineNode()
}

// Document represents the root of the AST
type Document struct {
	Children []BlockNode
}

func (d *Document) node()      {}
func (d *Document) blockNode() {}

func NewDocument() *Document {
	return &Document{}
}

type Paragraph struct {
	Children []InlineNode
}

func (p *Paragraph) node()      {}
func (p *Paragraph) blockNode() {}

type Text struct {
	Value string
}

func (t *Text) node()       {}
func (t *Text) inlineNode() {}

type Strong struct {
	Children []InlineNode
}

func (s *Strong) node()       {}
func (s *Strong) inlineNode() {}

type Heading struct {
	Level    int
	Children []InlineNode
}

func (h *Heading) node()      {}
func (h *Heading) blockNode() {}

type List struct {
	Items []*ListItem
}

func (l *List) node()      {}
func (l *List) blockNode() {}

type ListItem struct {
	Children []InlineNode
}

func (li *ListItem) node() {}
