(function () {
  const filterTriggers = document.querySelectorAll(".js-filter-button");
  const containers = document.querySelectorAll(".js-filter-container");
  const langCounter = document.querySelector(".js-lang-count");
  const owaspCounter = document.querySelector(".js-owasp-count");
  const ruleCounter = document.querySelector(".js-rule-count");
  filterTriggers.forEach((button) => {
    button.addEventListener("click", (e) => {
      const parent = button.parentElement;
      if (parent.classList.contains("filter-open")) {
        parent.classList.remove("filter-open");
      } else {
        containers.forEach((e) => {
          e.classList.remove("filter-open");
        });
        parent.classList.add("filter-open");
      }
    });
  });

  function resetForm(search) {
    search.value = "";
    document.querySelectorAll(".filter-toggle:checked").forEach((item) => {
      item.checked = false;
    });
    langCounter.innerHTML = "All";
    owaspCounter.innerHTML = "All";
    resetButton.disabled = true;
    filterResults();
  }
  function compare(item, query, filters) {
    let source = item.innerHTML.toLowerCase();
    let filtersActive = filters.length === 0 ? false : true;
    let found = false;
    if (filtersActive) {
      filters.forEach((f) => {
        if (source.includes(f.value.toLowerCase())) {
          found = true;
        }
      });
      if (source.includes(query) && found === true) {
        return true;
      }
    } else {
      if (source.includes(query)) {
        return true;
      }
      return false;
    }
  }

  function updateCounts(langs, owasp) {
    if (langs.length === 0) {
      langCounter.innerHTML = "All";
    } else if (langs.length > 0) {
      langCounter.innerHTML = langs.length;
    }
    if (owasp.length === 0) {
      owaspCounter.innerHTML = "All";
    } else if (owasp.length > 0) {
      owaspCounter.innerHTML = owasp.length;
    }
  }

  function filterResults() {
    updateURL();
    const rules = document.querySelectorAll(".js-rule");
    const langs = document.querySelectorAll(
      "#lang-filters .filter-toggle:checked"
    );
    const owasp = document.querySelectorAll(
      "#owasp-filters .filter-toggle:checked"
    );
    let ruleCount = 0;
    updateCounts(langs, owasp);
    const query = document.getElementById("search-input").value.toLowerCase();
    rules.forEach((rule) => {
      if (compare(rule, query, langs) && compare(rule, query, owasp)) {
        rule.classList.remove("hidden");

        ruleCount++;
      } else {
        if (!rule.classList.contains("hidden")) {
          rule.classList.add("hidden");
        }
      }
    });
    ruleCounter.innerHTML = ruleCount;
    if (
      document.querySelectorAll(".filter-toggle:checked").length ||
      query.length
    ) {
      resetButton.disabled = false;
    } else {
      resetButton.disabled = true;
    }
  }

  const input = document.getElementById("search-input");
  const checkboxes = document.querySelectorAll(".filter-toggle");
  const form = document.getElementById("rule-search");
  const resetButton = document.querySelector(".js-filter-reset");
  let timer;
  const delay = 300;
  form.addEventListener("submit", (e) => {
    e.preventDefault();
  });
  input.addEventListener("keyup", (e) => {
    clearTimeout(timer);
    timer = setTimeout(filterResults, delay);
  });

  checkboxes.forEach((checkbox) => {
    checkbox.addEventListener("change", (e) => {
      filterResults();
    });
  });

  resetButton.addEventListener("click", (e) => {
    resetForm(input);
  });

  function updateURL() {
    let params = new URLSearchParams();
    let checkedBoxes = document.querySelectorAll(".filter-toggle:checked");
    if (input.value.length > 0) {
      params.append(input.name, input.value);
    }
    checkedBoxes.forEach((checkbox) => {
      params.append(checkbox.name, checkbox.value);
    });
    if (params.size > 0) {
      window.history.replaceState({}, "", decodeURIComponent(`?${params}`));
    } else {
      window.history.replaceState(history.state, "", window.location.pathname);
    }
  }

  function load() {
    let params = new URLSearchParams(window.location.search);
    for (const [key, value] of params.entries()) {
      let input = document.querySelector(`input[name="${key}"]`);
      if (input.type === "search") {
        input.value = value;
      } else if (input.type === "checkbox") {
        input.checked = true;
      }
    }
    filterResults();
  }
  load();
})();
