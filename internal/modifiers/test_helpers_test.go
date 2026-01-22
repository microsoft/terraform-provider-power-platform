// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func newPlan(t *testing.T, sch schema.Schema, values map[string]tftypes.Value) tfsdk.Plan {
	t.Helper()
	ctx := context.Background()
	raw := tftypes.NewValue(sch.Type().TerraformType(ctx), values)
	return tfsdk.Plan{Schema: sch, Raw: raw}
}

func newState(t *testing.T, sch schema.Schema, values map[string]tftypes.Value) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	raw := tftypes.NewValue(sch.Type().TerraformType(ctx), values)
	return tfsdk.State{Schema: sch, Raw: raw}
}

func setPrivateData(t *testing.T, target any) reflect.Value {
	t.Helper()
	value := reflect.ValueOf(target)
	if value.Kind() != reflect.Pointer || value.Elem().Kind() != reflect.Struct {
		t.Fatalf("expected pointer to struct, got %T", target)
	}

	privateField := value.Elem().FieldByName("Private")
	if !privateField.IsValid() || privateField.Kind() != reflect.Pointer {
		t.Fatalf("expected Private field to be a pointer on %T", target)
	}

	privatePtr := reflect.New(privateField.Type().Elem())
	privateField.Set(privatePtr)
	return privatePtr
}

func privateSetKey(t *testing.T, privatePtr reflect.Value, ctx context.Context, key string, value []byte) {
	t.Helper()
	method := privatePtr.MethodByName("SetKey")
	if !method.IsValid() {
		t.Fatal("expected SetKey method on private data")
	}

	method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(key), reflect.ValueOf(value)})
}

func privateGetKey(t *testing.T, privatePtr reflect.Value, ctx context.Context, key string) []byte {
	t.Helper()
	method := privatePtr.MethodByName("GetKey")
	if !method.IsValid() {
		t.Fatal("expected GetKey method on private data")
	}

	values := method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(key)})
	if len(values) != 2 {
		t.Fatal("expected GetKey to return two values")
	}
	if !values[1].IsNil() {
		t.Fatalf("expected GetKey error to be nil, got %v", values[1].Interface())
	}

	if values[0].IsNil() {
		return nil
	}

	return values[0].Bytes()
}
