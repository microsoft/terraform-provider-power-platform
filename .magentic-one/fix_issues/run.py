import asyncio
import os
import glob
from pathlib import Path

from agent import CodeFixAgent

async def main() -> None:

    code_agent = CodeFixAgent()
    await code_agent.run()

if __name__ == "__main__":
    asyncio.run(main())
