import os
import json
from typing import Any, Dict, List, Union
from mitmproxy import io

# Directory where mitmproxy dump files are stored
DUMP_DIR = "/Users/lucas/Coding/mitmproxy-mcp/dumps"

# Cache for storing flows per session
FLOW_CACHE = {}

async def get_flows_from_dump(session_id: str) -> list:
    """
    Retrieves flows from the dump file, using the cache if available.
    """
    dump_file = os.path.join(DUMP_DIR, f"{session_id}.dump")
    if not os.path.exists(dump_file):
        raise FileNotFoundError("Session not found")

    if session_id in FLOW_CACHE:
        return FLOW_CACHE[session_id]
    else:
        with open(dump_file, "rb") as f:
            reader = io.FlowReader(f)
            flows = list(reader.stream())
        FLOW_CACHE[session_id] = flows
        return flows

def parse_json_content(content: bytes, headers: dict) -> Union[Dict, str, bytes]:
    """
    Attempts to parse content as JSON if the content type indicates JSON.
    Returns the parsed JSON or the raw content if parsing fails.
    """
    content_type = headers.get("Content-Type", "").lower() if headers else ""
    
    if "application/json" in content_type or "text/json" in content_type:
        try:
            return json.loads(content.decode(errors="ignore"))
        except json.JSONDecodeError:
            return content.decode(errors="ignore")
    return content.decode(errors="ignore")
