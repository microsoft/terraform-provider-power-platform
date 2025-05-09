# CodeQL Query Guidelines for Go (GitHub Copilot Agent Mode)

## Authoring Best Practices

* **Understand the root cause:** Always begin by clarifying the bug or vulnerability at a fundamental level. Identify the **root cause pattern** in code (for example, a missing check, an API misuse, a dangerous data flow) rather than focusing on superficial symptoms. This ensures your query finds the real issue, not just one instance of it.

* **Leverage CodeQL’s Go library:** Always start your query with `import go` to pull in the standard Go CodeQL classes. Use the provided classes (like `Expr`, `CallExpr`, `Function`, etc.) and predicates instead of writing low-level logic from scratch. These abstractions make your query more expressive and robust. For example, use `CallExpr.getArg(i)` to get function call arguments or `Method.hasQualifiedName("pkg", "Type", "Name")` to precisely match a known function/method – this is clearer and less error-prone than manual string matching.

* **Be precise and minimize false positives:** Define enough conditions to pinpoint the problematic pattern **without accidentally matching safe code**. Add constraints for types, names, values, or surrounding context as needed. For instance, if searching for a dangerous function call, ensure it’s the right one by matching its qualified name and perhaps the types of its arguments. If looking for missing cleanup (like opened resources not closed), ensure the resource was actually opened first. Each extra filter (an `and` clause, a type check, etc.) should rule out false alarms and hone in on real issues.

* **Focus on what “should” be true or false:** Write your query conditions to capture the bug’s *presence* or the *absence of a required element*. For example, to find unchecked errors, require that an `IfStmt` checking an error is **not** present when it should be. This approach targets the root cause (the missing check) instead of the symptom (some downstream crash). Always ask: *“What code *should* accompany this pattern?”* and then make the query ensure that expected code is **missing** or that an unexpected code pattern is **present**.

* **Use helper predicates for clarity:** Break complex logic into smaller **predicates** or classes within your query. This improves readability and maintainability. For example, define `predicate callsDangerousFunc(Function f)` to encapsulate the logic of “function `f` calls dangerous function X” and another predicate `hasNoValidation(f)` for “f does not perform a required check.” Then your main query can simply combine these predicates. This makes it easier to adjust or reuse parts of the logic later, and it helps the Copilot agent (and humans) understand the intent of each part.

* **Write idiomatic CodeQL:** Follow CodeQL query conventions and style. Use meaningful class and predicate names that describe their intent. Prefer **declarative** patterns (like existence of an AST node with certain properties) over complex procedural logic. Let the query engine handle iteration and flow—your job is to specify *what* to find. Using the built-in classes and their member predicates (like `getCondition()` for an `IfStmt` or `getLeftOperand()` for a `BinaryExpr`) results in more idiomatic and efficient queries.

* **Include complete query metadata:** Every query should start with a comment block providing metadata for CodeQL. Include a descriptive `@name` (what the query finds), a helpful `@description` (why it matters or an example scenario), and appropriate tags like `@kind` (e.g. `"problem"` or `"path-problem"`), `@problem.severity` (e.g. `"warning"` or `"error"` or `"recommendation"`), `@precision` (`high`, `medium`, or `low`), and any relevant security tags (like `external/cwe-089` for SQL injection, etc.). This information is crucial for integrating the query into larger query suites and for anyone reviewing the results. **Do not skip metadata.**

* **Aim for High Precision:** Especially for security queries, prefer **precision over recall** in most cases. It’s often better to miss a fringe case than flag numerous false positives in developers’ code. You can achieve this by adding necessary conditions (even at the risk of being slightly narrower). If needed, mention in the query description what limitations or assumptions your query makes. In Copilot Agent context, it’s easier to explain a narrower query than to overwhelm a user with many false alarms. That said, ensure the logic is not so tight that it misses the primary instances of the issue.

