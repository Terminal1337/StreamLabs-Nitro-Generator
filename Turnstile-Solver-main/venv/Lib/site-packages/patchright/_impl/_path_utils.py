import inspect
from pathlib import Path


def get_file_dirname() -> Path:
    """Returns the callee (`__file__`) directory name"""
    frame = inspect.stack()[1]
    module = inspect.getmodule(frame[0])
    assert module
    assert module.__file__
    return Path(module.__file__).parent.absolute()
