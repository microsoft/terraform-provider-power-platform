import asyncio
import os
import json
import re
import base64
from typing import Any, Dict, List, Optional, Union, Tuple
from mitmproxy import io

from mcp.server.models import InitializationOptions
import mcp.types as types
from mcp.server import NotificationOptions, Server
import mcp.server.stdio

# Directory where mitmproxy dump files are stored
DUMP_DIR = "/workspaces/terraform-provider-power-platform"

server = Server("mitmproxy-mcp")

# Cache for storing flows per session
FLOW_CACHE = {}

# Maximum content size in bytes before switching to structure preview
MAX_CONTENT_SIZE = 9999999

# Known bot protection systems and their signatures
BOT_PROTECTION_SIGNATURES = {
    "Cloudflare": [
        r"cf-ray",  # Cloudflare Ray ID header
        r"__cf_bm",  # Cloudflare Bot Management cookie
        r"cf_clearance",  # Cloudflare challenge clearance cookie
        r"\"why_captcha\"",  # Common in Cloudflare challenge responses
        r"challenge-platform",  # Used in challenge scripts
        r"turnstile\.js",  # Cloudflare Turnstile
    ],
    "Akamai Bot Manager": [
        r"_abck=",  # Akamai Bot Manager cookie
        r"akam_", # Akamai cookie prefix
        r"bm_sz", # Bot Manager cookie
        r"sensor_data", # Bot detection data
    ],
    "PerimeterX": [
        r"_px\d?=",  # PerimeterX cookies
        r"px\.js", # PerimeterX script
        r"px-captcha", # PerimeterX captcha
    ],
    "DataDome": [
        r"datadome=",  # DataDome cookie
        r"datadome\.js", # DataDome script
        r"_dd_s",  # DataDome session cookie
    ],
    "reCAPTCHA": [
        r"google\.com/recaptcha",
        r"recaptcha\.net",
        r"g-recaptcha",
    ],
    "hCaptcha": [
        r"hcaptcha\.com",
        r"h-captcha",
    ],
    "Generic Bot Detection": [
        r"bot=",  # Generic bot cookie
        r"captcha", # Generic captcha reference
        r"challenge",  # Generic challenge term
        r"detected automated traffic", # Common message
        r"verify you are human", # Common message
    ]
}

@server.list_tools()
async def handle_list_tools() -> list[types.Tool]:
    """
    List available tools.
    Each tool specifies its arguments using JSON Schema validation.
    """
    return [
        types.Tool(
            name="list_flows",
            description="Lists all detailed HTTP request/response data flows including headers, content (or structure preview for large JSON), and metadata from specified flows",
            inputSchema={
                "type": "object",
                "properties": {
                    "session_id": {
                        "type": "string",
                        "description": "The ID of the session to list flows from",
                        "default": "mitmproxy.dump"
                    }
                },
                "required": [""]
            }
        ),
        types.Tool(
            name="get_flow_details",
            description="Lists HTTP requests/responses from a mitmproxy capture session, showing method, URL, and status codes",
            inputSchema={
                "type": "object",
                "properties": {
                    "session_id": {
                        "type": "string",
                        "description": "The ID of the session",
                        "default": "mitmproxy.dump"
                    },
                    "flow_indexes": {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "description": "The indexes of the flows"
                    },
                    "include_content": {
                        "type": "boolean",
                        "description": "Whether to include full content in the response (default: true)",
                        "default": True
                    }
                },
                "required": ["session_id", "flow_indexes"]
            }
        ),
        types.Tool(
            name="extract_json_fields",
            description="Extract specific fields from JSON content in a flow using JSONPath expressions",
            inputSchema={
                "type": "object",
                "properties": {
                    "session_id": {
                        "type": "string",
                        "description": "The ID of the session",
                        "default": "mitmproxy.dump"
                    },
                    "flow_index": {
                        "type": "integer",
                        "description": "The index of the flow"
                    },
                    "content_type": {
                        "type": "string",
                        "enum": ["request", "response"],
                        "description": "Whether to extract from request or response content"
                    },
                    "json_paths": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "JSONPath expressions to extract (e.g. ['$.data.users', '$.metadata.timestamp'])"
                    }
                },
                "required": ["session_id", "flow_index", "content_type", "json_paths"]
            }
        ),
        types.Tool(
            name="analyze_protection",
            description="Analyze flow for bot protection mechanisms and extract challenge details",
            inputSchema={
                "type": "object",
                "properties": {
                    "session_id": {
                        "type": "string",
                        "description": "The ID of the session",
                        "default": "mitmproxy.dump"
                    },
                    "flow_index": {
                        "type": "integer",
                        "description": "The index of the flow to analyze"
                    },
                    "extract_scripts": {
                        "type": "boolean",
                        "description": "Whether to extract and analyze JavaScript from the response (default: true)",
                        "default": True
                    }
                },
                "required": ["session_id", "flow_index"]
            }
        )
    ]

