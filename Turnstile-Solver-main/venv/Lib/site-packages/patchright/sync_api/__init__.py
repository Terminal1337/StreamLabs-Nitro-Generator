"""
Python package `playwright` is a Python library to automate Chromium,
Firefox and WebKit with a single API. Playwright is built to enable cross-browser
web automation that is ever-green, capable, reliable and fast.
"""

from typing import Any, Optional, Union, overload
import patchright._impl._api_structures
import patchright._impl._errors
import patchright.sync_api._generated
from patchright._impl._assertions import (
    APIResponseAssertions as APIResponseAssertionsImpl,
)
from patchright._impl._assertions import LocatorAssertions as LocatorAssertionsImpl
from patchright._impl._assertions import PageAssertions as PageAssertionsImpl
from patchright.sync_api._context_manager import PlaywrightContextManager
from patchright.sync_api._generated import (
    Accessibility,
    APIRequest,
    APIRequestContext,
    APIResponse,
    APIResponseAssertions,
    Browser,
    BrowserContext,
    BrowserType,
    CDPSession,
    ConsoleMessage,
    Dialog,
    Download,
    ElementHandle,
    FileChooser,
    Frame,
    FrameLocator,
    JSHandle,
    Keyboard,
    Locator,
    LocatorAssertions,
    Mouse,
    Page,
    PageAssertions,
    Playwright,
    Request,
    Response,
    Route,
    Selectors,
    Touchscreen,
    Video,
    WebSocket,
    WebSocketRoute,
    Worker,
)

ChromiumBrowserContext = BrowserContext
Cookie = patchright._impl._api_structures.Cookie
FilePayload = patchright._impl._api_structures.FilePayload
FloatRect = patchright._impl._api_structures.FloatRect
Geolocation = patchright._impl._api_structures.Geolocation
HttpCredentials = patchright._impl._api_structures.HttpCredentials
PdfMargins = patchright._impl._api_structures.PdfMargins
Position = patchright._impl._api_structures.Position
ProxySettings = patchright._impl._api_structures.ProxySettings
ResourceTiming = patchright._impl._api_structures.ResourceTiming
SourceLocation = patchright._impl._api_structures.SourceLocation
StorageState = patchright._impl._api_structures.StorageState
ViewportSize = patchright._impl._api_structures.ViewportSize
Error = patchright._impl._errors.Error
TimeoutError = patchright._impl._errors.TimeoutError


def sync_playwright() -> PlaywrightContextManager:
    return PlaywrightContextManager()


class Expect:
    _unset: Any = object()

    def __init__(self) -> None:
        self._timeout: Optional[float] = None

    def set_options(self, timeout: Optional[float] = _unset) -> None:
        """
        This method sets global `expect()` options.

        Args:
            timeout (float): Timeout value in milliseconds. Default to 5000 milliseconds.

        Returns:
            None
        """
        if timeout is not self._unset:
            self._timeout = timeout

    @overload
    def __call__(
        self, actual: Page, message: Optional[str] = None
    ) -> PageAssertions: ...

    @overload
    def __call__(
        self, actual: Locator, message: Optional[str] = None
    ) -> LocatorAssertions: ...

    @overload
    def __call__(
        self, actual: APIResponse, message: Optional[str] = None
    ) -> APIResponseAssertions: ...

    def __call__(
        self, actual: Union[Page, Locator, APIResponse], message: Optional[str] = None
    ) -> Union[PageAssertions, LocatorAssertions, APIResponseAssertions]:
        if isinstance(actual, Page):
            return PageAssertions(
                PageAssertionsImpl(actual._impl_obj, self._timeout, message=message)
            )
        elif isinstance(actual, Locator):
            return LocatorAssertions(
                LocatorAssertionsImpl(actual._impl_obj, self._timeout, message=message)
            )
        elif isinstance(actual, APIResponse):
            return APIResponseAssertions(
                APIResponseAssertionsImpl(
                    actual._impl_obj, self._timeout, message=message
                )
            )
        raise ValueError(f"Unsupported type: {type(actual)}")


expect = Expect()
__all__ = [
    "expect",
    "Accessibility",
    "APIRequest",
    "APIRequestContext",
    "APIResponse",
    "Browser",
    "BrowserContext",
    "BrowserType",
    "CDPSession",
    "ChromiumBrowserContext",
    "ConsoleMessage",
    "Cookie",
    "Dialog",
    "Download",
    "ElementHandle",
    "Error",
    "FileChooser",
    "FilePayload",
    "FloatRect",
    "Frame",
    "FrameLocator",
    "Geolocation",
    "HttpCredentials",
    "JSHandle",
    "Keyboard",
    "Locator",
    "Mouse",
    "Page",
    "PdfMargins",
    "Position",
    "Playwright",
    "ProxySettings",
    "Request",
    "ResourceTiming",
    "Response",
    "Route",
    "Selectors",
    "SourceLocation",
    "StorageState",
    "sync_playwright",
    "TimeoutError",
    "Touchscreen",
    "Video",
    "ViewportSize",
    "WebSocket",
    "WebSocketRoute",
    "Worker",
]
