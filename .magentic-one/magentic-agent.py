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
            azure_endpoint="https://<<your_az_foundry_endpoin>>.openai.azure.com/",  # Replace with your endpoint
            azure_ad_token_provider=token_provider,  # Optional if you choose key-based authentication
        )
    
    # surfer = MultimodalWebSurfer(
    #     "WebSurfer",
    #     model_client=model_client,
    # )


    file_mcp_server = StdioServerParams(command="uv", args=[
        "--directory",
        "/workspaces/terraform-provider-power-platform/.mcp_servers/ai-mcp",
       # "-v",
        "run",
        "server.py",], read_timeout_seconds=60)
    tools = await mcp_server_tools(file_mcp_server)


    # shell_mcp_server = StdioServerParams(command="uv", args=[
    #     #"-v",
    #     "run",
    #     "mcp-shell-server",
    #    ], 
    #    env={
    #     "ALLOW_COMMANDS": "ls,pwd,echo,cat,head,tail,find,grep,wc,touch,cd"
    #    },
    #    read_timeout_seconds=60)
    # tools.extend(await mcp_server_tools(shell_mcp_server))



    # fetch_mcp_server = StdioServerParams(command="uv", args=[
    #     "run",
    #     #"-v",
    #     "mcp-server-fetch"], read_timeout_seconds=60)
    # tools = await mcp_server_tools(fetch_mcp_server)

    # time_mcp_server = StdioServerParams(command="uv", args=[
    #     "run",
    #     #"-v",
    #     "mcp-server-time"])
    # tools.extend(await mcp_server_tools(time_mcp_server))

    # mitmproxy_mcp_server = StdioServerParams(command="uv", args=[
    #     "--directory",
    #     "/workspaces/terraform-provider-power-platform/.mcp_servers/mitmproxy-mcp",
    #    # "-v",
    #     "run",
    #     "mitmproxy-mcp"])
    # tools.extend(await mcp_server_tools(mitmproxy_mcp_server))


    # Get the prompt from standard input
    print("Enter your prompt: ", end="", flush=True)
    user_prompt = await asyncio.get_event_loop().run_in_executor(None, input)
    
    codebase_path =  os.path.abspath("/workspaces/terraform-provider-power-platform")
    # Read the Copilot instructions file
    instructions_path = os.path.join(codebase_path, ".github", "copilot-instructions.md")

    try:
        # Concatenate instructions and user prompt
        hacks = "Try to use tools that ToolsAssistant has, before tryting to code a solution.\n"
       # hacks = "If you can't edit or save file, use tools.\n"
        #hacks += "Use tools provided by ToolsAssistant to edit or save files.\n"
        #hacks += "If you can't find the file, use the command 'find' to search for it.\n"
        #user_prompt = f"the absolute path where the codebase is here: {codebase_path}\n\nInformation about the repository is here {instructions_path}\n\n{hacks}\n\n{user_prompt}"
        user_prompt = f"\n\n{hacks}\n\n{user_prompt}"
    except FileNotFoundError:
        print(f"Warning: Instructions file not found at {instructions_path}")
        return

        # Initialize user memory
    user_memory = ListMemory()
    await user_memory.add(MemoryContent(content="Use tools that ToolsAssistant has, before trying to code a solution", mime_type=MemoryMimeType.TEXT))
    await user_memory.add(MemoryContent(content="Use ComputerTerminal Assistant to execute bash shell commands", mime_type=MemoryMimeType.TEXT))
    await user_memory.add(MemoryContent(content=f"Absolute path where the codebase is here: {codebase_path}", mime_type=MemoryMimeType.TEXT))
   # await user_memory.add(MemoryContent(content=f"Information about the repository is here {instructions_path}", mime_type=MemoryMimeType.TEXT))

    await user_memory.add(MemoryContent(content="""
Analyze all the .go file under /workspaces/terraform-provider-power-platform/internal/ and find any code clarity, security and other issues that you will find. For each issue found read the template .github/prompts/ai_bug_report.md and use it to create marddown file in .github/prompts/issues_found folder describing the issue.

File name should be markdown file based on template located in .github/prompts/ai_bug_report.md
Template file has elements using "<<>>" that should be replaced, fill all of them
Validate the markdown file you've create is valid
Markdown file name should correspond with the analyzed file name, issue and serverity
Markdown file content should be saved using the following rules:
    - issues that represent error handling, panic, or control flow issue shuld be saved in {folder_path}/issues_found/error_handling
    - issues that represent variable, function, or type naming should be saved in {folder_path}/issues_found/naming
    - issues that represent resource management, memory leaks and performance should be saved in {folder_path}/issues_found/resource_management
    - issues that represent API client and HTTP issues should be saved in {folder_path}/issues_found/api_client
    - issues that represent code structure, maintainability, or readability should be saved in {folder_path}/issues_found/structure
    - issues that represent testing and quality assurance should be saved in {folder_path}/issues_found/testing
    - issues that represent type safety, validation, or data consistency should be saved in {folder_path}/issues_found/type_safety
    - any other issues that do not fit the above categories should be saved in {folder_path}/issues_found/other
                                        
Analyze all the files that you will find in  /workspaces/terraform-provider-power-platform/internal/
CodeQualityAssistant should not analyze the files.
""", mime_type=MemoryMimeType.TEXT))
    reminder_assistant = AssistantAgent("ReminderAssistant", model_client=model_client, memory=[user_memory], system_message="You are a helpful assistant that provides reminders about what user wanted", description="AI Assistant that provides reminders about what user wanted")

    tools_assistant = AssistantAgent("ToolsAssistant", model_client=model_client, tools=tools, reflect_on_tool_use=True, system_message="You are a helpful assistant that providers useful tools. You only run them for others. Do not execute them on your own.", description="AI Assistant that provides tools: save_file, analyze_code")
   # code_quality_assistant = AssistantAgent("CodeQualityAssistant", model_client=model_client, tools=tools, system_message="You are a helpful assistant that provides code quality checks and suggestions", description="AI Assistant that provides code quality checks and suggestions.")
    #memory_assistant = AssistantAgent("MemoryAssistant", model_client=model_client, memory=[user_memory], system_message="You are a helpful assistant that provides memorized information about the most important facts to finish tasks", description="AI Assistant that provides memorized information about the most important facts to finish tasks.")
    #file_surfer = FileSurfer("FileSurfer",model_client=model_client, base_path=codebase_path)
    #coder = MagenticOneCoderAgent("Coder",model_client=model_client)
    terminal = CodeExecutorAgent("ComputerTerminal",code_executor=LocalCommandLineCodeExecutor(work_dir=codebase_path))
    team = MagenticOneGroupChat([tools_assistant, reminder_assistant, terminal], model_client=model_client)
    #team = MagenticOneGroupChat([file_surfer, coder, terminal], model_client=model_client)

    await Console(team.run_stream(task=user_prompt))

asyncio.run(main())
