(function () {
  const themeToggleButton = document.getElementById("dark-light-toggle");
  // Dark/light mode setup
  if (
    localStorage.theme === "dark" ||
    (!("theme" in localStorage) &&
      window.matchMedia("(prefers-color-scheme: dark)").matches)
  ) {
    document.documentElement.classList.add("dark");
    localStorage.theme = "dark";
  } else {
    document.documentElement.classList.remove("dark");
    localStorage.theme = "light";
  }

  themeToggleButton.addEventListener("click", () => {
    toggleTheme();
  });

  function toggleTheme() {
    if ("theme" in localStorage && localStorage.theme === "dark") {
      localStorage.theme = "light";
      document.documentElement.classList.remove("dark");
    } else {
      localStorage.theme = "dark";
      document.documentElement.classList.add("dark");
    }
  }
})();
