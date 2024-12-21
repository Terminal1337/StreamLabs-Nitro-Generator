import os
from typing import Tuple
import greenlet


def _greenlet_trace_callback(
    event: str, args: Tuple[greenlet.greenlet, greenlet.greenlet]
) -> None:
    if event in ("switch", "throw"):
        origin, target = args
        print(f"Transfer from {origin} to {target} with {event}")


if os.environ.get("INTERNAL_PW_GREENLET_DEBUG"):
    greenlet.settrace(_greenlet_trace_callback)


class MainGreenlet(greenlet.greenlet):

    def __str__(self) -> str:
        return "<MainGreenlet>"


class RouteGreenlet(greenlet.greenlet):

    def __str__(self) -> str:
        return "<RouteGreenlet>"


class LocatorHandlerGreenlet(greenlet.greenlet):

    def __str__(self) -> str:
        return "<LocatorHandlerGreenlet>"


class EventGreenlet(greenlet.greenlet):

    def __str__(self) -> str:
        return "<EventGreenlet>"
