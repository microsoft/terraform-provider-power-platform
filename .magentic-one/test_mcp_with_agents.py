import asyncio

from autogen_agentchat.agents import AssistantAgent
from autogen_ext.models.openai import AzureOpenAIChatCompletionClient
from azure.identity import DefaultAzureCredential
from autogen_ext.auth.azure import AzureTokenProvider
from autogen_ext.tools.mcp import StdioServerParams, mcp_server_tools
from autogen_agentchat.conditions import MaxMessageTermination, TextMentionTermination
from autogen_agentchat.teams import RoundRobinGroupChat

async def main() -> None:
    # Get the fetch tool from mcp-server-fetch.
    fetch_mcp_server = StdioServerParams(command="uv", args=["run","-v","mcp-server-fetch"])
    tools = await mcp_server_tools(fetch_mcp_server)

    token_provider = AzureTokenProvider(
        DefaultAzureCredential(),
        "https://cognitiveservices.azure.com/.default",
    )

    # Create an agent that can use the fetch tool.
    model_client = AzureOpenAIChatCompletionClient(
        azure_deployment="gpt-4o",  # Replace with your Azure deployment name
        model="gpt-4o",  # Replace with your model name
        api_version="2024-12-01-preview",
        azure_endpoint="https://<<your_az_foundry_endpoin>>.openai.azure.com/",  # Replace with your endpoint
        azure_ad_token_provider=token_provider,  # Optional if you choose key-based authentication
    )
    # agent = AssistantAgent(name="fetcher", model_client=model_client, tools=tools, reflect_on_tool_use=True)  # type: ignore
    # # Let the agent fetch the content of a URL and summarize it.
    # result = await agent.run(task="Summarize the content of https://en.wikipedia.org/wiki/Seattle")
    # print(result.messages[-1].content)

    shell_mcp_server = StdioServerParams(command="uv", args=[
        #"-v",
        "run",
        "mcp-shell-server",
       ], 
       env={
        "ALLOW_COMMANDS": "ls,pwd,echo,cat,head,tail,find,grep,wc"
       },
       read_timeout_seconds=60)
    tools = await mcp_server_tools(shell_mcp_server)

    primary_agent = AssistantAgent(
        "CodeQualityAgent",
        model_client=model_client,
        tools=tools,
        system_message="""You are a code quality agent. Your job is to review the code and provide feedback filling up the following markdown template:
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

    prompt = f"""Analyze code under following absolute path: /workspaces/terraform-provider-power-platform/internal/api/auth.go"""
    
    result = await round_robin_team.run(task=prompt, cancellation_token=None)
    print(result.messages[-2].content)

if __name__ == "__main__":
    asyncio.run(main())
