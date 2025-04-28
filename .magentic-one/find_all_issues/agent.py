from autogen_agentchat.agents import AssistantAgent
from autogen_ext.models.openai import AzureOpenAIChatCompletionClient
from azure.identity import DefaultAzureCredential
from autogen_ext.auth.azure import AzureTokenProvider
from autogen_agentchat.conditions import MaxMessageTermination, TextMentionTermination
from autogen_agentchat.teams import RoundRobinGroupChat
from autogen_agentchat.ui import Console
from prompts import CODE_QUALITY_AGENT_PROMPT, CRITIC_AGENT_PROMPT
from tools.shell import execute_shell_command_tool
from tools.file_manager import execute_save_file_tool

class CodeQualityAgent:
    """
    A class that analyzes code quality of files in the terraform-provider-power-platform project.
    """
    
    def __init__(self):
        """Initialize the CodeQualityAgent."""
    
    async def run(self, file_path):
        token_provider = AzureTokenProvider(
            DefaultAzureCredential(),
            "https://cognitiveservices.azure.com/.default",
        )

        # Create an agent that can use the fetch tool.
        model_client = AzureOpenAIChatCompletionClient(
            azure_deployment="gpt-4o",  # Replace with your Azure deployment name
            model="gpt-4o",  # Replace with your model name
            api_version="2024-12-01-preview",
            azure_endpoint="https://mawasileazureopenai.openai.azure.com/",  # Replace with your endpoint
            azure_ad_token_provider=token_provider,  # Optional if you choose key-based authentication
        )

        primary_agent = AssistantAgent(
            "CodeQualityAgent",
            model_client=model_client,
            tools=[execute_shell_command_tool, execute_save_file_tool],
            reflect_on_tool_use=True,
            system_message=CODE_QUALITY_AGENT_PROMPT,
        )

        critic_agent = AssistantAgent(
            "critic",
            model_client=model_client,
            system_message=CRITIC_AGENT_PROMPT,
        )

        max_msg_termination = MaxMessageTermination(max_messages=30)
        text_termination = TextMentionTermination("APPROVE")
        combined_termination = max_msg_termination | text_termination

        round_robin_team = RoundRobinGroupChat([primary_agent, critic_agent], termination_condition=combined_termination)

        folder_path = "/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues"

        prompt = f"""# Role and Objective
    You are an agent - please keep going until the user’s query is completely resolved, before ending your turn and yielding back to the user. Only terminate your turn when you are sure that the problem is solved.
    Use your tools to gather the relevant information: do NOT guess or make up an answer.
    Your task is to read all files from a directory give by the user and use send team one by one to the code analyze tool. 
    You will save the output from the tool in a markdown files as instructed below.

    # Instructions:
    - go to the file {file_path}
    - analyze the file and find all issues
    - for every issue found, provide feedback filling up the following markdown template for that issue:
    - save each output markdown output from the tool in a separate markdown file beased on the following rule:
        - severity critical will be saved in {folder_path}/issues_found/critical
        - severity high will be saved in{folder_path}/issues_found/high
        - severity medium  will be saved in {folder_path}/issues_found/medium
        - severity low  will be saved in {folder_path}/issues_found/low
    - markdown file name should correspond with the analyzed file name and issue
    - fix the markdown file content if it is has issues like unclosed code blocks or lack of blank lines under headers
    - always validate that the markdown file has all the content from the tool's output
    - analyze only single file given by the user
    - do not analyze additional files, or execute other actions.
    - do not ask for any clarifications, just start analyzing the file
    """
        await Console(round_robin_team.run_stream(task=prompt))
        #result = await round_robin_team.run(task=prompt, cancellation_token=None)
    