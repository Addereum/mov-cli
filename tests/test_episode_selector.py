import pytest
from mov_cli.utils.episode_selector import EpisodeSelector

def test_episode_selector_basic():
    selector = EpisodeSelector()
    assert selector.episode == 1
    assert selector.season == 1
    
    selector = EpisodeSelector(episode=5, season=3)
    assert selector.episode == 5
    assert selector.season == 3

def test_next_season():
    selector = EpisodeSelector(episode=10, season=2)
    selector._next_season()
    assert selector.episode == 1
    assert selector.season == 3

def test_previous_season():
    selector = EpisodeSelector(episode=5, season=4)
    # Mocking media_episodes dictionary {season: episode_count}
    media_episodes = {3: 12}
    
    selector._previous_season(media_episodes)
    assert selector.season == 3
    assert selector.episode == 12
    
    # If not in media_episodes, defaults to 1
    selector._previous_season({2: 8})
    assert selector.season == 2
    assert selector.episode == 8
    
    selector._previous_season({})
    assert selector.season == 1
    assert selector.episode == 1