async def get_flows_from_dump(session_id: str) -> list:
    """
    Retrieves flows from the dump file, using the cache if available.
    """
    dump_file = os.path.join(DUMP_DIR, f"{session_id}")
    if not os.path.exists(dump_file):
        raise FileNotFoundError(f"Session not found: {dump_file}")

    if session_id in FLOW_CACHE:
        return FLOW_CACHE[session_id]
    else:
        with open(dump_file, "rb") as f:
            reader = io.FlowReader(f)
            flows = list(reader.stream())
        FLOW_CACHE[session_id] = flows
        return flows

async def list_flows(arguments: dict) -> list[types.TextContent]:
    """
    Lists HTTP flows from a mitmproxy dump file.
    """
    session_id = "mitmproxy.dump" #arguments.get("session_id")
    #if not session_id:
    #    return [types.TextContent(type="text", text="Error: Missing session_id")]

    try:
        flows = await get_flows_from_dump(session_id)

        flow_list = []
        for i, flow in enumerate(flows):
            if flow.type == "http":
                request = flow.request
                response = flow.response
                flow_info = {
                    "index": i,
                    "method": request.method,
                    "url": request.url,
                    "status": response.status_code if response else None
                }
                flow_list.append(flow_info)

        return [types.TextContent(type="text", text=json.dumps(flow_list, indent=2))]
    except FileNotFoundError as e:
        return [types.TextContent(type="text", text=f"Error reading file: {str(e)}")]
    except Exception as e:
        return [types.TextContent(type="text", text=f"Error reading flows: {str(e)}")]

def generate_json_structure(json_data: Any, max_depth: int = 2, current_depth: int = 0) -> Any:
    """
    Generate a simplified structure of JSON content, showing keys and types
    but replacing actual values with type indicators after a certain depth.
    """
    if current_depth >= max_depth:
        if isinstance(json_data, dict):
            return {"...": f"{len(json_data)} keys"}
        elif isinstance(json_data, list):
            return f"[{len(json_data)} items]"
        else:
            return f"({type(json_data).__name__})"
    
    if isinstance(json_data, dict):
        result = {}
        for key, value in json_data.items():
            result[key] = generate_json_structure(value, max_depth, current_depth + 1)
        return result
    elif isinstance(json_data, list):
        if not json_data:
            return []
        # For lists, show structure of first item and count
        sample = generate_json_structure(json_data[0], max_depth, current_depth + 1)
        return [sample, f"... ({len(json_data)-1} more items)"] if len(json_data) > 1 else [sample]
    else:
        return f"({type(json_data).__name__})"

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

