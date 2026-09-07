from __future__ import annotations
from typing import TYPE_CHECKING, List

if TYPE_CHECKING:
    from typing import Optional

    from ..media import Media
    from ..utils.platform import SUPPORTED_PLATFORMS, Literal

import subprocess
from devgoldyutils import Colours

from ..errors import ReferrerNotSupportedError

from .player import Player

__all__ = ("MPV",)

class MPV(Player):
    def __init__(
        self, 
        platform: SUPPORTED_PLATFORMS, 
        args: Optional[List[str]] = None, 
        args_override: bool = False, 
        debug: bool = False, 
        **kwargs
    ) -> None:
        super().__init__(
            platform = platform, 
            args = args,
            debug = debug, 
            args_override = args_override
        )

    @property
    def display_name(self) -> str:
        return Colours.PURPLE.apply("MPV")

    def _get_args(self, platform: Literal["Windows", "Linux", "Android", "Darwin"], media: Media):
        if platform == "Android":
            return []

        args = [
            f"--force-media-title={media.display_name}",
        ]

        if media.referrer is not None:
            args.append(f"--referrer={media.referrer}")

        if media.subtitles is not None:

            for subtitle in media.subtitles:
                args.append(f"--sub-file={subtitle}")

        if self.debug is False:
            args.append("--no-terminal")

        args = self.handle_additional_args(args, self.args)

        return args

    def play(self, media: Media) -> Optional[subprocess.Popen]:
        """Plays this media in the MPV media player."""

        if self.platform == "Android":

            if media.referrer is not None:
                raise ReferrerNotSupportedError(
                    "The MPV player on Android does not support passing referrers, so this media cannot be played. :("
                )

            return subprocess.Popen(
                [
                    "am",
                    "start",
                    "-n",
                    "is.xyz.mpv/is.xyz.mpv.MPVActivity",
                    "-e",
                    "filepath",
                    media.url,
                ]
            )

        elif self.platform == "Linux" or self.platform == "Windows" or self.platform == "Darwin" or self.platform == "FreeBSD":
            mpv_executable = "mpv"
            if self.platform == "Windows":
                import shutil
                import os
                if shutil.which("mpv") is None:
                    fallback_paths = [
                        os.path.expanduser(r"~\Desktop\mov\mpv\mpv.exe"),
                        r"C:\Program Files\mpv\mpv.exe",
                        r"C:\Program Files (x86)\mpv\mpv.exe",
                        os.path.expandvars(r"%LOCALAPPDATA%\Microsoft\WindowsApps\mpv.exe")
                    ]
                    for p in fallback_paths:
                        if os.path.exists(p):
                            mpv_executable = p
                            break

            if "youtube.com" in media.url or "youtu.be" in media.url:
                yt_dlp_process = subprocess.Popen(
                    ["yt-dlp", "-o", "-", media.url, "--js-runtimes", "node", "--remote-components", "ejs:github"],
                    stdout=subprocess.PIPE,
                    stderr=subprocess.DEVNULL
                )
                mpv_process = subprocess.Popen(
                    [mpv_executable, "-"] + self._get_args(self.platform, media),
                    stdin=yt_dlp_process.stdout,
                    stdout=(subprocess.STDOUT if self.debug else subprocess.DEVNULL),
                    stderr=subprocess.STDOUT
                )
                yt_dlp_process.stdout.close()
                return mpv_process

            default_args = [
                mpv_executable, 
                media.url
            ]

            if media.audio_url is not None:
                default_args.append(f"--audio-file={media.audio_url}")

            return subprocess.Popen(
                default_args + self._get_args(self.platform, media),
                stdout = (subprocess.STDOUT if self.debug else subprocess.DEVNULL), # NOTE: https://github.com/mov-cli/mov-cli/issues/361
                stderr = subprocess.STDOUT
            )

        return None