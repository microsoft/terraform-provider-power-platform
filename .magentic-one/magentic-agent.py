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
from autogen_core.memory import ListMemory, MemoryContent, MemoryMimeType


async def main() -> None:

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
    
    # surfer = MultimodalWebSurfer(
    #     "WebSurfer",
    #     model_client=model_client,
    # )

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


    # Get the prompt from standard input
    print("Enter your prompt: ", end="", flush=True)
    user_prompt = await asyncio.get_event_loop().run_in_executor(None, input)
    
    codebase_path =  os.path.abspath("/workspaces/terraform-provider-power-platform")
    # Read the Copilot instructions file
    instructions_path = os.path.join(codebase_path, ".github", "copilot-instructions.md")

    try:
        # Concatenate instructions and user prompt
        hacks = "Try to use tools that ToolsAssistant has, before tryting to code a solution.\n"
        hacks = "If you can't edit or save file, use shell commands to do it.\n"
        hacks += "If you can't find the file, use the command 'find' to search for it.\n"
        user_prompt = f"the absolute path where the codebase is here: {codebase_path}\n\nInformation about the repository is here {instructions_path}\n\n{hacks}\n\n{user_prompt}"
        #user_prompt = f"\n\n{hacks}\n\n{user_prompt}"
    except FileNotFoundError:
        print(f"Warning: Instructions file not found at {instructions_path}")
        return

        # Initialize user memory
    # user_memory = ListMemory()
    # await user_memory.add(MemoryContent(content="Use tools that ToolsAssistant has, before trying to code a solution", mime_type=MemoryMimeType.TEXT))
    # await user_memory.add(MemoryContent(content="Use ComputerTerminal Assistant to execute bash shell commands", mime_type=MemoryMimeType.TEXT))
    # await user_memory.add(MemoryContent(content=f"Absolute path where the codebase is here: {codebase_path}", mime_type=MemoryMimeType.TEXT))
    # await user_memory.add(MemoryContent(content=f"Information about the repository is here {instructions_path}", mime_type=MemoryMimeType.TEXT))

    tools_assistant = AssistantAgent("ToolsAssistant", model_client=model_client, tools=tools, reflect_on_tool_use=True, system_message="You are a helpful assistant that providers mcp-server-fetch,mcp-server-time tool,mitmproxy-mcp for fetching data from the internet", description="AI Assistant that provides mcp-server-fetch,mcp-server-time tool,mitmproxy-mcp to help with tasks.")
    #memory_assistant = AssistantAgent("MemoryAssistant", model_client=model_client, memory=[user_memory], system_message="You are a helpful assistant that provides memorized information about the most important facts to finish tasks", description="AI Assistant that provides memorized information about the most important facts to finish tasks.")
    #file_surfer = FileSurfer("FileSurfer",model_client=model_client, base_path="/workspaces/terraform-provider-power-platform")
    #coder = MagenticOneCoderAgent("Coder",model_client=model_client)
    terminal = CodeExecutorAgent("ComputerTerminal",code_executor=LocalCommandLineCodeExecutor())
    
    team = MagenticOneGroupChat([tools_assistant ,terminal], model_client=model_client)
    #team = MagenticOneGroupChat([file_surfer, coder, terminal], model_client=model_client)

    await Console(team.run_stream(task=user_prompt))

asyncio.run(main())