def extract_with_jsonpath(json_data: Any, path: str) -> Any:
    """
    Basic implementation of JSONPath extraction.
    Supports simple dot notation and array indexing.
    For more complex cases, consider using a full JSONPath library.
    """
    # Handle root object reference
    if path == "$":
        return json_data
    
    # Strip leading $ if present
    if path.startswith("$"):
        path = path[1:]
    if path.startswith("."):
        path = path[1:]
        
    parts = []
    # Parse the path - handle both dot notation and brackets
    current = ""
    in_brackets = False
    for char in path:
        if char == "[":
            if current:
                parts.append(current)
                current = ""
            in_brackets = True
        elif char == "]":
            if in_brackets:
                try:
                    # Handle array index
                    parts.append(int(current.strip()))
                except ValueError:
                    # Handle quoted key
                    quoted = current.strip()
                    if (quoted.startswith("'") and quoted.endswith("'")) or \
                       (quoted.startswith('"') and quoted.endswith('"')):
                        parts.append(quoted[1:-1])
                    else:
                        parts.append(quoted)
                current = ""
                in_brackets = False
        elif char == "." and not in_brackets:
            if current:
                parts.append(current)
                current = ""
        else:
            current += char
    
    if current:
        parts.append(current)
    
    # Navigate through the data
    result = json_data
    for part in parts:
        try:
            if isinstance(result, dict):
                result = result.get(part)
            elif isinstance(result, list) and isinstance(part, int):
                if 0 <= part < len(result):
                    result = result[part]
                else:
                    return None
            else:
                return None
            
            if result is None:
                break
        except Exception:
            return None
    
    return result

async def get_flow_details(arguments: dict) -> list[types.TextContent]:
    """
    Gets details of specific flows from a mitmproxy dump file.
    For large JSON content, returns structure preview instead of full content.
    """
    session_id = "mitmproxy.dump" #arguments.get("session_id")
    flow_indexes = arguments.get("flow_indexes")
    include_content = arguments.get("include_content", True)

    if not session_id:
        return [types.TextContent(type="text", text="Error: Missing session_id")]
    if not flow_indexes:
        return [types.TextContent(type="text", text="Error: Missing flow_indexes")]

    try:
        flows = await get_flows_from_dump(session_id)
        flow_details_list = []

        for flow_index in flow_indexes:
            try:
                flow = flows[flow_index]

                if flow.type == "http":
                    request = flow.request
                    response = flow.response

                    # Parse content
                    request_content = parse_json_content(request.content, dict(request.headers))
                    response_content = None
                    if response:
                        response_content = parse_json_content(response.content, dict(response.headers))
                    
                    # Handle large content
                    request_content_preview = None
                    response_content_preview = None

                    flow_details = {}
                    
                    # Check if request content is large and is JSON
                    if include_content and len(request.content) > MAX_CONTENT_SIZE and isinstance(request_content, dict):
                        request_content_preview = generate_json_structure(request_content)
                        request_content = None  # Don't include full content
                    elif include_content and len(request.content) > MAX_CONTENT_SIZE:
                        if isinstance(request_content, str):
                            request_content = request_content[:MAX_CONTENT_SIZE] + " ...[truncated]"
                        else:
                            request_content = request_content[:MAX_CONTENT_SIZE].decode(errors="ignore") + " ...[truncated]"
                        flow_details["request_content_note"] = f"Request content truncated to {MAX_CONTENT_SIZE} bytes."
                    
                    # Check if response content is large and is JSON
                    if response and include_content and len(response.content) > MAX_CONTENT_SIZE and isinstance(response_content, dict):
                        response_content_preview = generate_json_structure(response_content)
                        response_content = None  # Don't include full content
                    elif response and include_content and len(response.content) > MAX_CONTENT_SIZE:
                        if isinstance(response_content, str):
                            response_content = response_content[:MAX_CONTENT_SIZE] + " ...[truncated]"
                        else:
                            response_content = response_content[:MAX_CONTENT_SIZE].decode(errors="ignore") + " ...[truncated]"
                        flow_details["response_content_note"] = f"Response content truncated to {MAX_CONTENT_SIZE} bytes."

                    # Build flow details
                    flow_details.update( {
                        "index": flow_index,
                        "method": request.method,
                        "url": request.url,
                        "request_headers": dict(request.headers),
                        "status": response.status_code if response else None,
                        "response_headers": dict(response.headers) if response else None,
                    })
                    
                    # Add content or previews based on size
                    if include_content:
                        if request_content is not None:
                            flow_details["request_content"] = request_content
                        if request_content_preview is not None:
                            flow_details["request_content_preview"] = request_content_preview
                            flow_details["request_content_size"] = len(request.content)
                            flow_details["request_content_note"] = "Content too large to display. Use extract_json_fields tool to get specific values."
                            
                        if response_content is not None:
                            flow_details["response_content"] = response_content
                        if response_content_preview is not None:
                            flow_details["response_content_preview"] = response_content_preview
                            flow_details["response_content_size"] = len(response.content) if response else 0
                            flow_details["response_content_note"] = "Content too large to display. Use extract_json_fields tool to get specific values."
                    
                    flow_details_list.append(flow_details)
                else:
                    flow_details_list.append({"error": f"Flow {flow_index} is not an HTTP flow"})

            except IndexError:
                flow_details_list.append({"error": f"Flow index {flow_index} out of range"})

        return [types.TextContent(type="text", text=json.dumps(flow_details_list, indent=2))]

    except FileNotFoundError:
        return [types.TextContent(type="text", text="Error: Session not found b")]
    except Exception as e:
        return [types.TextContent(type="text", text=f"Error reading flow details: {str(e)}")]

