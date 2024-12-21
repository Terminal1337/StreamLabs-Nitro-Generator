from typing import Dict, cast
from patchright._impl._artifact import Artifact
from patchright._impl._browser import Browser
from patchright._impl._browser_context import BrowserContext
from patchright._impl._browser_type import BrowserType
from patchright._impl._cdp_session import CDPSession
from patchright._impl._connection import ChannelOwner
from patchright._impl._dialog import Dialog
from patchright._impl._element_handle import ElementHandle
from patchright._impl._fetch import APIRequestContext
from patchright._impl._frame import Frame
from patchright._impl._js_handle import JSHandle
from patchright._impl._local_utils import LocalUtils
from patchright._impl._network import (
    Request,
    Response,
    Route,
    WebSocket,
    WebSocketRoute,
)
from patchright._impl._page import BindingCall, Page, Worker
from patchright._impl._playwright import Playwright
from patchright._impl._selectors import SelectorsOwner
from patchright._impl._stream import Stream
from patchright._impl._tracing import Tracing
from patchright._impl._writable_stream import WritableStream


class DummyObject(ChannelOwner):

    def __init__(
        self, parent: ChannelOwner, type: str, guid: str, initializer: Dict
    ) -> None:
        super().__init__(parent, type, guid, initializer)


def create_remote_object(
    parent: ChannelOwner, type: str, guid: str, initializer: Dict
) -> ChannelOwner:
    if type == "Artifact":
        return Artifact(parent, type, guid, initializer)
    if type == "APIRequestContext":
        return APIRequestContext(parent, type, guid, initializer)
    if type == "BindingCall":
        return BindingCall(parent, type, guid, initializer)
    if type == "Browser":
        return Browser(cast(BrowserType, parent), type, guid, initializer)
    if type == "BrowserType":
        return BrowserType(parent, type, guid, initializer)
    if type == "BrowserContext":
        return BrowserContext(parent, type, guid, initializer)
    if type == "CDPSession":
        return CDPSession(parent, type, guid, initializer)
    if type == "Dialog":
        return Dialog(parent, type, guid, initializer)
    if type == "ElementHandle":
        return ElementHandle(parent, type, guid, initializer)
    if type == "Frame":
        return Frame(parent, type, guid, initializer)
    if type == "JSHandle":
        return JSHandle(parent, type, guid, initializer)
    if type == "LocalUtils":
        local_utils = LocalUtils(parent, type, guid, initializer)
        if not local_utils._connection._local_utils:
            local_utils._connection._local_utils = local_utils
        return local_utils
    if type == "Page":
        return Page(parent, type, guid, initializer)
    if type == "Playwright":
        return Playwright(parent, type, guid, initializer)
    if type == "Request":
        return Request(parent, type, guid, initializer)
    if type == "Response":
        return Response(parent, type, guid, initializer)
    if type == "Route":
        return Route(parent, type, guid, initializer)
    if type == "Stream":
        return Stream(parent, type, guid, initializer)
    if type == "Tracing":
        return Tracing(parent, type, guid, initializer)
    if type == "WebSocket":
        return WebSocket(parent, type, guid, initializer)
    if type == "WebSocketRoute":
        return WebSocketRoute(parent, type, guid, initializer)
    if type == "Worker":
        return Worker(parent, type, guid, initializer)
    if type == "WritableStream":
        return WritableStream(parent, type, guid, initializer)
    if type == "Selectors":
        return SelectorsOwner(parent, type, guid, initializer)
    return DummyObject(parent, type, guid, initializer)
