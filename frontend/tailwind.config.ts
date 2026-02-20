import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        bg: "#f5f6f8",
        card: "#ffffff",
        ink: "#101418",
        accent: "#0f766e",
        muted: "#64748b"
      }
    }
  },
  plugins: []
};

export default config;
