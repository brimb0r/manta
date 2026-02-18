package provider

import (
	"context"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = (*isPalindromeFunction)(nil)

type isPalindromeFunction struct{}

func NewIsPalindromeFunction() function.Function {
	return &isPalindromeFunction{}
}

func (f *isPalindromeFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "is_palindrome"
}

func (f *isPalindromeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Checks whether a string is a palindrome",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "input",
				Description: "The string to check",
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f *isPalindromeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var input string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &input))
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, IsPalindrome(input)))
}

// IsPalindrome checks whether s is a palindrome.
// It is case-insensitive and ignores non-alphanumeric characters.
func IsPalindrome(s string) bool {
	s = strings.ToLower(s)

	// Keep only alphanumeric runes.
	filtered := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			filtered = append(filtered, r)
		}
	}

	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		if filtered[i] != filtered[j] {
			return false
		}
	}
	return true
}
