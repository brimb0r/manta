package functions

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = (*truncateFunction)(nil)

type truncateFunction struct{}

func NewTruncateFunction() function.Function {
	return &truncateFunction{}
}

func (f *truncateFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "truncate"
}

func (f *truncateFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Truncates a string to a maximum length, appending a unique hash suffix",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "input",
				Description: "The string to truncate",
			},
			function.Int64Parameter{
				Name:        "max_length",
				Description: "The maximum allowed length of the output",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *truncateFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var input string
	var maxLength int64
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &input, &maxLength))
	if resp.Error != nil {
		return
	}

	result, err := Truncate(input, int(maxLength))
	if err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result))
}

const hashSuffixLen = 9 // "-" + 8 hex chars

// Truncate shortens s to at most maxLength characters. If truncation is needed
// and maxLength is large enough, a dash and 8-character hex hash of the original
// string is appended to preserve uniqueness. Returns an error if maxLength < 1.
func Truncate(s string, maxLength int) (string, error) {
	if maxLength < 1 {
		return "", fmt.Errorf("max_length must be at least 1, got %d", maxLength)
	}

	if len(s) <= maxLength {
		return s, nil
	}

	// If maxLength is too small for content + hash suffix, just truncate.
	if maxLength <= hashSuffixLen {
		return s[:maxLength], nil
	}

	hash := sha256.Sum256([]byte(s))
	suffix := fmt.Sprintf("-%x", hash[:4])
	return s[:maxLength-hashSuffixLen] + suffix, nil
}