* **Optimize iteratively:** Write a first version of the query that is correct and clearly expresses the concept, then consider performance. Check if you are iterating over a huge set unnecessarily. Common performance tips include: constrain searches by using specific classes (for example, iterating over `Function`s is cheaper than iterating over all `AstNode`s), use predicates like `hasQualifiedName` to quickly filter by name, and avoid deep nested loops with no filters. If a query could return a very large number of results, consider adding a `limit` or additional guard conditions. The CodeQL engine is powerful, but giving it hints (through more specific where-clauses) can significantly speed up queries on big codebases. In an agent scenario, you want the query to run efficiently even on large repositories.

* **Use data flow libraries for taint issues:** If the pattern involves data flowing from a source to a sink (e.g., user input flowing into a dangerous function), use CodeQL’s data flow or taint tracking libraries instead of manually coding the flow logic. For Go, import `DataFlow` or `TaintTracking` from the standard library and configure sources and sinks via a `Configuration` class. This high-level approach is less error-prone and more powerful, as it accounts for intermediate assignments, function calls, and other flow complexity automatically. Aim to detect the **whole path** of tainted data to the sink – this yields a path-problem query that can show step-by-step how the vulnerability occurs. (Remember to set `@kind path-problem` for those and use `select path, ...` to include path information.)

* **Keep the query maintainable:** Write the query as if others will read or extend it. Use comments (outside the QL metadata block) to explain non-obvious logic if necessary. Ensure the query is logically organized (for example, list your `where` conditions in a sensible order, maybe group related ones together). If you find the query becoming too complex, consider splitting it or using multiple predicates for clarity. Maintaining a clean structure will also help the Copilot agent if it needs to modify or debug the query later.

By following these best practices, your CodeQL queries will be **correct, efficient, and easy to understand** – whether written by you or generated with the help of Copilot.

## Query Composition Examples

To illustrate the process of translating a bug description into a CodeQL query, here are two example scenarios with full CodeQL queries and explanations.

### Example 1: Missing Unlock After Mutex Lock

**Scenario:** We suspect that in some Go functions, a mutex is locked but never unlocked, leading to potential deadlocks. The bug pattern is: a call to `mu.Lock()` occurs without a corresponding `mu.Unlock()` call in the same function.

**Approach:** To detect this, we will find functions that call a `Lock` method on a `sync.Mutex` (or `*sync.Mutex`) and do *not* call the `Unlock` method. We use CodeQL’s ability to examine method call names and types. We’ll structure the query with helper predicates for clarity.

```ql
/**
 * @name Missing Unlock After Mutex Lock
 * @description Flags functions that call Mutex.Lock() without a corresponding Mutex.Unlock(), which can lead to deadlock.
 * @kind problem
 * @problem.severity warning
 * @precision medium
 * @id go/concurrency/missing-unlock
 */
import go

// Predicate to check if a function calls sync.Mutex.Lock
predicate callsMutexLock(Function func) {
  exists(CallExpr call |
    call.getEnclosingFunction() = func and
    call.getCalleeExpr() instanceof SelectorExpr sel and
    sel.getSelector().getName() = "Lock" and
    sel.getReceiverType().(NamedType).hasQualifiedName("sync", "Mutex")
  )
}

// Predicate to check if a function calls sync.Mutex.Unlock
predicate callsMutexUnlock(Function func) {
  exists(CallExpr call |
    call.getEnclosingFunction() = func and
    call.getCalleeExpr() instanceof SelectorExpr sel and
    sel.getSelector().getName() = "Unlock" and
    sel.getReceiverType().(NamedType).hasQualifiedName("sync", "Mutex")
  )
}

from Function func
where callsMutexLock(func) and not callsMutexUnlock(func)
select func, "Function `${func.getName()}` locks a mutex without unlocking it."
```

**Explanation:** This query imports the Go QL library and defines two predicates, `callsMutexLock` and `callsMutexUnlock`, which are true if a given function contains at least one call to `Mutex.Lock()` or `Mutex.Unlock()` respectively. We identify those calls by looking for `CallExpr` in the function where the callee is a `SelectorExpr` (meaning something like `object.Method`) and the method name is `"Lock"` or `"Unlock"`. We further ensure the receiver type of that selector is the `sync.Mutex` type (using `hasQualifiedName("sync", "Mutex")` on the type of the receiver). Finally, the query selects all functions that call `Lock` but do not call `Unlock`. The result reported is the function itself, with a message indicating the issue. This query focuses on the **absence** of an unlock call, which is the root cause of the deadlock bug, and avoids false positives by checking the method names and types precisely.

