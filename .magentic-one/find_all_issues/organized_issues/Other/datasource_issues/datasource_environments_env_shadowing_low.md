# Type Naming: Variable `env` Shadowing in Loop

##

/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_environments.go

## Problem

In the `Read` function, the loop variable is named `env` (`for _, env := range envs`), and then inside the loop, another variable is also named `env` from the result of `convertSourceModelFromEnvironmentDto`. This variable shadowing is confusing and can lead to subtle bugs or misunderstanding when refactoring or debugging.

## Impact

Reduces code readability and maintainability; makes debugging and further changes riskier due to unclear reference. **Severity: Low**.

## Location

```go
for _, env := range envs {
    ...
    env, err := convertSourceModelFromEnvironmentDto(env, ...)
    ...
    state.Environments = append(state.Environments, *env)
}
```

## Code Issue

```go
for _, env := range envs {
    ...
    env, err := convertSourceModelFromEnvironmentDto(env, ...)
    ...
    state.Environments = append(state.Environments, *env)
}
```

## Fix

Use distinct names for the loop variable and the converted value. For example:

```go
for _, env := range envs {
    ...
    convertedEnv, err := convertSourceModelFromEnvironmentDto(env, ...)
    if err != nil {
        // ...
    }
    state.Environments = append(state.Environments, *convertedEnv)
}
```

This makes the code clearer and avoids accidental reference errors.
