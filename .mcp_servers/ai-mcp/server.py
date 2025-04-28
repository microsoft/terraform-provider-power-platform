from typing import Any
#import httpx
from mcp.server.fastmcp import FastMCP
from init import init_model_client, init_tools, init_magentic_one_agent
from autogen_agentchat.agents import AssistantAgent
from autogen_ext.tools.mcp import StdioServerParams, mcp_server_tools
from autogen_agentchat.conditions import MaxMessageTermination, TextMentionTermination
from autogen_agentchat.teams import RoundRobinGroupChat
import os

model_client = init_model_client()

# Initialize FastMCP server
mcp = FastMCP("ai-mcp")

@mcp.tool()
async def save_file(absolute_file_path: str, file_content: str) -> str:
    """Save the file to the specified absolute path.
    Args:
        absolute_file_path: The absolute path to the file.
    Returns:
        a string indicating the file has been saved.
        'SUCCESS' if the file was saved successfully, 'FAILURE' otherwise.
    """

    # Check if the file exists
    if os.path.exists(absolute_file_path):
        return "File already exists."

    # Create the directory if it doesn't exist
    os.makedirs(os.path.dirname(absolute_file_path), exist_ok=True)

    # Create a new file and write some content to it
    with open(absolute_file_path, 'w') as f:
        f.write(file_content)
    # Check if the file was created successfully
    if not os.path.exists(absolute_file_path):
        return "FAILURE"
    # Return success message
    return "SUCCESS"

@mcp.tool()
async def analyze_code(absolute_file_path: str) -> str:
    """Analyze the code and return the result in a markdown format..

    Args:
        absolute_file_path: The absolute path to the file.
    
    Returns:
        a markdown string describing the result of the analysis.
    """

    shell_mcp_server = StdioServerParams(command="uv", args=[
        #"-v",
        "run",
        "mcp-shell-server",
       ], 
       env={
        "ALLOW_COMMANDS": "ls,pwd,echo,cat,head,tail,find,grep,wc,touch,mkdir"
       },
       read_timeout_seconds=60)
    tools = await mcp_server_tools(shell_mcp_server)

    primary_agent = AssistantAgent(
        "CodeQualityAgent",
        model_client=model_client,
        tools=tools,
        system_message="""You are a code quality agent, keep going until the user’s query is completely resolved, before ending your turn and yielding back to the user. Only terminate your turn when you are sure that the problem is solved.
        Your job is to review the code in the file provided and for every issue found, provide feedback filling up the following markdown template for that issue:
# Title

<<Title for the issue>>

##

<<Path to the file>>

## Problem

<<Description of the issue>>

## Impact

<<Explanation how does the issue impact the codebase and what is the sevirity: low, medium, high, critical>>

## Location

<<File location where the issue was found>>

## Code Issue

<<Copy of the piece of code where the issue is

```go

```
>>

## Fix

<<Code example how to fix the issue with explanation

```go

```
>>
      """,
    )


    critic_agent = AssistantAgent(
        "critic",
        model_client=model_client,
        system_message="Respond with 'CONTINUE' until markdown content with analysis is provided. When markdown contenat is provided respond with 'APPROVE' to approve the analysis.",
    )

    max_msg_termination = MaxMessageTermination(max_messages=10)
    text_termination = TextMentionTermination("APPROVE")
    combined_termination = max_msg_termination | text_termination

    round_robin_team = RoundRobinGroupChat([primary_agent, critic_agent], termination_condition=combined_termination)

    prompt = f"""Analyze code under following absolute path: {absolute_file_path}"""
    
    result = await round_robin_team.run(task=prompt, cancellation_token=None)
    return result.messages[-2].content

if __name__ == "__main__":
    # Initialize and run the server
    mcp.run(transport='stdio')
