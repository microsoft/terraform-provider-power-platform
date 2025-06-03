# Title

Direct Use of json.Unmarshal Without Type Checking or Unknown Field Handling

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/api_currencies.go

## Problem

You are using `json.Unmarshal` to populate the DTO, but there is no strict error check for unknown or unexpected fields unless it's built into your `currenciesDto`. If your struct is meant to resist API drift or you want to detect changes in the returned payloads, stricter checks or logging about unknown fields are prudent.

## Impact

**Low** (potentially higher if the API changes unexpectedly). Maintainability and robustness.

## Location

Unmarshalling JSON:

```go
err = json.Unmarshal(response.BodyAsBytes, &currencies)

if err != nil {
	return currencies, err
}
```

## Fix

For stricter checks, you can use a decoder with `DisallowUnknownFields()`:

```go
decoder := json.NewDecoder(bytes.NewReader(response.BodyAsBytes))
decoder.DisallowUnknownFields()
err = decoder.Decode(&currencies)
```
