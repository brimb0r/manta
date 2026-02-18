package functions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = (*deepMergeFunction)(nil)

type deepMergeFunction struct{}

func NewDeepMergeFunction() function.Function {
	return &deepMergeFunction{}
}

func (f *deepMergeFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "deep_merge"
}

func (f *deepMergeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Recursively merges two or more JSON-encoded maps",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "base",
				Description: "The base JSON-encoded map",
			},
			function.StringParameter{
				Name:        "override",
				Description: "The first override JSON-encoded map",
			},
		},
		VariadicParameter: function.StringParameter{
			Name:        "additional",
			Description: "Additional JSON-encoded maps to merge in order",
		},
		Return: function.StringReturn{},
	}
}

func (f *deepMergeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var base, override string
	var additional []string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &base, &override, &additional))
	if resp.Error != nil {
		return
	}

	all := append([]string{base, override}, additional...)

	result, err := DeepMerge(all...)
	if err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result))
}

// DeepMerge recursively merges JSON-encoded map strings.
// Later values override earlier ones for scalars; nested maps are merged recursively.
func DeepMerge(jsonMaps ...string) (string, error) {
	if len(jsonMaps) == 0 {
		return "{}", nil
	}

	var result map[string]any
	for _, j := range jsonMaps {
		var m map[string]any
		if err := json.Unmarshal([]byte(j), &m); err != nil {
			return "", fmt.Errorf("invalid JSON: %w", err)
		}
		if result == nil {
			result = m
		} else {
			result = deepMergeMaps(result, m)
		}
	}

	out, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to encode result: %w", err)
	}
	return string(out), nil
}

func deepMergeMaps(base, override map[string]any) map[string]any {
	result := make(map[string]any, len(base))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range override {
		if baseVal, exists := result[k]; exists {
			baseMap, baseIsMap := baseVal.(map[string]any)
			overrideMap, overrideIsMap := v.(map[string]any)
			if baseIsMap && overrideIsMap {
				result[k] = deepMergeMaps(baseMap, overrideMap)
				continue
			}
		}
		result[k] = v
	}
	return result
}
