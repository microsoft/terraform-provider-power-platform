// source code: <https://github.com/microsoft/autogen/tree/main>

# Installs

- python extension for vscode

- pip install "autogen-agentchat" "autogen-ext[magentic-one,openai,azure]"
- pip install autogen-core
- pip install -U "autogen-ext[mcp]"
- pip install mcp
- pip install json_schema_to_pydantic
- pip install uvx
- uvx install mcp-server-fetch


## Services

- Create Azure Open AI service in eastus2, because it has all the [models](https://learn.microsoft.com/en-us/azure/ai-services/openai/concepts/models?tabs=global-standard%2Cstandard-chat-completions#model-summary-table-and-region-availability)
- to use SPN and az login, you need to go to subscription and add "Cognitive Services OpenAI User" to your SPN
