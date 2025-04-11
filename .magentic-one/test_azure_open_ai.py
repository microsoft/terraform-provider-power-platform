#https://microsoft.github.io/autogen/stable/user-guide/agentchat-user-guide/tutorial/models.html#azure-openai

import asyncio
from autogen_ext.auth.azure import AzureTokenProvider
from autogen_ext.models.openai import AzureOpenAIChatCompletionClient
from azure.identity import DefaultAzureCredential
from autogen_core.models import UserMessage

async def main():
    try:
        # Create the token provider
        token_provider = AzureTokenProvider(
            DefaultAzureCredential(),
            "https://cognitiveservices.azure.com/.default",
        )

        # Initialize the Azure OpenAI client
        az_model_client = AzureOpenAIChatCompletionClient(
            azure_deployment="gpt-4o",  # Replace with your Azure deployment name
            model="gpt-4o",  # Replace with your model name
            api_version="2024-12-01-preview",
            azure_endpoint="https://<<your_az_foundry_endpoin>>.openai.azure.com/",  # Replace with your endpoint
            azure_ad_token_provider=token_provider,  # Optional if you choose key-based authentication
            # api_key="...",  # Uncomment and replace for key-based authentication
        )

        # Send a message to the model
        result = await az_model_client.create([UserMessage(content="What is the capital of France?", source="user")])
        print(result)

    except Exception as e:
        print(f"An error occurred: {e}")
    finally:
        # Ensure the client is closed
        await az_model_client.close()

# Run the async main function
if __name__ == "__main__":
    asyncio.run(main())
