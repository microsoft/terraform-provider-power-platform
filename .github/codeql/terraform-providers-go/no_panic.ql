/**
 * @name Detect panic calls for API rest scope or environment_id issues
 * @description Finds instances where panic is called with messages indicating missing scope or environment_id.
 * @kind problem
 * @problem.severity warning
 * @id go/terraform-provider/panic-issue-detection
 * @language go
 */

 import go

 from CallExpr panicCall
 where
   panicCall.getTarget().getName() = "panic" 
 
 select panicCall, "Avoid using panic for error handling; return an error instead."
 