async def extract_json_fields(arguments: dict) -> list[types.TextContent]:
    """
    Extract specific fields from JSON content in a flow using JSONPath expressions.
    """
    session_id = "mitmproxy.dump" #arguments.get("session_id")
    flow_index = arguments.get("flow_index")
    content_type = arguments.get("content_type")
    json_paths = arguments.get("json_paths")

    if not session_id:
        return [types.TextContent(type="text", text="Error: Missing session_id")]
    if flow_index is None:
        return [types.TextContent(type="text", text="Error: Missing flow_index")]
    if not content_type:
        return [types.TextContent(type="text", text="Error: Missing content_type")]
    if not json_paths:
        return [types.TextContent(type="text", text="Error: Missing json_paths")]

    try:
        flows = await get_flows_from_dump(session_id)
        
        try:
            flow = flows[flow_index]
            
            if flow.type != "http":
                return [types.TextContent(type="text", text=f"Error: Flow {flow_index} is not an HTTP flow")]
            
            request = flow.request
            response = flow.response
            
            # Determine which content to extract from
            content = None
            headers = None
            if content_type == "request":
                content = request.content
                headers = dict(request.headers)
            elif content_type == "response":
                if not response:
                    return [types.TextContent(type="text", text=f"Error: Flow {flow_index} has no response")]
                content = response.content
                headers = dict(response.headers)
            else:
                return [types.TextContent(type="text", text=f"Error: Invalid content_type. Must be 'request' or 'response'")]
            
            # Parse the content
            json_content = parse_json_content(content, headers)
            
            # Only extract from JSON content
            if not isinstance(json_content, (dict, list)):
                return [types.TextContent(type="text", text=f"Error: The {content_type} content is not valid JSON")]
            
            # Extract fields
            result = {}
            for path in json_paths:
                try:
                    extracted = extract_with_jsonpath(json_content, path)
                    result[path] = extracted
                except Exception as e:
                    result[path] = f"Error extracting path: {str(e)}"
            
            return [types.TextContent(type="text", text=json.dumps(result, indent=2))]
            
        except IndexError:
            return [types.TextContent(type="text", text=f"Error: Flow index {flow_index} out of range")]
            
    except FileNotFoundError:
        return [types.TextContent(type="text", text="Error: Session not found c")]
    except Exception as e:
        return [types.TextContent(type="text", text=f"Error extracting JSON fields: {str(e)}")]

