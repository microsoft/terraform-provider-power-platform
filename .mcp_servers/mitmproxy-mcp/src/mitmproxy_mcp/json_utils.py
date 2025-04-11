from typing import Any, Dict, List, Union

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
