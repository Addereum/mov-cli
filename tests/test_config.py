import pytest
from mov_cli.config import Config
from mov_cli.media.quality import Quality

def test_config_defaults():
    config = Config(override_config={})
    assert config.player == "mpv"
    assert config.debug is False
    assert config.hide_ip is True
    assert config.fzf_enabled in (True, False)

def test_override_config():
    config = Config(override_config={"player": "vlc", "debug": True, "hide_ip": False})
    assert config.player == "vlc"
    assert config.debug is True
    assert config.hide_ip is False

def test_resolution_quality_parsing():
    config_auto = Config(override_config={"quality": "auto"})
    assert config_auto.resolution == Quality.AUTO

    config_hd = Config(override_config={"quality": "hd"})
    assert config_hd.resolution == Quality.HD

    config_720 = Config(override_config={"quality": {"resolution": 720}})
    assert config_720.resolution == Quality.HD

    config_1080 = Config(override_config={"quality": {"resolution": 1080}})
    assert config_1080.resolution == Quality.FHD

def test_subtitle_language():
    config_subtitle = Config(override_config={"subtitle": {"language": "fr"}})
    assert config_subtitle.language.iso639_1 == "fr"

    config_subtitles = Config(override_config={"subtitles": {"language": "es"}})
    assert config_subtitles.language.iso639_1 == "es"

def test_player_args():
    config_dict = Config(override_config={"player": {"binary": "mpv", "args": ["--fs"], "args_override": True}})
    assert config_dict.player == "mpv"
    assert config_dict.player_args == ["--fs"]
    assert config_dict.player_args_override is True

    config_str = Config(override_config={"player": "vlc"})
    assert config_str.player == "vlc"
    assert config_str.player_args == []
    assert config_str.player_args_override is False

def test_fzf_enabled(monkeypatch):
    import shutil
    monkeypatch.setattr(shutil, "which", lambda x: "path/to/fzf" if x == "fzf" else None)
    assert Config(override_config={}).fzf_enabled is True

    monkeypatch.setattr(shutil, "which", lambda x: None)
    assert Config(override_config={}).fzf_enabled is False
