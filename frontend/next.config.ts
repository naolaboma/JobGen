import type { NextConfig } from "next";
import path from "path";

const nextConfig: NextConfig = {
  // Explicitly set Turbopack root to this frontend app to avoid multi-lockfile root inference warnings
  turbopack: {
    root: path.resolve(__dirname),
  },
};

export default nextConfig;
