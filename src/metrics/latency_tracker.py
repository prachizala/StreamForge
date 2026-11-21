"""
HLS Video Chunk Latency & Ingestion Tracker
Monitors chunk encoding latency percentiles across live media transcoding pipelines.
"""
from typing import List, Dict, Any
import math

class ChunkLatencyTracker:
    def __init__(self):
        self.latencies: List[float] = []

    def record_chunk_latency(self, latency_ms: float) -> None:
        if latency_ms >= 0:
            self.latencies.append(latency_ms)

    def compute_percentiles(self) -> Dict[str, float]:
        if not self.latencies:
            return {"p50": 0.0, "p95": 0.0, "p99": 0.0, "count": 0}
        s = sorted(self.latencies)
        n = len(s)
        def p(pct: float) -> float:
            idx = min(int(math.ceil(pct / 100.0 * n)) - 1, n - 1)
            return round(s[idx], 2)
        return {
            "p50": p(50),
            "p95": p(95),
            "p99": p(99),
            "count": n
        }

if __name__ == "__main__":
    tracker = ChunkLatencyTracker()
    for val in [10.0, 12.0, 15.0, 20.0, 50.0, 120.0]:
        tracker.record_chunk_latency(val)
    res = tracker.compute_percentiles()
    assert res["count"] == 6
    assert res["p50"] == 15.0
    print("Chunk latency tracker unit tests passed!")
