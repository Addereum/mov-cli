from __future__ import annotations
from typing import TYPE_CHECKING, TypedDict

if TYPE_CHECKING:
    from typing import List, Optional
    from .media import Metadata
    from .utils import EpisodeSelector
    from .utils.platform import SUPPORTED_PLATFORMS

import json
import time
from datetime import datetime
from devgoldyutils import LoggerAdapter, Colours

from .logger import mov_cli_logger
from .utils.paths import get_cache_directory
from .media import MetadataType

class HistoryEntry(TypedDict):
    title: str
    scraper_name: str
    metadata_id: str
    metadata_type: str
    episode: Optional[str]
    timestamp: float
    image_url: Optional[str]

logger = LoggerAdapter(
    mov_cli_logger, prefix = Colours.BLUE.apply("History")
)

class WatchHistory:
    """Manages persistent watch history for mov-cli."""
    
    def __init__(self, platform: SUPPORTED_PLATFORMS):
        self.history_file = get_cache_directory(platform) / "history.json"
        self.max_entries = 100

    def _load_history(self) -> List[HistoryEntry]:
        if not self.history_file.exists():
            return []
        try:
            with self.history_file.open("r", encoding="utf-8") as f:
                return json.load(f)
        except (json.JSONDecodeError, OSError) as e:
            logger.warning(f"Failed to load history: {e}")
            return []

    def _save_history(self, history: List[HistoryEntry]) -> None:
        try:
            with self.history_file.open("w", encoding="utf-8") as f:
                json.dump(history, f, indent=4)
        except OSError as e:
            logger.warning(f"Failed to save history: {e}")

    def add_entry(self, metadata: Metadata, episode: EpisodeSelector, scraper_name: str) -> None:
        history = self._load_history()
        
        episode_str = None
        if metadata.type == MetadataType.MULTI:
            episode_str = f"S{episode.season}E{episode.episode}"

        entry: HistoryEntry = {
            "title": metadata.title,
            "scraper_name": scraper_name,
            "metadata_id": str(metadata.id),
            "metadata_type": metadata.type.name if hasattr(metadata.type, "name") else str(metadata.type),
            "episode": episode_str,
            "timestamp": time.time(),
            "image_url": getattr(metadata, "image_url", None)
        }

        # Remove existing entry for the same media if it exists
        history = [h for h in history if h.get("metadata_id") != str(metadata.id) or h.get("scraper_name") != scraper_name]

        # Add to the beginning
        history.insert(0, entry)

        # Prune if too large
        if len(history) > self.max_entries:
            history = history[:self.max_entries]

        self._save_history(history)
        logger.debug(f"Added '{metadata.title}' to watch history.")

    def get_history(self, limit: int = 20) -> List[HistoryEntry]:
        history = self._load_history()
        return history[:limit]

    def clear_history(self) -> None:
        if self.history_file.exists():
            self.history_file.unlink()
        logger.info("Watch history cleared.")
