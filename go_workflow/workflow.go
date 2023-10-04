package go_workflow

import (
	"encoding/json"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

type MIDIRequest struct {
	Prompt      string
	RequestType string
}

func GenerateMIDIWorkflow(ctx workflow.Context, input MIDIRequest) (string, error) {
	activityoptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityoptions)
	var midiText string
	var validated bool = false
	var validationCount = 0
	for (!validated) && validationCount < 5 {
		validationCount++
		//typescript activity
		midiTextErr := workflow.ExecuteActivity(ctx, "SendMIDITextRequest", input.Prompt).Get(ctx, &midiText)
		if midiTextErr != nil {
			return "", midiTextErr

		}
		//Python activity
		validationErr := workflow.ExecuteActivity(ctx, "ValidateMIDIText", midiText).Get(ctx, &validated)
		if validationErr != nil {
			return "", validationErr

		}
	}
	if !validated {

		return "", fmt.Errorf("unable to validate MIDI text")
	}

	var s3Link string = ""
	var generationCount = 0
	for s3Link == "" && generationCount < 5 {
		//Python activity
		generateErr := workflow.ExecuteActivity(ctx, "GenerateMIDIFile", midiText).Get(ctx, &s3Link)
		if generateErr != nil {
			return "", generateErr
		}
		generationCount++
	}
	if s3Link != "" {
		jsonResult := struct {
			Link string `json:"link"`
		}{
			Link: s3Link,
		}

		result, err := json.Marshal(jsonResult)
		if err != nil {
			return "", err
		}
		return string(result), nil
	}
	return "", fmt.Errorf("unable to generate MIDI file")

}
