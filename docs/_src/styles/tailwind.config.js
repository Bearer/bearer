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
        500: "#3D3D3D",
        400: "#939393",
        300: "#D4D4D4",
        200: "#EAEAEA",
        100: "#F6F6F6",
      },
      main: {
        100: "#EEE8FD",
        200: "#DBD0FB",
        300: "#9472E2",
        400: "#7043ED",
        DEFAULT: "#7043ED",
        500: "#4C14E9",
        600: "#3D10BA",
        700: "#1E065F",
      },
    },
    extend: {
      spacing: {
        "heading-offset": "6rem",
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
}
