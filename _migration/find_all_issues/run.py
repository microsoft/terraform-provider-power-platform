import asyncio
import os
import glob
from pathlib import Path

from agent import CodeQualityAgent

def concatenate_files(directory_path: str) -> str:
    """
    Concatenates all files in a directory with path and content markers.
    
    Args:
        directory_path: Path to the directory to process
    
    Returns:
        Concatenated content as string
    """
    concatenated_content = ""
    
    # Get all files recursively, excluding _agent folder
    files = glob.glob(f"{directory_path}/**/*", recursive=True)
    #files = [f for f in files if os.path.isfile(f) and "_agent" not in f]

    for file_path in sorted(files):
        try:
            print(f"Processing file: {file_path}")
            with open(file_path, 'r', encoding='utf-8', errors='ignore') as file:
                content = file.read()
                
            # Add the formatted output
            concatenated_content += f"<path>{file_path}</path>\n"
            concatenated_content += f"<content>{content}</content>\n\n"
            
        except Exception as e:
            print(f"Error reading {file_path}: {e}")
            continue
    
    return concatenated_content

async def main() -> None:
    project_root = "/workspaces/Copilot-Studio-with-Azure-AI-Search/src/powerplatform/copilot_studio_gold_agent"
    
    # Option 1: Concatenate all files in the project and analyze with agent
    concatenated = concatenate_files(project_root)
    
    print("Files concatenated successfully!")
    print(f"Total content length: {len(concatenated)} characters")

    code_agent = CodeQualityAgent()
    await code_agent.run(content=concatenated)

    
    # Option 2: Analyze markdown files with the agent
    #md_files = glob.glob(f"{project_root}/**/*.md", recursive=True)
    
    # Filter out files in the _agent folder
    #md_files = [f for f in md_files if "_agent" not in f]

    # code_agent = CodeQualityAgent()
    # for file_path in md_files:
    #     absolute_path = os.path.abspath(file_path)
    #     print(f"ANALYZING: {absolute_path}")
    #     await code_agent.run(file_path=file_path)


if __name__ == "__main__":
    asyncio.run(main())
