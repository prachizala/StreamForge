/// Normalize raw PCM samples to fixed decibel scale
pub fn normalize_samples(samples: &[i16], target_peak: f32) -> Vec<f32> {
    let mut normalized = Vec::with_capacity(samples.len());
    let max_abs = samples.iter().map(|&s| s.abs()).max().unwrap_or(1) as f32;
    let scale = if max_abs > 0.0 { target_peak / max_abs } else { 1.0 };

    for &s in samples {
        normalized.push((s as f32 * scale).clamp(-1.0, 1.0));
    }
    normalized
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_normalize_samples() {
        let input = vec![1000, -2000, 3000];
        let out = normalize_samples(&input, 1.0);
        assert_eq!(out.len(), 3);
        assert!((out[2] - 1.0).abs() < 1e-4);
    }
}
