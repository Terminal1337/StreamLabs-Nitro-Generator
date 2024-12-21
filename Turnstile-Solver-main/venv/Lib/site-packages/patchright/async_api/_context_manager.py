import asyncio
from typing import Any
from patchright._impl._connection import Connection
from patchright._impl._object_factory import create_remote_object
from patchright._impl._transport import PipeTransport
from patchright.async_api._generated import Playwright as AsyncPlaywright


class PlaywrightContextManager:

    def __init__(self) -> None:
        self._connection: Connection
        self._exit_was_called = False

    async def __aenter__(self) -> AsyncPlaywright:
        loop = asyncio.get_running_loop()
        self._connection = Connection(
            None, create_remote_object, PipeTransport(loop), loop
        )
        loop.create_task(self._connection.run())
        playwright_future = self._connection.playwright_future
        done, _ = await asyncio.wait(
            {self._connection._transport.on_error_future, playwright_future},
            return_when=asyncio.FIRST_COMPLETED,
        )
        if not playwright_future.done():
            playwright_future.cancel()
        playwright = AsyncPlaywright(next(iter(done)).result())
        playwright.stop = self.__aexit__
        return playwright

    async def start(self) -> AsyncPlaywright:
        return await self.__aenter__()

    async def __aexit__(self, *args: Any) -> None:
        if self._exit_was_called:
            return
        self._exit_was_called = True
        await self._connection.stop_async()
