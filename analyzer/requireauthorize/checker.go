package requireauthorize

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/gqlgo/gqlanalysis"
	"github.com/vektah/gqlparser/v2/ast"
)

var Analyzer = &gqlanalysis.Analyzer{
	Name: "requireauthorizescope",
	Run:  run,
}

const directiveName = "authorize"

func run(pass *gqlanalysis.Pass) (any, error) {
	if pass.Schema.Query != nil {
		if err := check(pass, fieldKindQuery, pass.Schema.Query.Fields); err != nil {
			return nil, err
		}
	}
	if pass.Schema.Mutation != nil {
		if err := check(pass, fieldKindMutation, pass.Schema.Mutation.Fields); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func check(pass *gqlanalysis.Pass, fieldKind fieldKind, fields ast.FieldList) error {
	for _, f := range fields {
		if f.Position == nil { // injected
			continue
		}
		dir := f.Directives.ForName(directiveName)
		if dir == nil {
			pass.Reportf(f.Position, "directive %q not found", directiveName)
			continue
		}
		scopes := dir.Arguments.ForName("scopes")
		if scopes == nil {
			pass.Reportf(f.Position, "[BUG] scopes argument is not found")
			continue
		}
		expectedScope := buildExpectedScopeName(fieldKind, f.Name)
		var found bool
		for _, c := range scopes.Value.Children {
			if c.Value.Raw == expectedScope {
				found = true
				break
			}
		}
		if !found {
			pass.Reportf(f.Position, "the required scopes must include %q but no scopes found", expectedScope)
		}
	}
	return nil
}

func buildExpectedScopeName(fk fieldKind, fieldName string) string {
	b := new(bytes.Buffer)
	fmt.Fprint(b, fk.prefix())
	for _, r := range fieldName {
		if unicode.IsUpper(r) {
			b.WriteRune('_')
		}
		b.WriteRune(unicode.ToUpper(r))
	}
	return b.String()
}

type fieldKind int

func (fk fieldKind) prefix() string {
	switch fk {
	case fieldKindQuery:
		return "QUERY_"
	case fieldKindMutation:
		return "MUTATION_"
	default:
		return "UNKNOWN_"
	}
}

const (
	fieldKindQuery fieldKind = iota + 1
	fieldKindMutation
)
