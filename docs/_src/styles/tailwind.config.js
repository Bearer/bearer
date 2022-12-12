module.exports = {
  content: ["_site/**/*.html"],
  safelist: [],
  darkMode: "class",
  theme: {
    colors: {
      transparent: "transparent",
      current: "currentColor",
      black: "#0C0C0C",
      white: "#FFFFFF",
      code: "hsl(243,27%,35%)",
      neutral: {
        600: "hsl(222, 47%, 11%)",
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
        700: "#1F0F38",
      },
    },
    extend: {
      spacing: {
        "heading-offset": "7rem",
      },
      typography: (theme) => ({
        DEFAULT: {
          css: {
            "--tw-prose-body": theme("colors.black"),
            a: {
              color: theme("colors.main.400"),
              textDecoration: "none",
              "&:hover": {
                textDecoration: "underline",
              },
            },
            p: {},
            "code::before": {
              content: "normal",
            },
            "code::after": {
              content: "normal",
            },
          },
        },
        invert: {
          css: {
            p: {},
            a: {
              color: theme("colors.main.300"),
              textDecoration: "none",
              "&:hover": {
                textDecoration: "underline",
              },
            },
          },
        },
      }),
    },
  },
  plugins: [require("@tailwindcss/typography")],
};
