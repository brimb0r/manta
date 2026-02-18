package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestMask(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		showLast int
		want     string
	}{
		{"mask api key", "sk-1234567890abcdef", 4, "***************cdef"},
		{"show all", "short", 10, "short"},
		{"show none", "secret", 0, "******"},
		{"negative show_last", "secret", -1, "secret"},
		{"empty string", "", 4, ""},
		{"unicode", "\u2603\u2603\u2603\u2603", 2, "**\u2603\u2603"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mask(tt.input, tt.showLast)
			if got != tt.want {
				t.Errorf("Mask(%q, %d) = %q, want %q", tt.input, tt.showLast, got, tt.want)
			}
		})
	}
}

func TestMaskFunction_Run(t *testing.T) {
	f := NewMaskFunction()
	ctx := context.Background()

	result := function.NewResultData(basetypes.NewStringNull())
	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue("sk-1234567890abcdef"),
			types.Int64Value(4),
		}),
	}
	resp := function.RunResponse{Result: result}

	f.Run(ctx, req, &resp)

	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error)
	}

	got, ok := resp.Result.Value().(basetypes.StringValue)
	if !ok {
		t.Fatalf("result is not StringValue, got %T", resp.Result.Value())
	}
	if got.ValueString() != "***************cdef" {
		t.Errorf("mask result = %q, want %q", got.ValueString(), "***************cdef")
	}
}
