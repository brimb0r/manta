package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestMantaProvider_Metadata(t *testing.T) {
	p := New()

	req := provider.MetadataRequest{}
	resp := provider.MetadataResponse{}
	p.Metadata(context.Background(), req, &resp)

	if resp.TypeName != "manta" {
		t.Errorf("TypeName = %q, want %q", resp.TypeName, "manta")
	}
}

func TestMantaProvider_Schema(t *testing.T) {
	p := New()

	req := provider.SchemaRequest{}
	resp := provider.SchemaResponse{}
	p.Schema(context.Background(), req, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %s", resp.Diagnostics)
	}

	attr, ok := resp.Schema.Attributes["endpoint"]
	if !ok {
		t.Fatal("schema missing 'endpoint' attribute")
	}

	if attr.IsRequired() {
		t.Error("endpoint attribute should be optional, not required")
	}
}

func TestMantaProvider_Functions(t *testing.T) {
	p := New()

	pf, ok := p.(provider.ProviderWithFunctions)
	if !ok {
		t.Fatal("provider does not implement ProviderWithFunctions")
	}

	funcs := pf.Functions(context.Background())

	// Collect registered function names.
	registered := make(map[string]bool)
	for _, fn := range funcs {
		f := fn()
		var metaResp function.MetadataResponse
		f.Metadata(context.Background(), function.MetadataRequest{}, &metaResp)
		registered[metaResp.Name] = true
	}

	expected := []string{"deep_merge", "is_palindrome", "mask", "semver_compare", "truncate"}
	for _, name := range expected {
		if !registered[name] {
			t.Errorf("function %q not found in provider functions", name)
		}
	}
}
