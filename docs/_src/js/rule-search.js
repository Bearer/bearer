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
  }

  const input = document.getElementById("search-input");
  const checkboxes = document.querySelectorAll(".filter-toggle");
  const form = document.getElementById("rule-search");
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
})();
