package functions

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestDeepMerge(t *testing.T) {
	tests := []struct {
		name    string
		inputs  []string
		want    map[string]any
		wantErr bool
	}{
		{
			name:   "simple override",
			inputs: []string{`{"a":1}`, `{"b":2}`},
			want:   map[string]any{"a": float64(1), "b": float64(2)},
		},
		{
			name:   "scalar override",
			inputs: []string{`{"a":1}`, `{"a":2}`},
			want:   map[string]any{"a": float64(2)},
		},
		{
			name:   "nested merge",
			inputs: []string{`{"n":{"x":1,"y":2}}`, `{"n":{"y":3,"z":4}}`},
			want:   map[string]any{"n": map[string]any{"x": float64(1), "y": float64(3), "z": float64(4)}},
		},
		{
			name:   "three maps",
			inputs: []string{`{"a":1}`, `{"b":2}`, `{"c":3}`},
			want:   map[string]any{"a": float64(1), "b": float64(2), "c": float64(3)},
		},
		{
			name:   "override nested with scalar",
			inputs: []string{`{"a":{"b":1}}`, `{"a":"replaced"}`},
			want:   map[string]any{"a": "replaced"},
		},
		{
			name:    "invalid json",
			inputs:  []string{`not json`, `{"a":1}`},
			wantErr: true,
		},
		{
			name:   "empty inputs",
			inputs: []string{},
			want:   map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeepMerge(tt.inputs...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DeepMerge() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			var gotMap map[string]any
			if err := json.Unmarshal([]byte(got), &gotMap); err != nil {
				t.Fatalf("result is not valid JSON: %v", err)
			}

			wantJSON, _ := json.Marshal(tt.want)
			gotJSON, _ := json.Marshal(gotMap)
			if string(gotJSON) != string(wantJSON) {
				t.Errorf("DeepMerge() = %s, want %s", gotJSON, wantJSON)
			}
		})
	}
}

func TestDeepMergeFunction_Run(t *testing.T) {
	f := NewDeepMergeFunction()
	ctx := context.Background()

	result := function.NewResultData(basetypes.NewStringNull())
	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue(`{"a":1,"nested":{"x":1}}`),
			types.StringValue(`{"b":2,"nested":{"y":2}}`),
			types.TupleValueMust([]attr.Type{}, []attr.Value{}),
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

	var gotMap map[string]any
	if err := json.Unmarshal([]byte(got.ValueString()), &gotMap); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	nested, ok := gotMap["nested"].(map[string]any)
	if !ok {
		t.Fatal("expected nested to be a map")
	}
	if nested["x"] != float64(1) || nested["y"] != float64(2) {
		t.Errorf("nested merge incorrect: got %v", nested)
	}
}
