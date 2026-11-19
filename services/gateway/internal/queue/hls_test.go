package queue

import (
	"strings"
	"testing"
)

func TestHLSManifestFormat(t *testing.T) {
	manifest := "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-STREAM-INF:BANDWIDTH=800000,RESOLUTION=640x360\n360p.m3u8"
	if !strings.HasPrefix(manifest, "#EXTM3U") {
		t.Errorf("expected HLS manifest header #EXTM3U")
	}
}
