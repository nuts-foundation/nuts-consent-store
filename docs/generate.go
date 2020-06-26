package main

import (
	"github.com/nuts-foundation/nuts-consent-store/engine"
	"github.com/nuts-foundation/nuts-go-core/docs"
)

func main() {
	docs.GenerateConfigOptionsDocs("README_options.rst", engine.NewConsentStoreEngine().FlagSet)
}
