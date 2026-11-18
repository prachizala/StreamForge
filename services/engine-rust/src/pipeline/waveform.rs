use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WaveformResult {
    pub duration_seconds: f64,
    pub sample_rate: u32,
    pub total_peaks: usize,
    pub peaks: Vec<f32>,
}

/// Compute Root Mean Square (RMS) audio waveform peaks from raw PCM audio samples.
pub fn generate_waveform_peaks(samples: &[i16], buckets: usize) -> WaveformResult {
    if samples.is_empty() || buckets == 0 {
        return WaveformResult {
            duration_seconds: 0.0,
            sample_rate: 44100,
            total_peaks: 0,
            peaks: Vec::new(),
        };
    }

    let chunk_size = (samples.len() as f64 / buckets as f64).ceil() as usize;
    let mut peaks = Vec::with_capacity(buckets);

    for chunk in samples.chunks(chunk_size) {
        if chunk.is_empty() {
            continue;
        }

        // Calculate RMS amplitude
        let sum_sq: f64 = chunk
            .iter()
            .map(|&s| {
                let norm = s as f64 / 32768.0;
                norm * norm
            })
            .sum();

        let rms = (sum_sq / chunk.len() as f64).sqrt() as f32;
        // Clamp between 0.0 and 1.0
        peaks.push(rms.clamp(0.0, 1.0));
    }

    WaveformResult {
        duration_seconds: samples.len() as f64 / 44100.0,
        sample_rate: 44100,
        total_peaks: peaks.len(),
        peaks,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_empty_samples() {
        let res = generate_waveform_peaks(&[], 50);
        assert_eq!(res.total_peaks, 0);
        assert!(res.peaks.is_empty());
    }

    #[test]
    fn test_waveform_peak_normalization() {
        // Generate synthetic sine-wave burst
        let mut samples = Vec::new();
        for i in 0..1000 {
            let val = ((i as f64 * 0.1).sin() * 30000.0) as i16;
            samples.push(val);
        }

        let res = generate_waveform_peaks(&samples, 10);
        assert_eq!(res.total_peaks, 10);
        for peak in res.peaks {
            assert!(peak >= 0.0 && peak <= 1.0);
        }
    }
}
