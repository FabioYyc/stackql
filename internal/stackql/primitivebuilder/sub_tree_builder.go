package primitivebuilder

import "github.com/stackql/stackql/internal/stackql/primitivegraph"

type SubTreeBuilder struct {
	indirects []Builder
	children  []Builder
}

func NewSubTreeBuilder(children []Builder) Builder {
	return &SubTreeBuilder{
		children: children,
	}
}

func (st *SubTreeBuilder) Build() error {
	for _, child := range st.children {
		err := child.Build()
		if err != nil {
			return err
		}
	}
	return nil
}

func (st *SubTreeBuilder) GetRoot() primitivegraph.PrimitiveNode {
	return st.children[0].GetRoot()
}

func (st *SubTreeBuilder) GetTail() primitivegraph.PrimitiveNode {
	return st.children[len(st.children)-1].GetTail()
}
