// Copyright 2016 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	accumulo "github.com/larsw/accumulo-access-go/pkg"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/cmd"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"os"
)

func main() {
	rego.RegisterBuiltin2(
		&rego.Function{
			Name:             "accumulo.check_authorization",
			Decl:             types.NewFunction(types.Args(types.S, types.S), types.B),
			Memoize:          true,
			Nondeterministic: true,
		},
		func(ctx rego.BuiltinContext, a, b *ast.Term) (*ast.Term, error) {
			var expression, authorizations string

			if err := ast.As(a.Value, &expression); err != nil {
				return nil, err
			} else if err := ast.As(b.Value, &authorizations); err != nil {
				return nil, err
			}

			req, err := accumulo.CheckAuthorization(expression, authorizations)
			if err != nil {
				fmt.Printf("%v", err)
				return nil, err
			}

			return ast.BooleanTerm(req), nil
		},
	)

	if err := cmd.RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}

//go:generate build/gen-run-go.sh internal/cmd/genopacapabilities/main.go capabilities.json
//go:generate build/gen-run-go.sh internal/cmd/genbuiltinmetadata/main.go builtin_metadata.json
//go:generate build/gen-run-go.sh internal/cmd/genversionindex/main.go ast/version_index.json