def extract_javascript(html_content: str) -> List[Dict[str, Any]]:
    """
    Extract JavaScript from HTML content and provide basic analysis.
    Returns list of dictionaries with script info.
    """
    scripts = []
    
    # Extract inline scripts
    inline_pattern = r'<script[^>]*>(.*?)</script>'
    inline_scripts = re.findall(inline_pattern, html_content, re.DOTALL)
    
    for i, script in enumerate(inline_scripts):
        if len(script.strip()) > 0:
            script_info = {
                "type": "inline",
                "index": i,
                "size": len(script),
                "content": script if len(script) < 1000 else script[:1000] + "... [truncated]",
                "summary": analyze_script(script)
            }
            scripts.append(script_info)
    
    # Extract external script references
    src_pattern = r'<script[^>]*src=[\'"]([^\'"]+)[\'"][^>]*>'
    external_scripts = re.findall(src_pattern, html_content)
    
    for i, src in enumerate(external_scripts):
        script_info = {
            "type": "external",
            "index": i,
            "src": src,
            "suspicious": any(term in src.lower() for term in [
                "captcha", "challenge", "bot", "protect", "security", 
                "verify", "check", "shield", "defend", "guard"
            ])
        }
        scripts.append(script_info)
    
    return scripts

def analyze_script(script: str) -> Dict[str, Any]:
    """
    Analyze JavaScript content for common protection patterns.
    """
    analysis = {
        "potential_protection": False,
        "fingerprinting_indicators": [],
        "token_generation_indicators": [],
        "obfuscation_level": "none",
        "key_functions": []
    }
    
    # Check for fingerprinting techniques
    fingerprinting_patterns = [
        (r'navigator\.', "Browser navigator object"),
        (r'screen\.', "Screen properties"),
        (r'canvas', "Canvas fingerprinting"),
        (r'webgl', "WebGL fingerprinting"),
        (r'font', "Font enumeration"),
        (r'audio', "Audio fingerprinting"),
        (r'plugins', "Plugin enumeration"),
        (r'User-Agent', "User-Agent checking"),
        (r'platform', "Platform detection")
    ]
    
    for pattern, description in fingerprinting_patterns:
        if re.search(pattern, script, re.IGNORECASE):
            analysis["fingerprinting_indicators"].append(description)
    
    # Check for token generation
    token_patterns = [
        (r'(token|captcha|challenge|clearance)', "Token/challenge reference"),
        (r'(generate|calculate|compute)', "Computation terms"),
        (r'(Math\.random|crypto)', "Random generation"),
        (r'(cookie|setCookie|document\.cookie)', "Cookie manipulation"),
        (r'(xhr|XMLHttpRequest|fetch)', "Request sending")
    ]
    
    for pattern, description in token_patterns:
        if re.search(pattern, script, re.IGNORECASE):
            analysis["token_generation_indicators"].append(description)
    
    # Check for common obfuscation techniques
    if len(re.findall(r'eval\(', script)) > 3:
        analysis["obfuscation_level"] = "high"
    elif len(re.findall(r'\\x[0-9a-f]{2}', script)) > 10:
        analysis["obfuscation_level"] = "high"
    elif len(re.findall(r'String\.fromCharCode', script)) > 3:
        analysis["obfuscation_level"] = "high"
    elif re.search(r'function\(\w{1,2},\w{1,2},\w{1,2}\)\{', script):
        analysis["obfuscation_level"] = "medium"
    elif sum(1 for c in script if c == ';') > len(script) / 10:
        analysis["obfuscation_level"] = "medium"
    elif sum(len(w) > 30 for w in re.findall(r'\w+', script)) > 10:
        analysis["obfuscation_level"] = "medium"
    
    # Extract potential key function names
    function_pattern = r'function\s+(\w+)\s*\('
    functions = re.findall(function_pattern, script)
    
    suspicious_terms = ["challenge", "token", "captcha", "verify", "bot", "check", "security"]
    for func in functions:
        if any(term in func.lower() for term in suspicious_terms):
            analysis["key_functions"].append(func)
    
    # Determine if this is potentially protection-related
    analysis["potential_protection"] = (
        len(analysis["fingerprinting_indicators"]) > 2 or
        len(analysis["token_generation_indicators"]) > 2 or
        analysis["obfuscation_level"] != "none" or
        len(analysis["key_functions"]) > 0
    )
    
    return analysis

