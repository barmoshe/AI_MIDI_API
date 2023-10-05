package go_workflow

import (
	"context"

	"go.temporal.io/sdk/activity"
)

func SendMIDITextRequest(ctx context.Context, prompt string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("SendMIDITextRequest", "prompt", prompt)
	return "not implemented", nil
}
