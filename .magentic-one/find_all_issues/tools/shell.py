import subprocess
from autogen_core import CancellationToken
from autogen_core.tools import FunctionTool
from typing_extensions import Annotated

async def execute_shell_command(command_string: Annotated[str, "The shell command to execute"]) -> str:
    """
    Executes a shell command and returns its output as a string.
    
    Args:
        command_string (str): The shell command with parameters to execute
        
    Returns:
        str: The command's output (stdout)
        
    Raises:
        subprocess.CalledProcessError: If the command exits with a non-zero status
    """
    try:
        # Execute the command and capture its output
        result = subprocess.run(
            command_string,
            shell=True,
            check=True,
            text=True,
            capture_output=True
        )
        return result.stdout
    except subprocess.CalledProcessError as e:
        # If you want to handle errors differently, you can modify this part
        print(f"Error executing command: {e}")
        print(f"Command stderr: {e.stderr}")
        raise

execute_shell_command_tool = FunctionTool(execute_shell_command, description="Executes a shell command and returns its output as a string.")
