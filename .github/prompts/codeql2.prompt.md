# Guidelines for Authoring CodeQL Queries for Go

---

## 1  Purpose

Create precise, efficient, and maintainable CodeQL queries that locate root-cause patterns in Go codebases. Follow the imperative steps and examples below to ensure every query is production-ready and easily extended.

---

## 2  Authoring Workflow (Execute in Order)

1. **Clarify the Pattern**

   * Summarize the bug or vulnerability in one sentence.
   * Identify *where* it appears (function, call site, loop, etc.) and *what* is missing or wrong.

2. **Select Core Entities**

   * Choose the most specific CodeQL classes that model the pattern (e.g., `CallExpr`, `LoopStmt`).
   * Determine any required type or package constraints (`hasQualifiedName`).

3. **Draft Helper Predicates**

   * Isolate each logical test into a separate predicate.
   * Name predicates with verbs that describe intent (`callsMutexLock`, `callsMutexUnlock`).

4. **Compose the Main Query**

   * Import Go library: `import go`.
   * Combine helper predicates in a single `from … where … select` block.
   * Target the *absence* of required code or the *presence* of dangerous code.

5. **Add Metadata**

   ```ql
   /**
    * @name      Clear, human-readable title
    * @description Root-cause summary and risk
    * @kind      problem | path-problem
    * @problem.severity error | warning | recommendation
    * @precision high | medium | low
    * @id        go/<category>/<short-name>
    */
   ```

6. **Test and Refine**

   * Run on representative code; confirm it finds true positives only.
   * Tighten conditions or add exclusions until false positives are minimal.

7. **Optimize Performance**

   * Filter early by package, name, or type.
   * Avoid unrestricted transitive closures (`getAChild*`) without narrowing.

8. **Document Assumptions**

   * Add inline comments for any non-obvious logic.
   * Record known limitations in the metadata description.

---

## 3  Best-Practice Checklist

| #  | Imperative Rule                                                                      |
| -- | ------------------------------------------------------------------------------------ |
| 1  | Model the **root cause**, not the symptom.                                           |
| 2  | Match by **qualified name** whenever specificity matters.                            |
| 3  | Keep **one concern per query**; split if logic diverges.                             |
| 4  | Prefer **helper predicates** over long composite conditions.                         |
| 5  | Use **data-flow libraries** for source-to-sink analysis instead of manual traversal. |
| 6  | Set **@precision high** only when the pattern cannot misfire.                        |
| 7  | Guard performance—apply narrowing filters before joins.                              |
| 8  | Validate on **both** positive and negative examples.                                 |
| 9  | Include **all** mandatory metadata fields.                                           |
| 10 | Write **clear, imperative comments** that future maintainers can follow.             |

---

## 4  Illustrative Examples

### 4.1  Missing Unlock after Mutex Lock

```ql
/**
 * @name Missing Unlock After Mutex Lock
 * @description Detects functions that call Mutex.Lock() without a corresponding Mutex.Unlock().
 * @kind problem
 * @problem.severity warning
 * @precision medium
 * @id go/concurrency/missing-unlock
 */
import go

predicate callsMutexLock(Function f) {
  exists(CallExpr c |
    c.getEnclosingFunction() = f and
    c.getCalleeExpr() instanceof SelectorExpr s and
    s.getSelector().getName() = "Lock" and
    s.getReceiverType().(NamedType).hasQualifiedName("sync", "Mutex")
  )
}

predicate callsMutexUnlock(Function f) {
  exists(CallExpr c |
    c.getEnclosingFunction() = f and
    c.getCalleeExpr() instanceof SelectorExpr s and
    s.getSelector().getName() = "Unlock" and
    s.getReceiverType().(NamedType).hasQualifiedName("sync", "Mutex")
  )
}

from Function f
where callsMutexLock(f) and not callsMutexUnlock(f)
select f, "Function `" + f.getName() + "` locks a mutex without unlocking it."
```

### 4.2  Defer Inside Loop

```ql
/**
 * @name Defer in Loop
 * @description Flags any defer statement placed within a loop body.
 * @kind problem
 * @problem.severity recommendation
 * @precision high
 * @id go/performance/defer-in-loop
 */
import go

from LoopStmt loop, DeferStmt d
where loop.getAChild*() = d
select d, "Avoid deferring inside loops; deferred calls execute only after the loop finishes."
```

---

## 5  Go CodeQL Cheat Sheet

| Category    | Key Classes / Predicates                                    | Usage Hint                                             |
| ----------- | ----------------------------------------------------------- | ------------------------------------------------------ |
| Functions   | `Function`, `Method`                                        | `hasQualifiedName(pkg, type, name)` for specificity    |
| Calls       | `CallExpr` → `getCalleeExpr()` → `SelectorExpr`             | Inspect method name and receiver type                  |
| Statements  | `IfStmt`, `ForStmt`, `RangeStmt`, `DeferStmt`, `ReturnStmt` | Access body with `getThen()`, `getElse()`, `getBody()` |
| Expressions | `BinaryExpr`, `UnaryExpr`, `Literal` subclasses             | `getLeftOperand()`, `getRightOperand()`, `getValue()`  |
| Types       | `NamedType`, `PointerType`, `StructType`, `InterfaceType`   | Check with `hasQualifiedName` or `getName()`           |
| Data Flow   | `DataFlow::Configuration`, `Node`, `PathNode`               | Define `isSource`, `isSink`, then call `hasFlowPath`   |

---

## 6  Common Pitfalls — Avoid These

1. **Broad Name Matching** – Always include package/type qualifiers.
2. **Missing Metadata** – Queries without full metadata may be skipped.
3. **Complex One-Liners** – Break logic into helper predicates for clarity.
4. **Unconstrained Wildcards** – Limit `getAChild*` and similar transitive searches.
5. **Untested Logic** – Validate each query on real positive and negative examples.
6. **Ignoring Receiver Identity** – Ensure paired operations act on the same object.
7. **Excessive False Positives** – Tighten conditions or whitelist benign cases.
8. **Performance Bottlenecks** – Filter early; avoid Cartesian joins of large sets.
9. **Copy-Paste IDs** – Maintain unique, descriptive `@id` values.
10. **Outdated API Calls** – Verify every class and predicate exists in current Go pack.
