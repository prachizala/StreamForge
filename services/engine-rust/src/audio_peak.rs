//! Audio peak amplitude detector for dynamic normalization and waveform visualization.

/// Detects peak positive and negative amplitudes across PCM audio samples.
#[derive(Debug, Clone, PartialEq)]
pub struct AudioPeak {
    pub max_amplitude: f32,
    pub min_amplitude: f32,
    pub sample_count: usize,
}

impl AudioPeak {
    /// Scans a slice of 32-bit floating point audio samples to compute peak levels.
    pub fn compute(samples: &[f32]) -> Self {
        if samples.is_empty() {
            return Self {
                max_amplitude: 0.0,
                min_amplitude: 0.0,
                sample_count: 0,
            };
        }

        let mut max_val = f32::MIN;
        let mut min_val = f32::MAX;

        for &sample in samples {
            if sample > max_val {
                max_val = sample;
            }
            if sample < min_val {
                min_val = sample;
            }
        }

        Self {
            max_amplitude: max_val,
            min_amplitude: min_val,
            sample_count: samples.len(),
        }
    }

    /// Returns the absolute peak level in normalized [0.0, 1.0] range.
    pub fn normalized_peak(&self) -> f32 {
        self.max_amplitude.abs().max(self.min_amplitude.abs()).min(1.0)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_peak_computation() {
        let samples = vec![0.1, -0.85, 0.42, 0.76, -0.3];
        let peak = AudioPeak::compute(&samples);
        assert_eq!(peak.sample_count, 5);
        assert_eq!(peak.max_amplitude, 0.76);
        assert_eq!(peak.min_amplitude, -0.85);
        assert_eq!(peak.normalized_peak(), 0.85);
    }

    #[test]
    fn test_empty_samples() {
        let peak = AudioPeak::compute(&[]);
        assert_eq!(peak.sample_count, 0);
        assert_eq!(peak.normalized_peak(), 0.0);
    }
}
