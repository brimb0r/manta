package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = (*maskFunction)(nil)

type maskFunction struct{}

func NewMaskFunction() function.Function {
	return &maskFunction{}
}

func (f *maskFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "mask"
}

func (f *maskFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Masks a string, revealing only the last N characters",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "input",
				Description: "The string to mask",
			},
			function.Int64Parameter{
				Name:        "show_last",
				Description: "The number of trailing characters to leave visible",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *maskFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var input string
	var showLast int64
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &input, &showLast))
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, Mask(input, int(showLast))))
}

// Mask replaces all but the last showLast characters with asterisks.
func Mask(s string, showLast int) string {
	runes := []rune(s)
	if showLast >= len(runes) || showLast < 0 {
		return s
	}
	masked := len(runes) - showLast
	return strings.Repeat("*", masked) + string(runes[masked:])
}
