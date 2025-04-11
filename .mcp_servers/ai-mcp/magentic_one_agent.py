import asyncio
import os
from autogen_ext.auth.azure import AzureTokenProvider
from autogen_ext.models.openai import AzureOpenAIChatCompletionClient
from azure.identity import DefaultAzureCredential
from autogen_agentchat.teams import MagenticOneGroupChat
from autogen_agentchat.ui import Console
#from autogen_ext.agents.web_surfer import MultimodalWebSurfer
from autogen_ext.agents.file_surfer import FileSurfer
from autogen_ext.agents.magentic_one import MagenticOneCoderAgent
from autogen_agentchat.agents import CodeExecutorAgent, AssistantAgent
from autogen_ext.code_executors.local import LocalCommandLineCodeExecutor
from autogen_ext.tools.mcp import StdioServerParams, mcp_server_tools


def init_model_client():
    token_provider = AzureTokenProvider(
        DefaultAzureCredential(),
        "https://cognitiveservices.azure.com/.default",
    )

    model_client = AzureOpenAIChatCompletionClient(
        azure_deployment="gpt-4o",  # Replace with your Azure deployment name
        model="gpt-4o",  # Replace with your model name
        api_version="2024-12-01-preview",
        azure_endpoint="https://mawasileazureopenai.openai.azure.com/",  # Replace with your endpoint
        azure_ad_token_provider=token_provider,  # Optional if you choose key-based authentication
    )
    return model_client

async def init_tools():
    fetch_mcp_server = StdioServerParams(command="uv", args=[
        "run",
        #"-v",
        "mcp-server-fetch"])
    tools = await mcp_server_tools(fetch_mcp_server)

    time_mcp_server = StdioServerParams(command="uv", args=[
        "run",
        #"-v",
        "mcp-server-time"])
    tools.extend(await mcp_server_tools(time_mcp_server))

    mitmproxy_mcp_server = StdioServerParams(command="uv", args=[
        "--directory",
        "/workspaces/terraform-provider-power-platform/.mcp_servers/mitmproxy-mcp",
        # "-v",
        "run",
        "mitmproxy-mcp"])
    tools.extend(await mcp_server_tools(mitmproxy_mcp_server))
    
    return tools

def init_magentic_one_agent(model_client, tools):
    tools_assistant = AssistantAgent("ToolsAssistant", model_client=model_client, tools=tools, reflect_on_tool_use=True, system_message="You are a helpful assistant that providers mcp-server-fetch,mcp-server-time tool,mitmproxy-mcp for fetching data from the internet", description="AI Assistant that provides mcp-server-fetch,mcp-server-time tool,mitmproxy-mcp to help with tasks.")
    terminal = CodeExecutorAgent("ComputerTerminal",code_executor=LocalCommandLineCodeExecutor())
    
    team = MagenticOneGroupChat([tools_assistant ,terminal], model_client=model_client)

    return team

