import asyncio
from typing import Dict, Optional, cast
from pyee.asyncio import AsyncIOEventEmitter
from patchright._impl._connection import Channel
from patchright._impl._errors import TargetClosedError
from patchright._impl._helper import Error, ParsedMessagePayload
from patchright._impl._transport import Transport


class JsonPipeTransport(AsyncIOEventEmitter, Transport):

    def __init__(self, loop: asyncio.AbstractEventLoop, pipe_channel: Channel) -> None:
        super().__init__(loop)
        Transport.__init__(self, loop)
        self._stop_requested = False
        self._pipe_channel = pipe_channel

    def request_stop(self) -> None:
        self._stop_requested = True
        self._pipe_channel.send_no_reply("close", {})

    def dispose(self) -> None:
        self.on_error_future.cancel()
        self._stopped_future.cancel()

    async def wait_until_stopped(self) -> None:
        await self._stopped_future

    async def connect(self) -> None:
        self._stopped_future: asyncio.Future = asyncio.Future()

        def handle_message(message: Dict) -> None:
            if self._stop_requested:
                return
            self.on_message(cast(ParsedMessagePayload, message))

        def handle_closed(reason: Optional[str]) -> None:
            self.emit("close", reason)
            if reason:
                self.on_error_future.set_exception(TargetClosedError(reason))
            self._stopped_future.set_result(None)

        self._pipe_channel.on(
            "message", lambda params: handle_message(params["message"])
        )
        self._pipe_channel.on(
            "closed", lambda params: handle_closed(params.get("reason"))
        )

    async def run(self) -> None:
        await self._stopped_future

    def send(self, message: Dict) -> None:
        if self._stop_requested:
            raise Error("Playwright connection closed")
        self._pipe_channel.send_no_reply("send", {"message": message})
