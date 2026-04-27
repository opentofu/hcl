package typeexpr

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type TypeContext struct {
	Types    map[string]map[string]cty.Type
	Defaults map[string]map[string]*Defaults
}

// Type attempts to process the given expression as a type expression and, if
// successful, returns the resulting type. If unsuccessful, error diagnostics
// are returned.
func (t TypeContext) Type(expr hcl.Expression) (cty.Type, hcl.Diagnostics) {
	ty, _, diags := getType(expr, false, false, t)
	return ty, diags
}

// TypeConstraint attempts to parse the given expression as a type constraint
// and, if successful, returns the resulting type. If unsuccessful, error
// diagnostics are returned.
//
// A type constraint has the same structure as a type, but it additionally
// allows the keyword "any" to represent cty.DynamicPseudoType, which is often
// used as a wildcard in type checking and type conversion operations.
func (t TypeContext) TypeConstraint(expr hcl.Expression) (cty.Type, hcl.Diagnostics) {
	ty, _, diags := getType(expr, true, false, t)
	return ty, diags
}

// TypeConstraintWithDefaults attempts to parse the given expression as a type
// constraint which may include default values for object attributes. If
// successful both the resulting type and corresponding defaults are returned.
// If unsuccessful, error diagnostics are returned.
func (t TypeContext) TypeConstraintWithDefaults(expr hcl.Expression) (cty.Type, *Defaults, hcl.Diagnostics) {
	return getType(expr, true, true, t)
}

func (t TypeContext) TypeDependencies(expr hcl.Expression) (map[string][]string, hcl.Diagnostics) {
	_, _, diags := getType(expr, true, true, t)

	missing := map[string][]string{}
	var filtered hcl.Diagnostics

	for _, diag := range diags {
		if tm, ok := diag.Extra.(DiagnosticExtraTypeMissing); ok {
			missing[tm.Namespace] = append(missing[tm.Namespace], tm.TypeName)
		} else {
			filtered = append(filtered, diag)
		}
	}

	return missing, filtered
}

type DiagnosticExtraTypeMissing struct {
	Namespace string
	TypeName  string
}
