package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"simple palindrome", "racecar", true},
		{"single character", "a", true},
		{"empty string", "", true},
		{"mixed case", "RaceCar", true},
		{"with spaces", "taco cat", true},
		{"with punctuation", "A man, a plan, a canal: Panama", true},
		{"numeric palindrome", "12321", true},
		{"not a palindrome", "hello", false},
		{"almost palindrome", "abcba1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPalindrome(tt.input)
			if got != tt.want {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsPalindromeFunction_Run(t *testing.T) {
	f := NewIsPalindromeFunction()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"palindrome via function", "racecar", true},
		{"non-palindrome via function", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := function.NewResultData(basetypes.NewBoolNull())

			req := function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue(tt.input),
				}),
			}
			resp := function.RunResponse{
				Result: result,
			}

			f.Run(context.Background(), req, &resp)

			if resp.Error != nil {
				t.Fatalf("unexpected error: %s", resp.Error)
			}

			got, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok {
				t.Fatalf("result is not a BoolValue, got %T", resp.Result.Value())
			}

			if got.ValueBool() != tt.want {
				t.Errorf("Run(%q) result = %v, want %v", tt.input, got.ValueBool(), tt.want)
			}
		})
	}
}
