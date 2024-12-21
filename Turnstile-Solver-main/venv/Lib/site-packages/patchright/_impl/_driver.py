import inspect
import os
import sys
from pathlib import Path
from typing import Tuple
import patchright
from patchright._repo_version import version


def compute_driver_executable() -> Tuple[str, str]:
    driver_path = Path(inspect.getfile(patchright)).parent / "driver"
    cli_path = str(driver_path / "package" / "cli.js")
    if sys.platform == "win32":
        return (
            os.getenv("PLAYWRIGHT_NODEJS_PATH", str(driver_path / "node.exe")),
            cli_path,
        )
    return (os.getenv("PLAYWRIGHT_NODEJS_PATH", str(driver_path / "node")), cli_path)


def get_driver_env() -> dict:
    env = os.environ.copy()
    env["PW_LANG_NAME"] = "python"
    env["PW_LANG_NAME_VERSION"] = f"{sys.version_info.major}.{sys.version_info.minor}"
    env["PW_CLI_DISPLAY_VERSION"] = version
    return env
