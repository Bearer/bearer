module.exports = {
  content: ["_site/**/*.html"],
  safelist: [],
  theme: {
    colors: {
      transparent: "transparent",
      current: "currentColor",
      black: "#0C0C0C",
      white: "#FFFFFF",
      neutral: {
        500: "#3A3A3A",
        400: "#969696",
        300: "#C4C4C4",
        200: "#E9E9E9",
        100: "#F5F5F5",
      },
      main: {
        100: "#F1E9FD",
        200: "#D4BCF8",
        300: "#A46EF6",
        400: "#6E20E6",
        DEFAULT: "#6E20E6",
        500: "#42138A",
        600: "#230A49",
      },
    },
    extend: {
      typography: (theme) => ({
        DEFAULT: {
          css: {
            "--tw-prose-body": theme("colors.black"),
            p: {},
            "code::before": {
              content: "normal",
            },
            "code::after": {
              content: "normal",
            },
          },
          invert: {
            css: {
              "--tw-prose-invert-body": theme("colors.neutral.200"),
              // "--tw-prose-invert-headings": theme("colors.neutral.200"),
              p: {},
            },
          },
        },
      }),
    },
  },
  plugins: [require("@tailwindcss/typography")],
};
