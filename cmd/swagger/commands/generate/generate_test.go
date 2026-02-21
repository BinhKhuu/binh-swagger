package generate

import (
	"strings"
	"testing"
)

func Test_SpecToOperationResponse(t *testing.T) {
	op, mCommands := createMockData()

	res, err := SpecToOperationResponses(op, mCommands)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(res) != len(op.Responses) {
		t.Fatalf("Expected %d responses, got %d", len(op.Responses), len(res))
	}

	for key, res := range op.Responses {
		if strings.Compare(res.Description, op.Responses[key].Description) != 0 {
			t.Errorf("Expected description %q, got %q", op.Responses[key].Description, res.Description)
		}

		if op.Responses[key].Schema == nil && res.Schema != nil {
			t.Errorf("Expected no schema, got %v", res.Schema)
		} else if op.Responses[key].Schema != nil {
			if strings.Compare(res.Schema.Type, op.Responses[key].Schema.Type) != 0 {
				t.Errorf("Expected type %q, got %q", op.Responses[key].Schema.Type, res.Schema.Type)
			}

			if condition := res.Schema.Ref != op.Responses[key].Schema.Ref; condition {
				t.Errorf("Expected ref %q, got %q", op.Responses[key].Schema.Ref, res.Schema.Ref)
			}
		}

		// todo add check if type is array and check the items ref
		// todo check ref if type is object
	}
}

func Test_SpecToOperation(t *testing.T) {
	op, mCommands := createMockData()
	res, err := SpecToOperation(op, mCommands)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if res.Summary != op.Summary {
		t.Errorf("Expected summary %q, got %q", op.Summary, res.Summary)
	}

	if res.OperationID != op.OperationID {
		t.Errorf("Expected operation ID %q, got %q", op.OperationID, res.OperationID)
	}

	if len(res.Produces) != len(op.Produces) {
		t.Fatalf("Expected %d produces, got %d", len(op.Produces), len(res.Produces))
	}

	if res.ReturnType != "models.User" {
		t.Errorf("Expected no return type, got %q", res.ReturnType)
	}

	// todo add test if return type is array
}
