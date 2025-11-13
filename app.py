from llm_routing import get_tool_selection
from typing import Dict, List, Optional

def main(task_description: str, tools: List[Dict[str, str]]) -> None:
    """
    Main function to orchestrate the tool routing process.

    :param task_description: A description of the task.
    :param tools: A list of available tools.
    """
    print(f"Task: {task_description}")
    selected_tool = get_tool_selection(task_description, tools)

    if selected_tool:
        print(f"Selected tool: {selected_tool}")
    else:
        print("No suitable tool was found.")

if __name__ == "__main__":
    available_tools = [
        {
            "name": "code-analyzer",
            "description": "Analyzes code for potential bugs and vulnerabilities.",
        },
        {
            "name": "web-scraper",
            "description": "Scrapes data from websites.",
        },
    ]
    task = "Analyze the python code for potential security vulnerabilities"
    main(task, available_tools)