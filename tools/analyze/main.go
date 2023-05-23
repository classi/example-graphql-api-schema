package main

import (
	"github.com/classi/example-graphql-api-schema/analyzer/requireauthorize"
	"github.com/gqlgo/gqlanalysis/multichecker"
)

func main() {
	multichecker.Main(requireauthorize.Analyzer)
}
