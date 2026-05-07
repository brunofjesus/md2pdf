package node

import (
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown/ast"
)

// ProcessLink handles *ast.Link entering/leaving.
func ProcessLink(ctx PdfContext, n ast.Node, entering bool) {
	node := n.(*ast.Link)
	destination := string(node.Destination)
	if entering {
		if ctx.GetInputBaseURL() != "" && !strings.HasPrefix(destination, "http") {
			destination = ctx.GetInputBaseURL() + "/" + strings.Replace(destination, "./", "", 1)
		}
		x := &ContainerState{
			TextStyle:   ctx.GetTheme().Link,
			ListKind:    NotList,
			LeftMargin:  ctx.PeekState().LeftMargin,
			Destination: destination,
		}
		ctx.PushState(x)
		ctx.Tracer("Link (entering)",
			fmt.Sprintf("Destination[%v] Title[%v]",
				string(node.Destination),
				string(node.Title)))
	} else {
		ctx.Tracer("Link (leaving)", "")
		ctx.PopState()
	}
}
