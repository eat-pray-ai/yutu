package pkg

import "github.com/jedib0t/go-pretty/v6/table"

const (
	PartsUsage  = "Comma separated parts"
	MRUsage     = "The maximum number of items that should be returned, 0 for no limit"
	TableUsage  = "json|yaml|table"
	SilentUsage = "json|yaml|silent"
	JPUsage     = "JSONPath expression to filter the output"
	JsonMIME    = "application/json"
	PerPage     = 20

	getWdFailed    = "failed to get working directory"
	openRootFailed = "failed to open root directory"
)

var (
	TableStyle = table.Style{
		Name:   "StyleLight",
		Box:    table.StyleBoxLight,
		Color:  table.ColorOptionsDefault,
		Format: table.FormatOptionsDefault,
		HTML:   table.DefaultHTMLOptions,
		Options: table.Options{
			DrawBorder:      false,
			SeparateColumns: false,
			SeparateFooter:  false,
			SeparateHeader:  false,
			SeparateRows:    false,
		},
		Size:  table.SizeOptionsDefault,
		Title: table.TitleOptionsDefault,
	}
)
