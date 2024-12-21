from typing import Optional


def is_target_closed_error(error: Exception) -> bool:
    return isinstance(error, TargetClosedError)


class Error(Exception):

    def __init__(self, message: str) -> None:
        self._message = message
        self._name: Optional[str] = None
        self._stack: Optional[str] = None
        super().__init__(message)

    @property
    def message(self) -> str:
        return self._message

    @property
    def name(self) -> Optional[str]:
        return self._name

    @property
    def stack(self) -> Optional[str]:
        return self._stack


class TimeoutError(Error):
    pass


class TargetClosedError(Error):

    def __init__(self, message: str = None) -> None:
        super().__init__(message or "Target page, context or browser has been closed")


def rewrite_error(error: Exception, message: str) -> Exception:
    rewritten_exc = type(error)(message)
    if isinstance(rewritten_exc, Error) and isinstance(error, Error):
        rewritten_exc._name = error.name
        rewritten_exc._stack = error.stack
    return rewritten_exc
