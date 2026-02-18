package functions

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = (*semverCompareFunction)(nil)

type semverCompareFunction struct{}

func NewSemverCompareFunction() function.Function {
	return &semverCompareFunction{}
}

func (f *semverCompareFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "semver_compare"
}

func (f *semverCompareFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Compares two semantic version strings, returning -1, 0, or 1",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "version_a",
				Description: "The first semantic version",
			},
			function.StringParameter{
				Name:        "version_b",
				Description: "The second semantic version",
			},
		},
		Return: function.Int64Return{},
	}
}

func (f *semverCompareFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var versionA, versionB string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &versionA, &versionB))
	if resp.Error != nil {
		return
	}

	result, err := SemverCompare(versionA, versionB)
	if err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, int64(result)))
}

type semver struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
}

func parseSemver(s string) (semver, error) {
	s = strings.TrimPrefix(s, "v")

	// Strip build metadata (ignored in precedence).
	if idx := strings.Index(s, "+"); idx != -1 {
		s = s[:idx]
	}

	// Extract pre-release.
	var pre string
	if idx := strings.Index(s, "-"); idx != -1 {
		pre = s[idx+1:]
		s = s[:idx]
	}

	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return semver{}, fmt.Errorf("invalid semver %q: expected major.minor.patch", s)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return semver{}, fmt.Errorf("invalid major version %q: %w", parts[0], err)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return semver{}, fmt.Errorf("invalid minor version %q: %w", parts[1], err)
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return semver{}, fmt.Errorf("invalid patch version %q: %w", parts[2], err)
	}

	return semver{Major: major, Minor: minor, Patch: patch, Prerelease: pre}, nil
}

// SemverCompare compares two semantic version strings.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
func SemverCompare(a, b string) (int, error) {
	va, err := parseSemver(a)
	if err != nil {
		return 0, err
	}
	vb, err := parseSemver(b)
	if err != nil {
		return 0, err
	}

	if c := cmpInt(va.Major, vb.Major); c != 0 {
		return c, nil
	}
	if c := cmpInt(va.Minor, vb.Minor); c != 0 {
		return c, nil
	}
	if c := cmpInt(va.Patch, vb.Patch); c != 0 {
		return c, nil
	}

	// A version with pre-release has lower precedence than the release version.
	switch {
	case va.Prerelease == "" && vb.Prerelease == "":
		return 0, nil
	case va.Prerelease == "":
		return 1, nil
	case vb.Prerelease == "":
		return -1, nil
	default:
		return strings.Compare(va.Prerelease, vb.Prerelease), nil
	}
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