def analyze_cookies(headers: Dict[str, str]) -> List[Dict[str, Any]]:
    """
    Analyze cookies for common protection-related patterns.
    """
    cookie_header = headers.get("Cookie", "") or headers.get("Set-Cookie", "")
    if not cookie_header:
        return []
    
    # Split multiple cookies
    cookies = []
    for cookie_str in cookie_header.split(";"):
        parts = cookie_str.strip().split("=", 1)
        if len(parts) == 2:
            name, value = parts
            cookie = {
                "name": name.strip(),
                "value": value.strip() if len(value.strip()) < 50 else value.strip()[:50] + "... [truncated]",
                "protection_related": False,
                "vendor": "unknown"
            }
            
            # Check if this is a known protection cookie
            for vendor, signatures in BOT_PROTECTION_SIGNATURES.items():
                for sig in signatures:
                    if re.search(sig, name, re.IGNORECASE):
                        cookie["protection_related"] = True
                        cookie["vendor"] = vendor
                        break
                if cookie["protection_related"]:
                    break
            
            cookies.append(cookie)
    
    return cookies

def identify_protection_system(flow) -> List[Dict[str, Any]]:
    """
    Identify potential bot protection systems based on signatures.
    """
    protections = []
    
    # Combine all searchable content
    searchable_content = ""
    # Add request headers
    for k, v in flow.request.headers.items():
        searchable_content += f"{k}: {v}\n"
    
    # Check response if available
    if flow.response:
        # Add response headers
        for k, v in flow.response.headers.items():
            searchable_content += f"{k}: {v}\n"
        
        # Add response content if it's text
        content_type = flow.response.headers.get("Content-Type", "")
        if "text" in content_type or "javascript" in content_type or "json" in content_type:
            try:
                searchable_content += flow.response.content.decode('utf-8', errors='ignore')
            except Exception:
                pass
    
    # Check for protection signatures
    for vendor, signatures in BOT_PROTECTION_SIGNATURES.items():
        matches = []
        for sig in signatures:
            if re.search(sig, searchable_content, re.IGNORECASE):
                matches.append(sig)
        
        if matches:
            protections.append({
                "vendor": vendor,
                "confidence": len(matches) / len(signatures) * 100,
                "matching_signatures": matches
            })
    
    return sorted(protections, key=lambda x: x["confidence"], reverse=True)

def analyze_response_for_challenge(flow) -> Dict[str, Any]:
    """
    Analyze a response to determine if it contains a challenge.
    """
    if not flow.response:
        return {"is_challenge": False}
    
    result = {
        "is_challenge": False,
        "challenge_indicators": [],
        "status_code": flow.response.status_code,
        "challenge_type": "unknown"
    }
    
    # Check status code
    if flow.response.status_code in [403, 429, 503]:
        result["challenge_indicators"].append(f"Suspicious status code: {flow.response.status_code}")
    
    # Check for challenge headers
    challenge_headers = {
        "cf-mitigated": "Cloudflare mitigation",
        "cf-chl-bypass": "Cloudflare challenge bypass",
        "x-datadome": "DataDome protection",
        "x-px": "PerimeterX",
        "x-amz-captcha": "AWS WAF Captcha"
    }
    
    for header, description in challenge_headers.items():
        if any(h.lower() == header.lower() for h in flow.response.headers.keys()):
            result["challenge_indicators"].append(f"Challenge header: {description}")
    
    # Check for challenge content patterns
    content = flow.response.content.decode('utf-8', errors='ignore')
    challenge_patterns = [
        (r'captcha', "CAPTCHA"),
        (r'challenge', "Challenge term"),
        (r'blocked', "Blocking message"),
        (r'verify.*human', "Human verification"),
        (r'suspicious.*activity', "Suspicious activity message"),
        (r'security.*check', "Security check message"),
        (r'ddos', "DDoS protection message"),
        (r'automated.*request', "Automated request detection")
    ]
    
    for pattern, description in challenge_patterns:
        if re.search(pattern, content, re.IGNORECASE):
            result["challenge_indicators"].append(f"Content indicator: {description}")
    
    # Determine if this is a challenge response
    result["is_challenge"] = len(result["challenge_indicators"]) > 0
    
    # Determine challenge type
    if "CAPTCHA" in " ".join(result["challenge_indicators"]):
        result["challenge_type"] = "captcha"
    elif "JavaScript" in content and result["is_challenge"]:
        result["challenge_type"] = "javascript"
    elif result["is_challenge"]:
        result["challenge_type"] = "other"
    
    return result

