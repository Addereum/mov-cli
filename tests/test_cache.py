import pytest
import time
from mov_cli.cache import Cache

@pytest.fixture
def mock_cache_dir(tmp_path, monkeypatch):
    monkeypatch.setattr("mov_cli.cache.get_cache_directory", lambda _: tmp_path)
    return tmp_path

def test_set_get_cache(mock_cache_dir):
    cache = Cache("Windows")
    cache.set_cache("key1", "value1")
    assert cache.get_cache("key1") == "value1"

def test_cache_expiration(mock_cache_dir):
    cache = Cache("Windows")
    cache.set_cache("expiring_key", "value", seconds_until_expired=1)
    assert cache.get_cache("expiring_key") == "value"
    # Sleep to ensure it expires
    time.sleep(1.1)
    assert cache.get_cache("expiring_key") is None

def test_clear_cache(mock_cache_dir):
    cache = Cache("Windows")
    cache.set_cache("key1", "value1")
    cache.set_cache("key2", "value2")
    
    cache.clear_cache("key1")
    assert cache.get_cache("key1") is None
    assert cache.get_cache("key2") == "value2"
    
    cache.clear_all_cache()
    assert cache.get_cache("key2") is None

def test_section_based_caching(mock_cache_dir):
    cache1 = Cache("Windows", section="sec1")
    cache2 = Cache("Windows", section="sec2")
    
    cache1.set_cache("key1", "value_sec1")
    cache2.set_cache("key1", "value_sec2")
    
    assert cache1.get_cache("key1") == "value_sec1"
    assert cache2.get_cache("key1") == "value_sec2"
    
    cache1.clear_all_cache()
    assert cache1.get_cache("key1") is None
    assert cache2.get_cache("key1") == "value_sec2"
