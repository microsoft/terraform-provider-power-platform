from typing import Any
#import httpx
from mcp.server.fastmcp import FastMCP
from magentic_one_agent import init_model_client, init_tools, init_magentic_one_agent

model_client = init_model_client()

# Initialize FastMCP server
mcp = FastMCP("ai-mcp")

@mcp.tool()
async def get_magentic_one_agent(prompt: str) -> str:
    """Ask Magentic One multi-agent to do something for you

    Args:
        prompt: The prompt to ask the Magentic One agent.
    
    Returns:
        A string describing the result of the task.
    """
    tools = await init_tools()

    # Initialize the Magentic One agent
    agent = init_magentic_one_agent(model_client, tools)
    result = await agent.run(task=prompt, cancellation_token=None)
    return result.messages[0].content

@mcp.tool()
async def get_plan(task: str) -> str:
    """Get plan from a planner assitant

    Args:
        task: The task to plan for.
    Returns:
        A string describing the plan.
    """
    
    tools = await init_tools()

    # Initialize the Magentic One agent
    agent = init_magentic_one_agent(model_client, tools)

    prompt = f"""
        You are a planner assistant. Your job is to create a plan for the task.
        Task: {task}
    """

    result = await agent.run(task=prompt, cancellation_token=None)
    return result.messages[0].content

if __name__ == "__main__":
    # Initialize and run the server
    mcp.run(transport='stdio')
