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
        """Initialize the CodeQualityAgent with model client."""
        # Initialize token provider and model client once
        self.token_provider = AzureTokenProvider(
            DefaultAzureCredential(),
            "https://cognitiveservices.azure.com/.default",
        )

        self.model_client = AzureOpenAIChatCompletionClient(
            azure_deployment="gpt-4.1",  # Replace with your Azure deployment name
            model="gpt-4.1",  # Replace with your model name
            api_version="2024-12-01-preview",
            azure_endpoint="https://mawasileazureopenai.openai.azure.com/",  # Replace with your endpoint
            azure_ad_token_provider=self.token_provider,  # Optional if you choose key-based authentication
        )

    def _create_agents_and_team(self):
        """Create the agents and team setup that's shared between both run methods."""
        primary_agent = AssistantAgent(
            "CodeQualityAgent",
            model_client=self.model_client,
            tools=[execute_shell_command_tool, execute_save_file_tool],
            reflect_on_tool_use=True,
            system_message=CODE_QUALITY_AGENT_PROMPT,
        )

        critic_agent = AssistantAgent(
            "critic",
            model_client=self.model_client,
            system_message=CRITIC_AGENT_PROMPT,
        )

        max_msg_termination = MaxMessageTermination(max_messages=30)
        text_termination = TextMentionTermination("APPROVE")
        combined_termination = max_msg_termination | text_termination

        round_robin_team = RoundRobinGroupChat([primary_agent, critic_agent], termination_condition=combined_termination)
        
        return round_robin_team

    async def _run_on_content(self, content):
        """Analyze content directly without reading from a file."""
        round_robin_team = self._create_agents_and_team()
        
        prompt = f"""# Role and Objective
    You are an agent - please keep going until the user's query is completely resolved, before ending your turn and yielding back to the user. Only terminate your turn when you are sure that the problem is solved.
    Use your tools to gather the relevant information: do NOT guess or make up an answer.
    Your task is to analyze the provided content and find all issues.
    You will save the output from the tool in a markdown files as instructed below.

    # Instructions:
    - analyze the provided content below
    - consider all provider files as a single content piece. Expect dependencies between files, so analyze them together.
    - find all issues in the content
        - for every issue found, provide feedback filling up the following markdown template for that issue. 
        - try to create single issue for same findings in different files, if they are similar enough
        - markdown file name should correspond with the content type, issue and severity
        - fix the markdown file content if it has issues like unclosed code blocks or lack of blank lines under headers
        - always validate that the markdown file has all the content from the tool's output
        - analyze only the content provided by the user
        - do not analyze additional files, or execute other actions.
        - do not ask for any clarifications, just start analyzing the content
        - ignore the comment lines in the content
    - Always save the findings in files on the disk. 

    # Content to analyze:
    {content}
    """
        await Console(round_robin_team.run_stream(task=prompt))
        
    async def _run_on_file(self, file_path):
        """Analyze a specific file by reading it first."""
        round_robin_team = self._create_agents_and_team()

        prompt = f"""# Role and Objective
    You are an agent - please keep going until the user’s query is completely resolved, before ending your turn and yielding back to the user. Only terminate your turn when you are sure that the problem is solved.
    Use your tools to gather the relevant information: do NOT guess or make up an answer.
    Your task is to read all files from a directory give by the user and use send team one by one to the code analyze tool. 
    You will save the output from the tool in a markdown files as instructed below.

    # Instructions:
    - go to the file {file_path}
    - analyze the file and find all issues
    - for every issue found, provide feedback filling up the following markdown template for that issue:
    - markdown file name should correspond with the analyzed file name, issue and serverity
    - fix the markdown file content if it is has issues like unclosed code blocks or lack of blank lines under headers
    - always validate that the markdown file has all the content from the tool's output
    - analyze only single file given by the user
    - do not analyze additional files, or execute other actions.
    - do not ask for any clarifications, just start analyzing the file
    - ignore the comment lines in the file
    """
        await Console(round_robin_team.run_stream(task=prompt))

    async def run(self, file_path=None, content=None):
        """
        Backward-compatible run method that delegates to the appropriate specific method.
        
        Args:
            file_path: Path to file to analyze (optional)
            content: Content to analyze directly (optional)
        """
        if content is not None:
            await self._run_on_content(content)
        elif file_path is not None:
            await self._run_on_file(file_path)
        else:
            raise ValueError("Either file_path or content must be provided")
