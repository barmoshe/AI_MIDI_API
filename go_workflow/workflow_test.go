package go_workflow

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_GenerateMIDIWorkflow_bestCase(t *testing.T) {

	// Create a new test workflow environment
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	//logger := testSuite.GetLogger()
	env.RegisterActivity(SendMIDITextRequest)
	env.RegisterActivity(ValidateMIDIText)
	env.RegisterActivity(GenerateMIDIFile)

	// Mock activity implementations
	mockMIDIText := "Mocked MIDIText"
	mockS3Link := "mocked-s3-link"

	// Mock the "SendMIDITextRequest" activity
	env.OnActivity("SendMIDITextRequest", mock.Anything).Return(mockMIDIText, nil)

	// Mock the "ValidateMIDIText" activity
	env.OnActivity("ValidateMIDIText", mockMIDIText).Return(true, nil)

	// Mock the "GenerateMIDIFile" activity
	env.OnActivity("GenerateMIDIFile", mockMIDIText).Return(mockS3Link, nil)

	// Execute the workflow function with a mock input
	input := MIDIRequest{Prompt: "Test Prompt", RequestType: "Test Request Type"}
	env.ExecuteWorkflow(GenerateMIDIWorkflow, input)

	// Check workflow completion and errors
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())

	// Check the workflow result
	var result string
	require.NoError(t, env.GetWorkflowResult(&result))

	// Verify the expected S3 link in the result
	expectedJSON := `{"link":"mocked-s3-link"}`
	require.Equal(t, expectedJSON, result)
}

func TestGenerateMIDIWorkflow_ValidationFails(t *testing.T) {
	// Create a new test workflow environment
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	env.RegisterActivity(SendMIDITextRequest)
	env.RegisterActivity(ValidateMIDIText)
	env.RegisterActivity(GenerateMIDIFile)

	// Mock activity implementations
	mockMIDIText := "Mocked MIDIText"

	// Mock the "SendMIDITextRequest" activity
	env.OnActivity("SendMIDITextRequest", mock.Anything).Return(mockMIDIText, nil)

	// Mock the "ValidateMIDIText" activity
	env.OnActivity("ValidateMIDIText", mockMIDIText).Return(false, nil)

	// Execute the workflow function with a mock input
	input := MIDIRequest{Prompt: "Test Prompt", RequestType: "Test Request Type"}
	env.ExecuteWorkflow(GenerateMIDIWorkflow, input)

	// Check workflow completion and errors
	require.True(t, env.IsWorkflowCompleted())

	// Check for an error in the workflow
	require.Error(t, env.GetWorkflowError())

	// Verify the expected error message
	expectedErrorMessage := "unable to validate MIDI text"
	require.Contains(t, env.GetWorkflowError().Error(), expectedErrorMessage)
}

func Test_GenerateMIDIWorkflow_GenerateMIDIFileFails(t *testing.T) {
	// Create a new test workflow environment
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	// logger := testSuite.GetLogger()
	env.RegisterActivity(SendMIDITextRequest)
	env.RegisterActivity(ValidateMIDIText)
	env.RegisterActivity(GenerateMIDIFile)

	// Mock activity implementations
	mockMIDIText := "Mocked MIDIText"

	// Mock the "SendMIDITextRequest" activity
	env.OnActivity("SendMIDITextRequest", mock.Anything).Return(mockMIDIText, nil)

	// Mock the "ValidateMIDIText" activity
	env.OnActivity("ValidateMIDIText", mockMIDIText).Return(true, nil)

	// Mock the "GenerateMIDIFile" activity
	env.OnActivity("GenerateMIDIFile", mockMIDIText).Return("", nil)

	// Execute the workflow function with a mock input
	input := MIDIRequest{Prompt: "Test Prompt", RequestType: "Test Request Type"}
	env.ExecuteWorkflow(GenerateMIDIWorkflow, input)

	// Check workflow completion and errors
	require.True(t, env.IsWorkflowCompleted())

	// Check for an error in the workflow
	require.Error(t, env.GetWorkflowError())

	// Verify the expected error message
	expectedErrorMessage := "unable to generate MIDI file"
	require.Contains(t, env.GetWorkflowError().Error(), expectedErrorMessage)
}

func Test_GenerateMIDIWorkflow_EdgeCases(t *testing.T) {
	testCases := []struct {
		Name       string
		Input      MIDIRequest
		Expected   string // Expected workflow result
		Activities map[string]struct {
			ReturnVal interface{}
			Err       error
		}
	}{
		{
			Name:     "ValidInput",
			Input:    MIDIRequest{Prompt: "Test Prompt", RequestType: "Test Request Type"},
			Expected: `{"link":"mocked-s3-link"}`,
			Activities: map[string]struct {
				ReturnVal interface{}
				Err       error
			}{
				"SendMIDITextRequest": {ReturnVal: "Mocked MIDIText", Err: nil},
				"ValidateMIDIText":    {ReturnVal: true, Err: nil},
				"GenerateMIDIFile":    {ReturnVal: "mocked-s3-link", Err: nil},
			},
		},
		{
			Name:     "ValidationFails",
			Input:    MIDIRequest{Prompt: "Test Prompt", RequestType: "Test Request Type"},
			Expected: "",
			Activities: map[string]struct {
				ReturnVal interface{}
				Err       error
			}{
				"SendMIDITextRequest": {ReturnVal: "Mocked MIDIText", Err: nil},
				"ValidateMIDIText":    {ReturnVal: false, Err: nil},
				"GenerateMIDIFile":    {ReturnVal: "", Err: nil},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestWorkflowEnvironment()
			env.RegisterActivity(SendMIDITextRequest)
			env.RegisterActivity(ValidateMIDIText)
			env.RegisterActivity(GenerateMIDIFile)

			for activityName, mockResult := range tc.Activities {
				env.OnActivity(activityName, mock.Anything).Return(mockResult.ReturnVal, mockResult.Err)
			}

			env.ExecuteWorkflow(GenerateMIDIWorkflow, tc.Input)

			require.True(t, env.IsWorkflowCompleted())
			if tc.Expected != "" {
				var result string
				require.NoError(t, env.GetWorkflowResult(&result))
				require.Equal(t, tc.Expected, result)
			} else {
				require.Error(t, env.GetWorkflowError())
			}
		})
	}
}

// mock activity implementations

func ValidateMIDIText(midiText string) (bool, error) {
	return true, nil
}
func GenerateMIDIFile(midiText string) (string, error) {
	return "mocked-s3-link", nil
}
