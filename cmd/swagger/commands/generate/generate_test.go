package generate

import (
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"strings"
	"testing"
)

func Test_SpecToOperation_ObjectType(t *testing.T) {
	op, mCommands := createMockDataWithObjectResponse()

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
	}
}

func Test_SpecToOperation_ArrayType(t *testing.T) {
	op, mCommands := createMockDataWithArrayResponse()

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

			if condition := res.Schema.Items.Ref != op.Responses[key].Schema.Items.Ref; condition {
				t.Errorf("Expected items ref %q, got %q", op.Responses[key].Schema.Items.Ref, res.Schema.Items.Ref)
			}
		}
	}
}

func Test_SetReturnCode(t *testing.T) {
	tests := []struct {
		name         string
		schema       *spec.SchemaSpec
		expectedCode string
	}{
		{
			name: "array",
			schema: &spec.SchemaSpec{
				Type: "array",
				Items: spec.Items{
					Ref: spec.Ref{Ref: "#/definitions/User"},
				},
			},
			expectedCode: "c.JSON(200,[]models.User{})",
		},
		{
			name: "object",
			schema: &spec.SchemaSpec{
				Type: "object",
				Ref:  spec.Ref{Ref: "#/definitions/User"},
			},
			expectedCode: "c.JSON(200, models.User{})",
		},
	}

	mCdm := &ModelCommand{
		Name: "User",
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := spec.ResponseSpec{
				Description: "Successful response",
				Schema:      tc.schema,
			}
			resCmd := &ResponseCommand{
				Type:        res.Schema.Type,
				Description: res.Description,
			}
			if res.Schema.Type == "array" {
				resCmd.Ref = res.Schema.Items.Ref.Ref
			} else {
				resCmd.Ref = res.Schema.Ref.Ref
			}

			setReturnCode(res, mCdm, resCmd, 200)

			if resCmd.SuccessReturnCode != tc.expectedCode {
				t.Errorf("expected %q, got %q", tc.expectedCode, resCmd.SuccessReturnCode)
			}
		})
	}
}

func Test_CombineReturnCodes(t *testing.T) {
	responses := map[int]ResponseCommand{
		200: {SuccessReturnCode: "c.JSON(200, models.User{})"},
		400: {SuccessReturnCode: "c.JSON(400, models.Error{})"},
		500: {SuccessReturnCode: ""},
	}

	expected := "c.JSON(200, models.User{})\nc.JSON(400, models.Error{})"
	combined := combineReturnCodes(responses)
	if combined != expected {
		t.Errorf("expected %q, got %q", expected, combined)
	}
}
