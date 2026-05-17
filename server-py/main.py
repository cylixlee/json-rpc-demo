import asyncio
import sys

from ajsonrpc.dispatcher import Dispatcher
from ajsonrpc.manager import AsyncJSONRPCResponseManager

dispatcher = Dispatcher()


@dispatcher.add_function
def pythonic_hello(name: str) -> str:
    return f">>> Hello, {name}"


async def main() -> None:
    print("[SERVER PY] Starting JSON-RPC server...")
    loop = asyncio.get_event_loop()
    manager = AsyncJSONRPCResponseManager(dispatcher)
    while True:
        try:
            # accept one line from stdin
            line = await loop.run_in_executor(None, sys.stdin.readline)
            if not line:
                break

            # strip JSON-RPC request
            request = line.strip()
            if not request:
                continue

            # JSON-RPC response
            response = await manager.get_payload_for_payload(request)
            sys.stdout.write(response + "\n")
            sys.stdout.flush()
        except KeyboardInterrupt, asyncio.CancelledError:
            break
    print("[SERVER PY] Bye~")


if __name__ == "__main__":
    asyncio.run(main())
