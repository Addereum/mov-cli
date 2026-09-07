from __future__ import annotations
from typing import TYPE_CHECKING

import os
import shutil
import subprocess
import unicodedata

from rich.console import Console

__all__ = ("Download",)

if TYPE_CHECKING:
    from .config import Config
    from .media import Multi, Single

from .logger import mov_cli_logger
from .utils.cookies import get_cookie_file
from devgoldyutils import LoggerAdapter

logger = LoggerAdapter(mov_cli_logger, "Downloader")
console = Console()

class Download():
    def __init__(self, config: Config) -> None:
        self.config = config

    def download(self, media: Multi | Single, subtitles: str = None) -> subprocess.Popen:
        title = unicodedata.normalize('NFKD', media.display_name).encode('ascii', 'ignore').decode('ascii').replace("/", " ") # normalize title

        download_dir = getattr(self.config, 'downloads_dir', self.config.download_location)
        file_path = os.path.join(download_dir, title + ".mp4")

        console.print(f"[bold green]Downloading {title}...[/bold green]")

        use_yt_dlp = self.config.use_yt_dlp

        if shutil.which("yt-dlp") is None:
            logger.warning("yt-dlp was not found, defaulting to ffmpeg!")
            use_yt_dlp = False

        elif media.audio_url is not None:
            logger.warning("Can't use yt-dlp as this media contains an audio url, defaulting to ffmpeg!")
            use_yt_dlp = False

        if use_yt_dlp:
            logger.info("Downloading via yt-dlp...")

            args = [
                "yt-dlp", 
                media.url, 
                "-o", 
                file_path,
                "--downloader", 
                "ffmpeg", 
                "--hls-use-mpegts",
                "--js-runtimes",
                "node",
                "--remote-components",
                "ejs:github"
            ]

            cookie_file = get_cookie_file(self.config)
            if cookie_file is not None:
                args.extend(["--cookies", cookie_file])

            if self.config.debug is False:
                args.append("--quiet")

            if media.referrer is not None:
                args.extend(["--add-header", f"Referer:{media.referrer}"])

        else:
            logger.info("Downloading via ffmpeg...")

            args = [
                "ffmpeg",
                "-n",
            ]

            cookie_file = get_cookie_file(self.config)
            if cookie_file is not None:
                args.extend(["-cookies", cookie_file])

            if media.referrer is not None:
                args.extend(["-headers", f"Referer: {media.referrer}"])

            args.extend(["-i", media.url])

            if media.audio_url:
                args.extend(["-i", media.audio_url])

            if subtitles:
                args.extend(["-vf", f"subtitle={subtitles}", file_path])
            else:
                args.extend(["-c", "copy", file_path])

        return subprocess.Popen(args)