async def analyze_protection(arguments: dict) -> list[types.TextContent]:
    """
    Analyze a flow for bot protection mechanisms and extract challenge details.
    """
    session_id = "mitmproxy.dump" #arguments.get("session_id")
    flow_index = arguments.get("flow_index")
    extract_scripts = arguments.get("extract_scripts", True)
    
    if not session_id:
        return [types.TextContent(type="text", text="Error: Missing session_id")]
    if flow_index is None:
        return [types.TextContent(type="text", text="Error: Missing flow_index")]
    
    try:
        flows = await get_flows_from_dump(session_id)
        
        try:
            flow = flows[flow_index]
            
            if flow.type != "http":
                return [types.TextContent(type="text", text=f"Error: Flow {flow_index} is not an HTTP flow")]
            
            # Analyze the flow for protection mechanisms
            analysis = {
                "flow_index": flow_index,
                "method": flow.request.method,
                "url": flow.request.url,
                "protection_systems": identify_protection_system(flow),
                "request_cookies": analyze_cookies(dict(flow.request.headers)),
                "has_response": flow.response is not None,
            }
            
            if flow.response:
                # Add response analysis
                content_type = flow.response.headers.get("Content-Type", "")
                is_html = "text/html" in content_type
                
                analysis.update({
                    "status_code": flow.response.status_code,
                    "response_cookies": analyze_cookies(dict(flow.response.headers)),
                    "challenge_analysis": analyze_response_for_challenge(flow),
                    "content_type": content_type,
                    "is_html": is_html,
                })
                
                # If HTML and script extraction is requested, extract and analyze JavaScript
                if is_html and extract_scripts:
                    try:
                        html_content = flow.response.content.decode('utf-8', errors='ignore')
                        analysis["scripts"] = extract_javascript(html_content)
                    except Exception as e:
                        analysis["script_extraction_error"] = str(e)
            
            # Add remediation suggestions based on findings
            analysis["suggestions"] = generate_suggestions(analysis)
            
            return [types.TextContent(type="text", text=json.dumps(analysis, indent=2))]
            
        except IndexError:
            return [types.TextContent(type="text", text=f"Error: Flow index {flow_index} out of range")]
            
    except FileNotFoundError:
        return [types.TextContent(type="text", text="Error: Session not found d")]
    except Exception as e:
        return [types.TextContent(type="text", text=f"Error analyzing protection: {str(e)}")]

