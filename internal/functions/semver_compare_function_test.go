package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestSemverCompare(t *testing.T) {
	tests := []struct {
		name    string
		a, b    string
		want    int
		wantErr bool
	}{
		{"equal", "1.2.3", "1.2.3", 0, false},
		{"major less", "1.0.0", "2.0.0", -1, false},
		{"major greater", "2.0.0", "1.0.0", 1, false},
		{"minor less", "1.2.0", "1.3.0", -1, false},
		{"minor greater", "1.3.0", "1.2.0", 1, false},
		{"patch less", "1.2.3", "1.2.4", -1, false},
		{"patch greater", "1.2.4", "1.2.3", 1, false},
		{"prerelease less than release", "1.0.0-alpha", "1.0.0", -1, false},
		{"release greater than prerelease", "1.0.0", "1.0.0-alpha", 1, false},
		{"prerelease ordering", "1.0.0-alpha", "1.0.0-beta", -1, false},
		{"v prefix stripped", "v1.2.3", "1.2.3", 0, false},
		{"build metadata ignored", "1.2.3+build1", "1.2.3+build2", 0, false},
		{"invalid version", "not.a.ver", "1.0.0", 0, true},
		{"too few parts", "1.2", "1.0.0", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SemverCompare(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SemverCompare(%q, %q) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("SemverCompare(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSemverCompareFunction_Run(t *testing.T) {
	f := NewSemverCompareFunction()
	ctx := context.Background()

	result := function.NewResultData(basetypes.NewInt64Null())
	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue("1.2.3"),
			types.StringValue("1.3.0"),
		}),
	}
	resp := function.RunResponse{Result: result}

	f.Run(ctx, req, &resp)

	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error)
	}

	got, ok := resp.Result.Value().(basetypes.Int64Value)
	if !ok {
		t.Fatalf("result is not Int64Value, got %T", resp.Result.Value())
	}
	if got.ValueInt64() != -1 {
		t.Errorf("semver_compare(1.2.3, 1.3.0) = %d, want -1", got.ValueInt64())
	}
}