### Example 2: Defer Call Inside Loop

**Scenario:** Using the `defer` statement inside loops is a known performance issue in Go. Each iteration will add a deferred function to the stack, which won’t execute until the function returns, potentially exhausting resources if the loop runs many times. We want to detect any `defer` inside a loop.

**Approach:** We will find any `DeferStmt` (defer statement) that is lexically inside a loop (`for` or range loop). The CodeQL library has a class for loop statements and a class for defer statements, which we can directly use to check the parent-child relationship in the AST.

```ql
/**
 * @name Defer in Loop
 * @description Detects a 'defer' statement inside a loop, which can cause increased memory usage and delayed execution of deferred calls.
 * @kind problem
 * @problem.severity recommendation
 * @precision high
 * @id go/performance/defer-in-loop
 */
import go

from LoopStmt loop, DeferStmt deferStmt
where loop.getAChild*() = deferStmt
select deferStmt, "Avoid deferring a function call inside a loop; the deferred call will not run until the loop finishes, leading to potential performance issues."
```

**Explanation:** This query uses `LoopStmt`, a class representing any loop (`for`, including `range` loops), and `DeferStmt` for any deferred function call. The `where` clause uses `loop.getAChild*() = deferStmt` to ensure the `DeferStmt` is a descendant of the loop in the abstract syntax tree (the `*` in `getAChild*` allows matching at any depth within the loop’s body, so even if the defer is inside an if-statement or another block inside the loop, it will be caught). We mark the result at the location of the `defer` statement and output a message advising against deferring inside loops. The query uses a high precision because any defer inside a loop is almost certainly unintentional or ill-advised – thus we can be confident in flagging it.

These examples demonstrate taking a textual pattern description and turning it into CodeQL conditions. The first example targets a **resource management bug** by correlating two related function calls, and the second targets a **performance concern** by structural context (a statement inside a loop). In both cases, the queries include metadata and are written in an imperative, clear style suitable for automation via Copilot. The use of helper predicates and precise matching (qualified names, specific AST node types) helps ensure the queries are accurate and maintainable.

## Copilot Prompting Tips

Writing CodeQL queries with the help of GitHub Copilot (in Agent/Chat mode) can boost productivity. Here are some tips to craft prompts that guide the AI effectively:

* **Provide clear context and requirements:** Start your conversation by clearly describing the bug or code smell you want to detect. Include relevant details like function names, code snippets, or specific conditions if possible. For example: *“I have a bug where `os.Exit` is called in library code, which is undesirable. I need a CodeQL query to find any call to `os.Exit` in package code.”* Being specific about what you want ensures the agent generates a focused query.

* **Specify the output format upfront:** Tell Copilot you want a CodeQL query and remind it to include the necessary sections. For example: *“Write a CodeQL query (in Go QL) with proper metadata comments that…”*. By stating this, the agent will know to produce a full query with the `/** ... */` metadata block, the `import go`, and the query logic. You can even list the metadata fields you expect (name, description, etc.) to make sure none are missed.

* **Use an iterative approach for complex queries:** Break down the problem into steps in your prompts. You might first ask, *“Show me how to get all calls to function X in Go using CodeQL.”* Once the agent provides that, you can refine: *“Now modify the query to filter those calls that are missing a preceding check Y.”* This step-by-step prompting aligns with how a human would develop the query and helps the agent correct course incrementally. Copilot can keep the context of earlier steps, so it will combine them in the refinement.

* **Leverage few-shot learning:** If you have examples of similar queries, provide them to Copilot as part of the prompt. For instance, you could show a small query that finds a similar pattern (perhaps in a different language or a simpler scenario) and then ask Copilot to write a new query following that style for your Go scenario. Example prompt: *“Here is a CodeQL query that finds hard-coded credentials in Java:\n`ql\n(…Java query example…)\n`\nWrite a similar query in Go that finds hard-coded AWS secret keys.”* This helps the agent mimic the structure and practices from the example.

