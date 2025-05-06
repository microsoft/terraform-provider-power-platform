/**
 * @name Detect panic calls
 * @description Detects direct calls to the built-in `panic` function. In Terraform providers, `panic` should be avoided as it leads to an abrupt termination of the plugin, bypassing graceful error handling and resource cleanup. Errors should be returned instead to allow Terraform Core to manage the lifecycle and user feedback appropriately.
 * @kind problem
 * @problem.severity warning
 * @precision high
 * @id go/terraform-provider/panic-issue-detection
 */

 import go

 from CallExpr panicCall
 where
  panicCall.getTarget().getQualifiedName() = "panic"

 select panicCall, "Avoid using panic for error handling; return an error instead. Panics can abruptly terminate the provider."
