import asyncio
from typing import TYPE_CHECKING, Any, Optional, cast
from greenlet import greenlet
from patchright._impl._connection import ChannelOwner, Connection
from patchright._impl._errors import Error
from patchright._impl._greenlets import MainGreenlet
from patchright._impl._object_factory import create_remote_object
from patchright._impl._playwright import Playwright
from patchright._impl._transport import PipeTransport
from patchright.sync_api._generated import Playwright as SyncPlaywright

if TYPE_CHECKING:
    from asyncio.unix_events import AbstractChildWatcher


class PlaywrightContextManager:

    def __init__(self) -> None:
        self._playwright: SyncPlaywright
        self._loop: asyncio.AbstractEventLoop
        self._own_loop = False
        self._watcher: Optional[AbstractChildWatcher] = None
        self._exit_was_called = False

    def __enter__(self) -> SyncPlaywright:
        try:
            self._loop = asyncio.get_running_loop()
        except RuntimeError:
            self._loop = asyncio.new_event_loop()
            self._own_loop = True
        if self._loop.is_running():
            raise Error(
                "It looks like you are using Playwright Sync API inside the asyncio loop.\nPlease use the Async API instead."
            )

        def greenlet_main() -> None:
            self._loop.run_until_complete(self._connection.run_as_sync())

        dispatcher_fiber = MainGreenlet(greenlet_main)
        self._connection = Connection(
            dispatcher_fiber,
            create_remote_object,
            PipeTransport(self._loop),
            self._loop,
        )
        g_self = greenlet.getcurrent()

        def callback_wrapper(channel_owner: ChannelOwner) -> None:
            playwright_impl = cast(Playwright, channel_owner)
            self._playwright = SyncPlaywright(playwright_impl)
            g_self.switch()

        self._connection.call_on_object_with_known_name("Playwright", callback_wrapper)
        dispatcher_fiber.switch()
        playwright = self._playwright
        playwright.stop = self.__exit__
        return playwright

    def start(self) -> SyncPlaywright:
        return self.__enter__()

    def __exit__(self, *args: Any) -> None:
        if self._exit_was_called:
            return
        self._exit_was_called = True
        self._connection.stop_sync()
        if self._watcher:
            self._watcher.close()
        if self._own_loop:
            tasks = asyncio.all_tasks(self._loop)
            for t in [t for t in tasks if not (t.done() or t.cancelled())]:
                t.cancel()
            self._loop.run_until_complete(self._loop.shutdown_asyncgens())
            self._loop.close()