* **Ask for explanations or verify understanding:** After Copilot produces a query, you can ask it to explain the query to ensure it understood the task. For example: *“Explain how this query works.”* Reviewing the explanation can reveal misunderstandings or mistakes in the logic, which you can then address with further prompts. If something is off, instruct the agent with a clear directive, e.g., *“The query missed filtering by package name; please incorporate a check that the function belongs to package X.”*

* **Iterate and refine:** It’s rare to get a perfect query in one shot. Use Copilot Chat to iteratively refine the query. You might prompt, *“The query is returning some false positives where the pattern is actually safe. How can we refine the where clause to exclude those?”* By iterating, you guide the agent to adjust the logic, add conditions, or improve performance (perhaps by adding an index or using a different approach).

* **Keep prompts and code separate:** When drafting a query with Copilot, it can be useful to use comments or pseudo-code in the prompt as scaffolding. For example: *“Let’s write a query. First, find all `sql.Query` calls. Then filter those whose first argument is a string concatenation. We’ll need to use `BinaryExpr` to detect string concatenation. Finally, select those calls. Now provide the complete query.”* Structuring the prompt like this (almost like writing pseudocode for the query) can help Copilot follow the intended structure.

* **Validate the output:** Once Copilot provides a query, **review it carefully**. Even with good prompting, the agent might produce a logically flawed or suboptimal query. Check that all metadata is filled in correctly, the logic matches the intended pattern, and there are no obvious performance issues (like an unbounded loop over all AST nodes). If you have a small test code snippet, run the query (if possible) on it to see if the results make sense. If the query needs tuning, explain what needs to change to Copilot and let it suggest improvements.

By using these prompting strategies, you can effectively guide Copilot to generate high-quality CodeQL queries. Treat Copilot as a capable assistant: give it clear instructions, verify its work, and progressively steer it toward the correct and optimized solution.

## Go CodeQL Cheat Sheet

When writing CodeQL for Go, you’ll frequently use certain classes and predicates from the Go standard library. Below is a handy cheat sheet of common CodeQL classes and how to use them:

* **`Function` / `Method`:** Represents a function or method definition in the code. Use `Function.getName()` to get its name. You can get parameters or return variables via `getParameter(index)` or `getResultVar(index)`. A `Method` is a specialized function with a receiver (for example, a method on a struct or interface). For methods, you can use `Method.hasQualifiedName("pkg", "Type", "methodName")` to identify a specific method by package, receiver type, and name. Both Function and Method allow `getACall()` to find call sites, or you can iterate through `CallExpr` (see below) to find where they’re invoked.

* **`CallExpr`:** An expression representing a function or method call. `CallExpr.getCalleeExpr()` gives the expression for what’s being called (for instance, an `Ident` for a simple function call, or a `SelectorExpr` for a method call like `obj.foo`). You can get call arguments with `CallExpr.getArg(i)` for the i-th argument. Often you will cast the callee to a `SelectorExpr` to get details on method calls (e.g., `call.getCalleeExpr().(SelectorExpr).getSelector().getName()` for the method name, or get the base object via `getBase()`). Use `CallExpr` to find where functions are used and to inspect their arguments or receivers.

* **`SelectorExpr`:** An expression of the form `X.Y` (like accessing a field or calling a method Y on object X). `SelectorExpr.getBase()` returns the expression for `X` and `getSelector()` returns the `Ident` for `Y`. You can use this to distinguish method calls vs field accesses. For example, if you have `sel = call.getCalleeExpr().(SelectorExpr)`, then `sel.getSelector().getName()` might be `"Lock"` and `sel.getBase().getType()` would give you the type of the object on which the method is called.

* **`Ident`:** An identifier in the code (a name). `Ident.getName()` returns the string name. Idents can refer to variables, functions, types, etc. Often you will cast an expression to `Ident` to check if it’s a reference to a particular name (but remember that just matching a name might not be unique — scope and qualification matter).

