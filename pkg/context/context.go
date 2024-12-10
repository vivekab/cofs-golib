package golibcontext

import (
	"context"

	golibconstants "github.com/vivekab/golib/pkg/constants"
	"google.golang.org/grpc/metadata"
)

// GetFromContext will check for given key in incomingContext, outgoingContext and ctx till a value if found
func GetFromContext(ctx context.Context, key string) string {
	var value string
	var found bool

	// Check incoming context first
	incomingMetadata, ok := metadata.FromIncomingContext(ctx)
	if ok {
		incomingValues, exists := incomingMetadata[key]
		if exists && len(incomingValues) > 0 {
			value = incomingValues[0]
			found = true
		}
	}

	// If not found, check outgoing context
	if !found {
		outgoingMetadata, ok := metadata.FromOutgoingContext(ctx)
		if ok {
			outgoingValues, exists := outgoingMetadata[key]
			if exists && len(outgoingValues) > 0 {
				value = outgoingValues[0]
				found = true
			}
		}
	}

	// If still not found, check the context directly
	if !found {
		valueStr, ok := ctx.Value(key).(string)
		if ok {
			value = valueStr
		}
	}

	return value
}

func GetRequestIdFromContext(ctx context.Context) string {
	return GetFromContext(ctx, golibconstants.HeaderRequestID)
}
