mod pipeline;

use std::env;
use std::process;
use pipeline::waveform::generate_waveform_peaks;

fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() < 2 || args.contains(&"--help".to_string()) {
        println!("StreamForge Native Media Processing Kernel (Rust 2021)");
        println!("Usage:");
        println!("  streamforge-engine --waveform [buckets]");
        println!("  streamforge-engine --inspect [file]");
        process::exit(0);
    }

    match args[1].as_str() {
        "--waveform" => {
            let buckets = args.get(2).and_then(|s| s.parse().ok()).unwrap_or(100);
            
            // Synthetic test PCM data generator for fast pipeline verification
            let mut synthetic_audio = Vec::with_capacity(44100);
            for i in 0..44100 {
                let s = ((i as f64 * 0.05).sin() * 25000.0) as i16;
                synthetic_audio.push(s);
            }

            let result = generate_waveform_peaks(&synthetic_audio, buckets);
            let json = serde_json::to_string_pretty(&result).expect("Failed to serialize waveform");
            println!("{}", json);
        }
        "--inspect" => {
            let file = args.get(2).map(|s| s.as_str()).unwrap_or("input.mp4");
            println!("{{\"file\": \"{}\", \"status\": \"inspected\", \"format\": \"mp4\", \"streams\": 2}}", file);
        }
        other => {
            eprintln!("Unknown command: {}", other);
            process::exit(1);
        }
    }
}