* **Statement classes:** Various subclasses of `Stmt` represent Go statements:

  * **`IfStmt`:** An `if` statement. Use `IfStmt.getCond()` to get the condition expression, and `getThen()` / `getElse()` to get the consequent and alternative statements (the bodies). You can further cast those to `BlockStmt` if needed to get inside them. For example, `ifStmt.getThen()` gives you the `Stmt` that is executed on true (often a `BlockStmt` for a braced block).
  * **`ForStmt` / `RangeStmt`:** Loop statements. `ForStmt` represents the traditional `for` with an initialization, condition, and post statement (or a simplified form like `for condition { }` or infinite loop). `RangeStmt` represents `for ... range` loops. They both subclass `LoopStmt`. You can get the loop’s body via `LoopStmt.getBody()` which returns the `BlockStmt` inside the loop.
  * **`ReturnStmt`:** A return statement. `ReturnStmt.getResult(i)` gives the i-th expression being returned (if multiple returns).
  * **`DeferStmt`:** A defer statement (calling a function to be deferred). Use `DeferStmt.getCall()` to get the deferred call (as a `CallExpr`).
  * **`GoStmt`:** A go statement (starting a new goroutine). Similar to DeferStmt, it has `getCall()` for the call that is invoked concurrently.

* **Expression classes:** Subclasses of `Expr` for various kinds of expressions:

  * **`BinaryExpr`:** A binary operation like `x + y` or `a && b`. You can get the left and right sides with `getLeftOperand()` and `getRightOperand()`. For comparisons, there are subclasses like `EqualityTestExpr` (for `==` and `!=`) and `RelationalComparisonExpr` (for `<`, `>`, etc.), which have extra predicates like `getPolarity()` for equality or `isStrict()` for `<` vs `<=`.
  * **`UnaryExpr`:** A unary operation like `!flag` or `*ptr`. Use `getOperand()` to get the expression it applies to, and `getOperator()` to distinguish `!`, `*`, `&`, etc.
  * **`CallExpr`:** (described above, but remember it’s also an Expr subclass).
  * **`Literal` classes:** Constants in code. For example, `IntLit`, `StringLit`, `BoolLit` for integer, string, boolean literals. Each has a `getValue()` predicate to get the actual constant value (as a string or number).
  * **`CompositeLit`:** A composite literal (like struct literal or array literal). You can get its elements or key/value pairs if it’s a map or struct literal.
  * **`FuncLit`:** A function literal (anonymous function). You can get its `getBody()` (block of statements) and note it doesn’t have a name.
  * **`IndexExpr` / `SliceExpr`:** Indexing or slicing expressions like `arr[i]` or `slice[lo:hi]`. Use `IndexExpr.getBase()` / `getIndex()`, and for `SliceExpr` use `getBase()`, `getLow()`, `getHigh()` etc.

* **Type-related classes:** CodeQL represents types as well:

  * **`Type` (and subclasses):** Represents a type in Go. For example, `StructType` for struct definitions, `InterfaceType` for interface definitions, `PointerType` for pointers, etc. These classes let you query type information. A `StructType` can have members (fields and methods). An `InterfaceType` has methods (representing the interface’s method set). You might use `Type.getUnqualifiedName()` to get the name of a type and `Type.getPackage()` to get the package it’s defined in.
  * **`Field`:** Represents a struct field **entity** (as opposed to a field access expression). You can check a field’s name and type. Use `Field.hasQualifiedName("pkg", "StructType", "fieldName")` to identify a specific struct’s field. This is useful if you’re looking for usage of a particular important field (like a sensitive field).
  * **`Variable` and `Parameter`:** Represents local or global variables and function parameters. A `Parameter` object (which may be a subclass of `Variable`) can give you the function it belongs to and its index. You might use these to track data flow from parameters or to ensure certain parameters have specific properties (like checking if a parameter name starts with `_` to find unused parameters, etc.).

