# Incorrect `MarshallTo` usage of `json.NewDecoder().Decode()`

##

/workspaces/terraform-provider-power-platform/internal/api/request.go

## Problem

The code passes `&obj` to `json.NewDecoder().Decode()`, but `obj` is already of type `any` (interface{}), typically a pointer to struct should be passed, not its address. Taking `&obj` gives a pointer to interface, which is almost always incorrect and leads to runtime errors.

## Impact

Severity: High

This can cause decoding to fail, with JSON not being correctly deserialized into the provided struct.

## Location

```go
func (apiResponse *Response) MarshallTo(obj any) error {
	err := json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(&obj)
	if err != nil {
		return err
	}
	return nil
}
```

## Code Issue

```go
	err := json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(&obj)
```

## Fix

Pass `obj` directly to Decode (it should already be a pointer):

```go
func (apiResponse *Response) MarshallTo(obj any) error {
	return json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(obj)
}
```
