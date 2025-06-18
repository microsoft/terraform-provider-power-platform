import os
from pathlib import Path
from typing import Union, Optional
from autogen_core import CancellationToken
from autogen_core.tools import FunctionTool

async def save_file(file_path: Union[str, Path], content: Union[str, bytes], 
              encoding: str = 'utf-8') -> Optional[str]:
    """
    Saves content to a file at the specified path.
    
    Args:
        file_path: Path where the file should be saved
        content: Content to write to the file (string or bytes). This will overwrite any existing file content.
        encoding: File encoding (default: utf-8, ignored if content is bytes)
    
    Returns:
        None if successful, error message string if failed
    
    Note:
        This function requires that the directory structure already exists.
        It will not create missing directories.
    """
    try:
        # Convert to Path object if string was provided
        path_obj = Path(file_path)
        
        # Check if parent directory exists
        if not path_obj.parent.exists():
            return f"Error: Directory '{path_obj.parent}' does not exist. Please create the directory first."
        
        # Check if parent is actually a directory
        if not path_obj.parent.is_dir():
            return f"Error: '{path_obj.parent}' exists but is not a directory."
        
        # Write content to file
        write_mode = 'w' if isinstance(content, str) else 'wb'
        
        if write_mode == 'w':
            with open(path_obj, write_mode, encoding=encoding) as f:
                f.write(content)
        else:
            with open(path_obj, write_mode) as f:
                f.write(content)
                
        return None  # Success
    
    except PermissionError:
        return f"Error: Permission denied when writing to '{file_path}'. Check file permissions."
    except IsADirectoryError:
        return f"Error: '{file_path}' is a directory, not a file."
    except Exception as e:
        return f"Error saving file '{file_path}': {str(e)}"

execute_save_file_tool = FunctionTool(
    save_file,
    description="Saves content to a file at the specified path. Returns None if successful, or an error message string if failed."
)