* **Data flow and path classes:**

  * **`DataFlow::Node`:** An abstraction of an expression or value in the program used for data flow analysis. You get a `DataFlow::Node` from an `Expr` by `expr.asExpr().getASuccessor()` or more commonly by using a library predicate to find source or sink nodes. If you are doing a custom taint tracking, you’ll define a `Configuration` and use its `hasFlow()` or `hasFlowPath()` predicates with `DataFlow::Node` or `DataFlow::PathNode`.
  * **`DataFlow::Configuration`:** The base class to configure a data flow analysis (or use `TaintTracking::Configuration` for convenience if dealing with taint-style source-to-sink). You override `isSource(Node)` and `isSink(Node)` to define where data comes from and where it should (not) go. There are often utility modules like `import go::DataFlow` that give you default configurations for common scenarios (user input to output, etc.), but you can always make a custom one for your specific pattern.
  * **`PathNode` / `PathProblem`:** When writing path-problem queries (to show how data flows or how a series of statements lead to a bug), CodeQL provides classes to represent nodes along a path. You typically use `hasFlowPath(source, sink, path)` or simply select `path` in a query with `@kind path-problem` to output the trace.

* **Locations and elements:**

  * **`AstNode.getLocation()`**: Any AST node (statements, expressions, declarations) can give you a `Location` which contains file, line, column info. Usually you don’t need to use this explicitly in a query (the CodeQL engine will report the location of whatever element you select), but it’s available if needed.
  * **`Element` and `Entity`:** `Element` is a very generic class that many CodeQL classes extend (like almost everything in a CodeQL database). `Entity` is a subclass for program elements that have an identity beyond the AST (like functions, types, variables as opposed to a syntactic construct). Typically you don’t use these directly, but it’s good to know that e.g. `Function` and `Field` are Entities, meaning they have a `getName()` and other common traits.

**Usage hints:** When writing a query, choose the most specific class that fits what you need. For example, if you need to find function declarations, use `FuncDecl` (a specific AST node for a top-level function) or `FuncDef` (covers both functions and function literals) or simply `Function` (if you want the entity and its properties), depending on context. If you need to ensure something is a method call, use the presence of a `SelectorExpr`. If you want to match by name and package, prefer `hasQualifiedName` on an Entity like `Function`, `Method`, or `Field` because it’s precise (it matches exactly a package/path, type, and name). Use `getName()` string comparisons only when package/type context isn’t important or available.

This cheat sheet can be a quick reference as you formulate queries or instruct Copilot. By knowing these common classes and their key predicates, you can more quickly tell the agent what to do (e.g., “iterate over all `IfStmt` and check if `getCond()` is a `CallExpr`” or “find any `Field` named ‘Password’”). It also helps in reading and understanding the queries Copilot produces.

## Common Pitfalls and Anti-Patterns

Even with clear guidelines, there are some common mistakes to watch out for when authoring CodeQL queries for Go. Avoiding these will save time and make your queries more effective:

* **Overly broad patterns:** A query that is too general will return many false positives. For example, searching for all calls to `Println` might flag legitimate debug logging as an issue if you intended to find sensitive info leaks. Always refine your pattern with additional conditions. If you find your query results include a lot of noise, identify what differentiates the true positives and incorporate that into the query. **Anti-pattern:** `from Ident i where i.getName() = "Foo" select i` – this finds any identifier named "Foo" (way too broad). Instead, qualify it: if “Foo” is a function in a certain package, match that specifically via a `Function` or `CallExpr` with `hasQualifiedName`.

* **Missing or incorrect metadata:** Forgetting to update the query’s metadata comments is a pitfall especially when generating queries quickly. An inaccurate `@description` or an overly generic `@name` can confuse users later, and missing the proper `@tags` (like CWE IDs for security queries) can prevent the query from being categorized correctly. Always double-check that the metadata reflects the query’s intent and use-case. **Anti-pattern:** Copy-pasting a query and not changing the `@id` or description, leading to duplicates or misleading info. Make sure each query has a unique, meaningful ID (typically a path like `go/<category>/<query-name>`).

* **Not using qualified names for symbols:** As mentioned, using `getName()` alone can be dangerous if that name exists in multiple contexts. A common pitfall is to match a function name without specifying the package or receiver type, which might catch unrelated code. Always consider if you should use `hasQualifiedName` or include the package in your logic. This ensures you don’t flag, say, `config.Load()` in *any* package when you only meant `myapp/config.Load` in your application’s own package.

