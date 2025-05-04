/**
 * @name Function (or method) call returning an error
 * @description Flags any function or method invocation whose result list
 *              includes the predeclared `error` interface.
 * @kind problem
 * @id go/call-returning-error
 * @problem.severity recommendation
 * @tags reliability
 *       error-handling
 *       go
 */
import go
import semmle.go.Types   // supplies ErrorType

/** A `return` that contains a bare `error`-typed identifier. */
class BareErrorReturn extends ReturnStmt {
  BareErrorReturn() {
    exists(int idx, Ident id |
      // one of the returned expressions is the identifier itself
      this.getExpr(idx) = id and
      // the identifier’s declaration has (or implements) the `error` type
      id.getType() instanceof ErrorType
    )
  }
}

from BareErrorReturn ret
select ret,
  "Error value is returned directly; consider wrapping it with fmt.Errorf(" +
  "\"<context>: %w\", err) or similar for better diagnostics."

//  import go
//  import semmle.go.Types   // brings in SignatureType and ErrorType
 
//  /** A call-expression whose signature’s result types include `error`. */
//  class CallReturningError extends CallExpr {
//    CallReturningError() {
//      exists(int i |
//        i < this.getCalleeType().getNumResult() and
//        this.getCalleeType().getResultType(i) instanceof ErrorType
//      )
//    }
//  }
 
// /** A `return` that contains such a call literally. */
// predicate isReturnedDirectly(CallReturningError call, ReturnStmt ret) {
//     ret.getAnExpr() = call   // one of the return expressions *is* the call node
//   }
  
//   from CallReturningError call, ReturnStmt ret
//   where isReturnedDirectly(call, ret)
//   select ret,
//     "Call that returns an `error` is returned directly; consider adding context or checks."