def generate_suggestions(analysis: Dict[str, Any]) -> List[str]:
    """
    Generate remediation suggestions based on the protection analysis.
    """
    suggestions = []
    
    # Check if any protection system was detected
    if analysis.get("protection_systems"):
        top_system = analysis["protection_systems"][0]["vendor"]
        confidence = analysis["protection_systems"][0]["confidence"]
        
        if confidence > 50:
            suggestions.append(f"Detected {top_system} with {confidence:.1f}% confidence.")
            
            # Add vendor-specific suggestions
            if "Cloudflare" in top_system:
                suggestions.append("Cloudflare often uses JavaScript challenges. Check for cf_clearance cookie.")
                suggestions.append("Consider using proven techniques like cfscrape or cloudscraper libraries.")
            elif "Akamai" in top_system:
                suggestions.append("Akamai uses sensor_data for browser fingerprinting.")
                suggestions.append("Focus on _abck cookie which contains browser verification data.")
            elif "PerimeterX" in top_system:
                suggestions.append("PerimeterX relies on JavaScript execution and browser fingerprinting.")
                suggestions.append("Look for _px cookies which are essential for session validation.")
            elif "DataDome" in top_system:
                suggestions.append("DataDome uses advanced behavioral and fingerprinting techniques.")
                suggestions.append("The datadome cookie is critical for maintaining sessions.")
            elif "CAPTCHA" in top_system:
                suggestions.append("This site uses CAPTCHA challenges which may require manual solving or specialized services.")
    
    # Add suggestions based on challenge type
    if analysis.get("challenge_analysis", {}).get("is_challenge", False):
        challenge_type = analysis.get("challenge_analysis", {}).get("challenge_type", "unknown")
        
        if challenge_type == "javascript":
            suggestions.append("This response contains a JavaScript challenge that must be solved.")
            suggestions.append("Consider using a headless browser to execute the challenge JavaScript.")
            
            # If we have script analysis, add more specific suggestions
            if "scripts" in analysis:
                obfuscated_scripts = [s for s in analysis["scripts"] if s.get("summary", {}).get("obfuscation_level") in ["medium", "high"]]
                if obfuscated_scripts:
                    suggestions.append(f"Found {len(obfuscated_scripts)} obfuscated script(s) that likely contain challenge logic.")
                
                fingerprinting_scripts = [s for s in analysis["scripts"] if s.get("summary", {}).get("fingerprinting_indicators")]
                if fingerprinting_scripts:
                    techniques = set()
                    for script in fingerprinting_scripts:
                        techniques.update(script.get("summary", {}).get("fingerprinting_indicators", []))
                    suggestions.append(f"Detected browser fingerprinting techniques: {', '.join(techniques)}.")
                    
        elif challenge_type == "captcha":
            suggestions.append("This response contains a CAPTCHA challenge.")
            suggestions.append("Consider using a CAPTCHA solving service or manual intervention.")
    
    # Check for important cookies
    protection_cookies = [c for c in analysis.get("response_cookies", []) if c.get("protection_related")]
    if protection_cookies:
        cookie_names = [c["name"] for c in protection_cookies]
        suggestions.append(f"Important protection cookies to maintain: {', '.join(cookie_names)}.")
    
    # General suggestions
    if analysis.get("protection_systems") or analysis.get("challenge_analysis", {}).get("is_challenge", False):
        suggestions.append("General recommendations:")
        suggestions.append("- Maintain consistent User-Agent between requests")
        suggestions.append("- Preserve all cookies from the session")
        suggestions.append("- Add appropriate referer and origin headers")
        suggestions.append("- Consider adding delays between requests to avoid rate limiting")
        suggestions.append("- Use rotating IP addresses if available")
    
    return suggestions

@server.call_tool()
async def handle_call_tool(
    name: str, arguments: dict | None
) -> list[types.TextContent | types.ImageContent | types.EmbeddedResource]:
    """
    Handle tool execution requests.
    Delegates to specific functions based on the tool name.
    """
    if not arguments:
        raise ValueError("Missing arguments")

    if name == "list_flows":
        return await list_flows(arguments)
    elif name == "get_flow_details":
        return await get_flow_details(arguments)
    elif name == "extract_json_fields":
        return await extract_json_fields(arguments)
    elif name == "analyze_protection":
        return await analyze_protection(arguments)
    else:
        raise ValueError(f"Unknown tool: {name}")

async def main():
    # Run the server using stdin/stdout streams
    async with mcp.server.stdio.stdio_server() as (read_stream, write_stream):
        await server.run(
            read_stream,
            write_stream,
            InitializationOptions(
                server_name="mitmproxy-mcp",
                server_version="0.1.0",
                capabilities=server.get_capabilities(
                    notification_options=NotificationOptions(),
                    experimental_capabilities={},
                ),
            ),
        )