* **Ignoring object identity in method patterns:** When dealing with methods like the Lock/Unlock example, it’s easy to forget to ensure you’re talking about the same object instance. Our simple example looked only within a function, which is usually sufficient. But if you ever need to ensure a method pair (like open/close, lock/unlock) operate on the same object, make sure to capture the receiver or object and use it in both parts of the pattern. Failing to do so can cause false positives or misses. **Anti-pattern:** Selecting a function if it calls `a.Lock()` and somewhere calls `b.Unlock()` – that technically has both calls but not on the same object. Guard against this by comparing the object expressions (`a` vs `b` in this case) in the query conditions.

* **Too much logic in one place:** Writing one giant `where` clause with many conditions and nested exist/non-exists can become hard to read and maintain. It’s often better to split logic into predicates (as shown in the examples). A pitfall is thinking that more nesting in one query = efficiency. In fact, the CodeQL engine will optimize across predicates as well; using predicates for clarity usually does not hurt performance and can even help the engine understand your intent. So avoid the anti-pattern of an overly complex monolithic query. Instead, create named helper predicates for different aspects of the pattern (e.g., “isTainted(source, sink)” or “locksWithoutUnlock(func)”) and then combine them. This also makes it easier to test parts of your query in isolation during development.

* **Not considering performance on large codebases:** A query might work fine on a small example but time out on a big repository. Common performance pitfalls include: iterating over all `Expr` or all `Stmt` in a program without any early filtering, or doing a join between two huge sets. To avoid this, apply filters as early as possible. For instance, if you only care about functions in a certain package, add `Function f where f.getPackage().getPath() = "mypkg"` up front. If looking for calls to a specific function, use its name or import path in the condition rather than checking every call and filtering later. Also, be mindful of recursion or transitive closures (`getAChild*` can be expensive if unconstrained). If your query is inherently heavy (like a data flow search), ensure you’ve set up sources and sinks narrowly (e.g., only tainting specific input sources, not *every* string in the program, unless necessary). In summary, think about how the query will scale: a good practice is to mentally estimate the size of each predicate’s result set and make sure you’re not joining two enormous sets without at least one selective condition.

* **Overlooking harmless patterns (false positives):** Sometimes different coding patterns can look similar to a bug pattern. For example, a developer might intentionally skip an error check in a situation where it’s safe to ignore the error. If your query flags that, it’s a false positive. To minimize these, consider common benign cases. For instance, if looking for unused error returns, you might ignore calls to `println(err)` or `_ = err` (explicitly ignoring error) as those indicate the developer consciously handled/ignored it. Likewise, if searching for dangerous function usage, consider whitelisting certain files or packages (maybe test files or generated code) if the pattern is known to appear there harmlessly. You can implement this with `not` conditions or additional checks in the `where` clause. **Anti-pattern:** Blanketly reporting something as an issue without context. Try to encode context where possible, or at least document in the query what is a potential false positive so a user can understand it.

* **Not testing the query on real code:** After writing a query, especially with Copilot’s help, failing to test it is a pitfall. You should run it (if you have a CodeQL setup) on at least a portion of the target code or similar open-source code to see if it catches what it should and nothing more. In agent mode, if you can’t run it directly, at least mentally walk through an example. Imagine a snippet of code that should be caught – does each condition in your query hold true for that snippet? Now imagine a snippet that should *not* be caught – is there any way your query would mistakenly match it? This thought experiment can reveal logical mistakes. Copilot might introduce subtle errors, like using the wrong CodeQL class or missing a condition, and testing is how you catch them.

* **Ignoring Copilot’s mistakes or hallucinations:** When using an AI assistant, sometimes it might use a class or predicate that doesn’t exist or is outdated. For example, it might try to use a wrong import or assume a predicate from a different language’s library. Always verify that the classes and predicates in the generated query actually exist in the Go CodeQL libraries. If something looks unfamiliar, check the CodeQL docs or codebase for it. Don’t assume everything Copilot writes is correct. A common anti-pattern is to trust an incorrect suggestion, which leads to a query that won’t run or will silently do the wrong thing. Stay vigilant and cross-check questionable parts of the query.
