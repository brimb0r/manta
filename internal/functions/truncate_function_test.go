package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		wantErr   bool
		check     func(t *testing.T, got string)
	}{
		{
			name:      "no truncation needed",
			input:     "short",
			maxLength: 10,
			check:     func(t *testing.T, got string) { assertEqual(t, got, "short") },
		},
		{
			name:      "exact length",
			input:     "exact",
			maxLength: 5,
			check:     func(t *testing.T, got string) { assertEqual(t, got, "exact") },
		},
		{
			name:      "truncated with hash",
			input:     "this-is-a-very-long-resource-name-that-needs-truncation",
			maxLength: 24,
			check: func(t *testing.T, got string) {
				if len(got) != 24 {
					t.Errorf("len = %d, want 24", len(got))
				}
				// Should end with -XXXXXXXX (hex hash)
				if got[15] != '-' {
					t.Errorf("expected dash at position 15, got %q", got)
				}
			},
		},
		{
			name:      "deterministic hash",
			input:     "this-is-a-very-long-resource-name-that-needs-truncation",
			maxLength: 24,
			check: func(t *testing.T, got string) {
				got2, _ := Truncate("this-is-a-very-long-resource-name-that-needs-truncation", 24)
				assertEqual(t, got, got2)
			},
		},
		{
			name:      "very small max_length",
			input:     "abcdefghijklmnop",
			maxLength: 5,
			check: func(t *testing.T, got string) {
				if len(got) != 5 {
					t.Errorf("len = %d, want 5", len(got))
				}
			},
		},
		{
			name:      "max_length zero",
			input:     "test",
			maxLength: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Truncate(tt.input, tt.maxLength)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Truncate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTruncateFunction_Run(t *testing.T) {
	f := NewTruncateFunction()
	ctx := context.Background()

	result := function.NewResultData(basetypes.NewStringNull())
	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue("my-very-long-resource-name-that-exceeds-the-limit"),
			types.Int64Value(24),
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
	if len(got.ValueString()) != 24 {
		t.Errorf("len = %d, want 24", len(got.ValueString()))
	}
}
