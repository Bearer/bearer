;(function () {
  const themeToggleButton = document.getElementById("dark-light-toggle")
  // Dark/light mode setup
  if (
    localStorage.theme === "dark" ||
    (!("theme" in localStorage) &&
      window.matchMedia("(prefers-color-scheme: dark)").matches)
  ) {
    document.documentElement.classList.add("dark")
    localStorage.theme = "dark"
  } else {
    document.documentElement.classList.remove("dark")
    localStorage.theme = "light"
  }

  themeToggleButton.addEventListener("click", () => {
    toggleTheme()
  })

  function toggleTheme() {
    if ("theme" in localStorage && localStorage.theme === "dark") {
      localStorage.theme = "light"
      document.documentElement.classList.remove("dark")
    } else {
      localStorage.theme = "dark"
      document.documentElement.classList.add("dark")
    }
  }

  // mobile nav open/close
  const navButton = document.getElementById("toggle-nav")

  navButton.addEventListener("click", () => {
    document.querySelector("#doc-nav").classList.toggle("open")
  })

  // toc open/close mobile
  const toggleToc = document.getElementById("js-toggle-toc")
  toggleToc.addEventListener("click", () => {
    document.querySelector("#toc-container").classList.toggle("open")
    toggleToc.querySelector("div").classList.toggle("flip")
  })
})()
