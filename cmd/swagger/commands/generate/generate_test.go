package generate

import (
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"strings"
	"testing"
)

func Test_SpecToModelCommand(t *testing.T) {
	fieldSpecs := []spec.FieldSpec{
		{Name: "ID", Type: "int", JSON: "id"},
		{Name: "Name", Type: "string", JSON: "name"},
	}

	modelSpec := spec.ModelSpec{
		Name:   "User",
		Fields: fieldSpecs,
	}

	res, err := SpecToModelCommand(modelSpec)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if strings.Compare(res.Name, modelSpec.Name) != 0 {
		t.Errorf("Expected name %q, got %q", modelSpec.Name, res.Name)
	}

	if len(res.Fields) != len(modelSpec.Fields) {
		t.Fatalf("Expected %d fields, got %d", len(modelSpec.Fields), len(res.Fields))
	}

	for i, field := range modelSpec.Fields {
		if strings.Compare(field.Name, res.Fields[i].Name) != 0 {
			t.Errorf("Expected field name %q at index %d, got %q", field.Name, i, res.Fields[i].Name)
		}
		if strings.Compare(field.Type, res.Fields[i].Type) != 0 {
			t.Errorf("Expected field type %q at index %d, got %q", field.Type, i, res.Fields[i].Type)
		}
		if strings.Compare(field.JSON, res.Fields[i].JSON) != 0 {
			t.Errorf("Expected field JSON tag %q at index %d, got %q", field.JSON, i, res.Fields[i].JSON)
		}
	}
}

func Test_SpecToOperation(t *testing.T) {
	op, mCommands := createMockDataWithObjectResponse()

	res, err := SpecToOperation(op, mCommands)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if strings.Compare(res.Summary, op.Summary) != 0 {
		t.Errorf("Expected summary %q, got %q", op.Summary, res.Summary)
	}

	if strings.Compare(res.OperationID, op.OperationID) != 0 {
		t.Errorf("Expected operation ID %q, got %q", op.OperationID, res.OperationID)
	}

	if len(res.Produces) != len(op.Produces) {
		t.Fatalf("Expected %d produces, got %d", len(op.Produces), len(res.Produces))
	}

	for i, produce := range op.Produces {
		if strings.Compare(produce, res.Produces[i]) != 0 {
			t.Errorf("Expected produce %q at index %d, got %q", produce, i, res.Produces[i])
		}
	}

	// relies on deriveReturnType and combineReturnCodes working correctly, but we can still check if they produce the expected results
	if strings.Compare(res.ReturnType, deriveReturnType(op)) != 0 {
		t.Errorf("Expected return type %q, got %q", deriveReturnType(op), res.ReturnType)
	}

	// relies on combineReturnCodes working correctly, but we can still check if it produces the expected result
	if res.ReturnCode != combineReturnCodes(res.Responses) {
		t.Errorf("Expected return code %q, got %q", combineReturnCodes(res.Responses), res.ReturnCode)
	}
}

func Test_SpecToOperationResponses_ObjectType(t *testing.T) {
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

func Test_SpecToOperationResponses_ArrayType(t *testing.T) {
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
