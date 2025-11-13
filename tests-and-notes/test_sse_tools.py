import asyncio
import json
import aiohttp
from aiohttp_sse_client import client as sse_client

async def get_tools_list():
    session_id = f"session-{int(asyncio.get_event_loop().time())}-{hash(str(asyncio.get_event_loop())) % 10000}"
    
    async with sse_client.EventSource('http://localhost:3020/sse', session_id=session_id) as source:
        # Send tools/list request
        request = {
            "jsonrpc": "2.0",
            "method": "tools/list",
            "id": 1,
            "params": {"sessionId": session_id}
        }
        
        async with aiohttp.ClientSession() as http_session:
            async with http_session.post('http://localhost:3020/sse', 
                                       json=request,
                                       headers={'Content-Type': 'application/json'}) as resp:
                if resp.status == 200:
                    response = await resp.text()
                    print("Tools List Response:")
                    print(response)
                else:
                    print(f"Error: {resp.status}")
                    print(await resp.text())

# Run the async function
asyncio.run(get_tools_list())
