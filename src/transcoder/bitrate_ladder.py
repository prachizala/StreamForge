"""
Adaptive Bitrate (ABR) Ladder Planner
Generates optimized HLS renditions balancing video quality and egress network bandwidth.
"""
from typing import List, Dict, Any

class BitrateLadderPlanner:
    STANDARD_PROFILES = [
        {"resolution": "1920x1080", "bitrate_kbps": 4500, "fps": 60},
        {"resolution": "1280x720", "bitrate_kbps": 2400, "fps": 30},
        {"resolution": "854x480", "bitrate_kbps": 1200, "fps": 30},
        {"resolution": "640x360", "bitrate_kbps": 600, "fps": 30},
    ]

    def plan_renditions(self, max_bandwidth_kbps: int) -> List[Dict[str, Any]]:
        return [p for p in self.STANDARD_PROFILES if p["bitrate_kbps"] <= max_bandwidth_kbps]

if __name__ == "__main__":
    planner = BitrateLadderPlanner()
    renditions = planner.plan_renditions(3000)
    assert len(renditions) == 3
    assert renditions[0]["resolution"] == "1280x720"
    print("Bitrate ladder planner unit tests passed!")
