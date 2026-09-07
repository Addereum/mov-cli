import pytest
from mov_cli.media.media import Single, Multi
from mov_cli.media.metadata import Metadata, MetadataType
from mov_cli.media.quality import Quality
from mov_cli.utils.episode_selector import EpisodeSelector
from devgoldyutils import Colours

def test_single_display_name():
    single1 = Single(url="http", title="Test Movie")
    assert single1.display_name == "Test Movie"
    
    single2 = Single(url="http", title="Test Movie", year="2023")
    assert single2.display_name == "Test Movie (2023)"

def test_multi_display_name():
    multi = Multi(url="http", title="Test Show", episode=EpisodeSelector(episode=5, season=2))
    assert multi.display_name == "Test Show - S2 EP5"

def test_metadata_display_name():
    meta1 = Metadata(id="1", title="Test Meta", type=MetadataType.SINGLE)
    assert meta1.display_name == f"{Colours.BLUE}Test Meta{Colours.RESET}"
    
    meta2 = Metadata(id="2", title="Test Meta 2", type=MetadataType.MULTI, year="2022")
    assert meta2.display_name == f"{Colours.PINK_GREY}Test Meta 2{Colours.RESET} (2022)"

def test_metadata_type_enum():
    assert MetadataType.MULTI.value == 0
    assert MetadataType.SINGLE.value == 1

def test_quality_enum():
    assert Quality.SD.value == 480
    assert Quality.HD.value == 720
    assert Quality.FHD.value == 1080
    assert Quality.QHD.value == 1440
    assert Quality.UHD.value == 2160
    assert Quality._2K.value == 1440
    assert Quality._4K.value == 2160
    assert Quality.AUTO.value == 0
    
    assert Quality.HD.apply_p() == "720p"
