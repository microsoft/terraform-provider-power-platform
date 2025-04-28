import asyncio
import os
import glob
from pathlib import Path

from agent import CodeQualityAgent

async def main() -> None:
    project_root = "/workspaces/terraform-provider-power-platform"
    go_files = glob.glob(f"{project_root}/**/*.go", recursive=True)

    code_agent = CodeQualityAgent()
    for file_path in go_files:
        absolute_path = os.path.abspath(file_path)
        print(f"ANALYZING: {absolute_path}")
        await code_agent.run(file_path="/workspaces/terraform-provider-power-platform/internal/api/auth.go")

if __name__ == "__main__":
    asyncio.run(main())
