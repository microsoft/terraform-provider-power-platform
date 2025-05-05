
# CRITIC_AGENT_PROMPT = """Respond with 'CONTINUE' until markdown content with analysis is provided for the given file. When markdown content is provided  and saved in files, respond with 'APPROVE' to approve the analysis. 
# Do not allow to analyze additional files, or execute other actions. Start responding with 'APPROVE' if the analysis is correct.
# Multiple markdown files will be created based on the severity of the issues found in the code. Do not respond with 'APPROVE' until all markdown files are created and saved.
# """

CRITIC_AGENT_PROMPT = """Respond with 'CONTINUE' until unittest and linter is executed successfully. when all tests are passed and linter is fixed, respond with 'APPROVE' to approve the analysis.
Do not run any other tools, or execute other actions. Start responding with 'APPROVE' when the task is completed.
"""

def prep_code_agent_prompt(issue_file_path: str, respository_folder: str, instructions_file_path: str) -> str:
    """
    Prepares the prompt for the code agent based on the file path and folder path.
    """

    return f"""# Role and Objective
You are a code agent responsible for fixing code issues described in {issue_file_path}
You are responsible for running validation tools for the repository focated in '{respository_folder}'.
keep going until the user query is completely resolved, before ending your turn and yielding back to the user. Only terminate your turn when task is completed and you are sure that the problem is solved.

To create changelog entry, you need to run the changie tool.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed (see [copilot-commit-message-instructions.md](../copilot-commit-message-instructions.md))
- `<issue_number>` pick a number given in instructions.
"""
