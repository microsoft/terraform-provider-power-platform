import asyncio

from autogen_agentchat.agents import AssistantAgent
from autogen_ext.models.openai import AzureOpenAIChatCompletionClient
from azure.identity import DefaultAzureCredential
from autogen_ext.auth.azure import AzureTokenProvider
from autogen_ext.tools.mcp import StdioServerParams, mcp_server_tools


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
        azure_deployment="gpt-4.5-preview",  # Replace with your Azure deployment name
        model="gpt-4.5-preview",  # Replace with your model name
        api_version="2024-12-01-preview",
        azure_endpoint="https://mawasileazureopenai.openai.azure.com/",  # Replace with your endpoint
        azure_ad_token_provider=token_provider,  # Optional if you choose key-based authentication
    )
    agent = AssistantAgent(name="fetcher", model_client=model_client, tools=tools, reflect_on_tool_use=True)  # type: ignore

    # Let the agent fetch the content of a URL and summarize it.
    result = await agent.run(task="Summarize the content of https://en.wikipedia.org/wiki/Seattle")
    print(result.messages[-1])


if __name__ == "__main__":
    asyncio.run(main())
