```golang

sahar.Block(
    sahar.Stack,
    sahar.A4(), // set the width and height of the block

    sahar.Block(
        sahar.Group,
        sahar.Height(100),
        sahar.Block(
            sahar.Stack,
            sahar.Width(100),
        ),
        sahar.Block(
            sahar.Stack,
            sahar.Alignments(sahar.Left, sahar.Top),
            sahar.Block(sahar.Stack),
        ),
    ),

    sahar.Block(
        sahar.Stack,
        sahar.Block(
            sahar.Stack,
            sahar.Height(100),
            sahar.Background("red"),
        ),
        sahar.Block(sahar.Stack),
        sahar.Block(
            sahar.Group,
            sahar.Height(10),
            sahar.Block(
                sahar.Stack,
                sahar.FontSize(14),
                sahar.Text(
                    "Contant me at",
                ),
            ),
            sahar.Block(
                sahar.Stack,
                sahar.FontSize(14),
                sahar.Text("sahar[at]ella.to"),
                sahar.Attr("href", "mailto:sahar@ella.to"),
            ),
        ),
    ),
)

```
