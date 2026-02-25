import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./src/**/*.{js,ts,jsx,tsx}"],
  darkMode: "media",
  theme: {
    extend: {
      colors: {
        surface: {
          primary: "var(--surface-primary)",
          secondary: "var(--surface-secondary)",
          tertiary: "var(--surface-tertiary)",
          elevated: "var(--surface-elevated)",
        },
        label: {
          primary: "var(--label-primary)",
          secondary: "var(--label-secondary)",
          tertiary: "var(--label-tertiary)",
          quaternary: "var(--label-quaternary)",
        },
        tint: {
          green: "var(--tint-green)",
          red: "var(--tint-red)",
          orange: "var(--tint-orange)",
          blue: "var(--tint-blue)",
          purple: "var(--tint-purple)",
          gold: "var(--tint-gold)",
        },
        separator: "var(--separator)",
        /* backwards compat */
        bg: "var(--surface-primary)",
        card: "var(--surface-secondary)",
        ink: "var(--label-primary)",
        accent: "var(--tint-blue)",
        muted: "var(--label-tertiary)",
      },
      fontSize: {
        "large-title": ["34px", { lineHeight: "41px", letterSpacing: "0.37px", fontWeight: "700" }],
        "title-1": ["28px", { lineHeight: "34px", letterSpacing: "0.36px", fontWeight: "700" }],
        "title-2": ["22px", { lineHeight: "28px", letterSpacing: "0.35px", fontWeight: "700" }],
        "title-3": ["20px", { lineHeight: "25px", letterSpacing: "0.38px", fontWeight: "600" }],
        headline: ["17px", { lineHeight: "22px", letterSpacing: "-0.41px", fontWeight: "600" }],
        body: ["17px", { lineHeight: "22px", letterSpacing: "-0.41px", fontWeight: "400" }],
        callout: ["16px", { lineHeight: "21px", letterSpacing: "-0.32px", fontWeight: "400" }],
        subheadline: ["15px", { lineHeight: "20px", letterSpacing: "-0.24px", fontWeight: "400" }],
        footnote: ["13px", { lineHeight: "18px", letterSpacing: "-0.08px", fontWeight: "400" }],
        "caption-1": ["12px", { lineHeight: "16px", fontWeight: "400" }],
        "caption-2": ["11px", { lineHeight: "13px", letterSpacing: "0.07px", fontWeight: "400" }],
      },
      borderRadius: {
        sm: "8px",
        md: "12px",
        lg: "16px",
        xl: "20px",
      },
      boxShadow: {
        "elevation-0": "none",
        "elevation-1": "0 1px 3px rgba(0,0,0,0.04), 0 1px 2px rgba(0,0,0,0.06)",
        "elevation-2": "0 4px 12px rgba(0,0,0,0.06), 0 1px 3px rgba(0,0,0,0.04)",
        "elevation-3": "0 8px 24px rgba(0,0,0,0.08), 0 2px 6px rgba(0,0,0,0.04)",
      },
      transitionTimingFunction: {
        apple: "cubic-bezier(0.25, 0.1, 0.25, 1)",
        "apple-spring": "cubic-bezier(0.34, 1.56, 0.64, 1)",
        "apple-decel": "cubic-bezier(0.0, 0.0, 0.2, 1)",
      },
      keyframes: {
        shimmer: {
          "0%": { backgroundPosition: "-200% 0" },
          "100%": { backgroundPosition: "200% 0" },
        },
        "fade-in": {
          "0%": { opacity: "0", transform: "translateY(8px)" },
          "100%": { opacity: "1", transform: "translateY(0)" },
        },
        "scale-in": {
          "0%": { opacity: "0", transform: "scale(0.95)" },
          "100%": { opacity: "1", transform: "scale(1)" },
        },
        "slide-up": {
          "0%": { opacity: "0", transform: "translateY(20px)" },
          "100%": { opacity: "1", transform: "translateY(0)" },
        },
        "pulse-soft": {
          "0%, 100%": { opacity: "1" },
          "50%": { opacity: "0.7" },
        },
      },
      animation: {
        shimmer: "shimmer 1.5s ease-in-out infinite",
        "fade-in": "fade-in 0.3s cubic-bezier(0,0,0.2,1) forwards",
        "scale-in": "scale-in 0.2s cubic-bezier(0.25,0.1,0.25,1) forwards",
        "slide-up": "slide-up 0.5s cubic-bezier(0,0,0.2,1) forwards",
        "pulse-soft": "pulse-soft 2s cubic-bezier(0.4,0,0.6,1) infinite",
      },
    },
  },
  plugins: [],
};

export default config;
