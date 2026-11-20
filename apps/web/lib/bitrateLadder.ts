/**
 * Adaptive Bitrate (ABR) Ladder Definition and Validation Module.
 * Standardizes multi-bitrate HLS output profiles for creator video rendering.
 */

export interface BitrateProfile {
  name: string;
  resolution: string;
  targetBitrateKbps: number;
  framerate: number;
}

export const STANDARD_ABR_LADDER: BitrateProfile[] = [
  { name: '1080p_FHD', resolution: '1920x1080', targetBitrateKbps: 4500, framerate: 60 },
  { name: '720p_HD', resolution: '1280x720', targetBitrateKbps: 2400, framerate: 30 },
  { name: '480p_SD', resolution: '854x480', targetBitrateKbps: 1000, framerate: 30 },
  { name: '360p_Mobile', resolution: '640x360', targetBitrateKbps: 600, framerate: 30 },
];

/**
 * Recommends optimal initial playback profile given measured connection speed in Kbps.
 */
export function getRecommendedProfile(connectionSpeedKbps: number): BitrateProfile {
  if (connectionSpeedKbps >= 5000) return STANDARD_ABR_LADDER[0];
  if (connectionSpeedKbps >= 2500) return STANDARD_ABR_LADDER[1];
  if (connectionSpeedKbps >= 1200) return STANDARD_ABR_LADDER[2];
  return STANDARD_ABR_LADDER[3];
}
