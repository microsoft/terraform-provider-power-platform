from autogen_agentchat.agents import AssistantAgent
from autogen_ext.models.openai import AzureOpenAIChatCompletionClient
from azure.identity import DefaultAzureCredential
from autogen_ext.auth.azure import AzureTokenProvider
from autogen_agentchat.conditions import MaxMessageTermination, TextMentionTermination
from autogen_agentchat.teams import RoundRobinGroupChat
from autogen_agentchat.ui import Console
from prompts import prep_code_agent_prompt, CRITIC_AGENT_PROMPT
from tools.shell import execute_shell_command_tool
from tools.file_manager import execute_save_file_tool

class CodeFixAgent:
    
    def __init__(self):
        """Initialize the CodeFixAgent with the necessary tools and agents."""
    
    async def run(self):
        repo_dir_path = "/workspaces/terraform-provider-power-platform"
        instructions_file_path = f"{repo_dir_path}/.github/copilot-instructions.md"
        issue_file_path = f"{repo_dir_path}/.magentic-one/fix_issues/issue.md"
        issue_id = "12345"


        token_provider = AzureTokenProvider(
            DefaultAzureCredential(),
            "https://cognitiveservices.azure.com/.default",
        )

        # Create an agent that can use the fetch tool.
        model_client = AzureOpenAIChatCompletionClient(
            azure_deployment="gpt-4o",  # Replace with your Azure deployment name
            model="gpt-4o",  # Replace with your model name
            api_version="2024-12-01-preview",
            azure_endpoint="https://<<link>>.openai.azure.com/",  # Replace with your endpoint
            azure_ad_token_provider=token_provider,  # Optional if you choose key-based authentication
        )


        primary_agent = AssistantAgent(
            "CodeFixAgent",
            model_client=model_client,
            tools=[execute_shell_command_tool, execute_save_file_tool],
            reflect_on_tool_use=True,
            system_message=prep_code_agent_prompt(issue_file_path, repo_dir_path, instructions_file_path),
        )

        critic_agent = AssistantAgent(
            "critic",
            model_client=model_client,
            system_message=CRITIC_AGENT_PROMPT,
        )

        max_msg_termination = MaxMessageTermination(max_messages=100)
        text_termination = TextMentionTermination("APPROVE")
        combined_termination = max_msg_termination | text_termination

        round_robin_team = RoundRobinGroupChat([primary_agent, critic_agent], termination_condition=combined_termination)


        prompt = f"""# Role and Objective
    You are an agent - please keep going until the user query is completely resolved, before ending your turn and yielding back to the user. Only terminate your turn when you are sure that the problem is solved.
    Use your tools to gather the relevant information: do NOT guess or make up an answer.

    You are responsible for fixing code issue described in {issue_file_path}
    To run the command tools you need to follow the instructions provided in the instructions file.
    You have to be in the directory '{repo_dir_path}' to run the command tools.
    Issue number to use is '{issue_id}'.
    When editing files, always write the whole content of the file, not just the change you want to make.

    # Instructions:
        - Read the instructions file located at '{instructions_file_path}' that will guide you on how to run the tools and how the repository is structured.
        - Go to the file {issue_file_path} and fix the issue described in the file
        - Run linter and if it returns issues, find the files mentioned and fix them.
        - Add changelog using `changie new`
    """
        await Console(round_robin_team.run_stream(task=prompt))

    