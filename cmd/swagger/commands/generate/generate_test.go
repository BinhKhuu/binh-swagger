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

func Test_SetReturnCode_ArrayType(t *testing.T) {
	res := spec.ResponseSpec{
		Description: "Successful response",
		Schema: &spec.SchemaSpec{
			Type: "array",
			Items: spec.Items{
				Ref: spec.Ref{
					Ref: "#/definitions/User",
				},
			},
		},
	}
	mCdm := map[string]*ModelCommand{
		"User": {
			Name: "User",
			Fields: []spec.FieldSpec{
				{Name: "ID", Type: "int", JSON: "id"},
				{Name: "Name", Type: "string", JSON: "name"},
			},
		},
	}
	ref := "User"
	resCmd := &ResponseCommand{
		Type:        res.Schema.Type,
		Ref:         res.Schema.Items.Ref.Ref,
		Description: res.Description,
	}
	responseKey := 200
	setReturnCode(res, mCdm, ref, resCmd, responseKey)
	code := resCmd.SuccessReturnCode
	expectedCode := "c.JSON(200,[]models.User{})"
	if code != expectedCode {
		t.Errorf("Expected code %q, got %q", expectedCode, code)
	}
}

func Test_SetReturnCode_ObjectType(t *testing.T) {
	res := spec.ResponseSpec{
		Description: "Successful response",
		Schema: &spec.SchemaSpec{
			Type: "object",
			Ref: spec.Ref{
				Ref: "#/definitions/User",
			},
		},
	}
	mCdm := map[string]*ModelCommand{
		"User": {
			Name: "User",
			Fields: []spec.FieldSpec{
				{Name: "ID", Type: "int", JSON: "id"},
				{Name: "Name", Type: "string", JSON: "name"},
			},
		},
	}
	ref := "User"
	resCmd := &ResponseCommand{
		Type:        res.Schema.Type,
		Ref:         res.Schema.Ref.Ref,
		Description: res.Description,
	}
	responseKey := 200
	setReturnCode(res, mCdm, ref, resCmd, responseKey)
	code := resCmd.SuccessReturnCode
	expectedCode := "c.JSON(200, models.User{})"
	if code != expectedCode {
		t.Errorf("Expected code %q, got %q", expectedCode, code)
	}
}
