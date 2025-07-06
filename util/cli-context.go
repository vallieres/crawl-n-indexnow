package util

import (
	"context"
	"errors"
)

type cliContextKey string

// CLIContext will hold the services so that they can be accessed in various commands and sub-commands
// when run with ExecuteContext(ctx).
type CLIContext struct {
	Domain      *string
	IndexNowKey *string
	context.Context
}

// NewCLIContext initializes the CLI context with provided services.
func NewCLIContext(domain *string, indexNowKey *string) *CLIContext {
	return &CLIContext{
		Domain:      domain,
		IndexNowKey: indexNowKey,
	}
}

// WithService initializes the CLIContext with the provided services.
func WithService(ctx context.Context, domain *string, indexNowKey *string) context.Context {
	cliCtx := &CLIContext{
		Domain:      domain,
		IndexNowKey: indexNowKey,
	}
	return context.WithValue(ctx, cliContextKey("cliContext"), cliCtx)
}

// GetDomain will pull the Domain from the CLIContext and return it.
func GetDomain(ctx context.Context) (*string, error) {
	cliCtx, ok := ctx.Value(cliContextKey("cliContext")).(*CLIContext)
	if !ok {
		return nil, errors.New("failed to get the DOmain through CLIContext from context")
	}

	return cliCtx.Domain, nil
}

// GetIndexNowKey will pull the IndexNowKey from the CLIContext and return it.
func GetIndexNowKey(ctx context.Context) (*string, error) {
	cliCtx, ok := ctx.Value(cliContextKey("cliContext")).(*CLIContext)
	if !ok {
		return nil, errors.New("failed to get the IndexNowKey through CLIContext from context")
	}

	return cliCtx.IndexNowKey, nil
}
