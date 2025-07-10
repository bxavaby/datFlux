document.addEventListener("DOMContentLoaded", function () {
  const navItems = document.querySelectorAll(".nav-item");
  const navTrack = document.querySelector(".nav-track");
  const navIndicator = document.querySelector(".nav-indicator");
  const mobileToggle = document.querySelector(".mobile-nav-toggle");

  setupSmoothScrolling();
  setupScrollSpy();
  setupMobileNav();
  updateNavIndicator(document.querySelector(".nav-item.active"));

  function setupSmoothScrolling() {
    navItems.forEach((item) => {
      item.addEventListener("click", function (e) {
        e.preventDefault();

        const targetId = this.getAttribute("href");
        const targetSection = document.querySelector(targetId);

        if (targetSection) {
          const headerOffset = document.querySelector(".neo-nav").offsetHeight;
          const targetPosition = targetSection.offsetTop - headerOffset;

          window.scrollTo({
            top: targetPosition,
            behavior: "smooth",
          });

          history.pushState(null, null, targetId);

          if (document.body.classList.contains("mobile-nav-open")) {
            document.body.classList.remove("mobile-nav-open");
          }

          navItems.forEach((nav) => nav.classList.remove("active"));
          this.classList.add("active");

          updateNavIndicator(this);
        }
      });
    });
  }

  function setupScrollSpy() {
    const sections = document.querySelectorAll("section");

    window.addEventListener(
      "scroll",
      debounce(() => {
        let current = "";

        sections.forEach((section) => {
          const sectionTop = section.offsetTop;
          const sectionHeight = section.offsetHeight;
          const headerHeight = document.querySelector(".neo-nav").offsetHeight;

          if (
            window.scrollY >= sectionTop - headerHeight - 20 &&
            window.scrollY < sectionTop + sectionHeight - headerHeight
          ) {
            current = section.getAttribute("id");
          }
        });

        if (current) {
          navItems.forEach((item) => {
            item.classList.remove("active");
            if (item.getAttribute("href") === `#${current}`) {
              item.classList.add("active");
              updateNavIndicator(item);
            }
          });
        }
      }, 50),
    );
  }

  function setupMobileNav() {
    if (mobileToggle) {
      mobileToggle.addEventListener("click", () => {
        document.body.classList.toggle("mobile-nav-open");
      });
    }
  }

  function updateNavIndicator(activeItem) {
    if (navIndicator && activeItem && window.innerWidth > 768) {
      const itemRect = activeItem.getBoundingClientRect();
      const trackRect = navTrack.getBoundingClientRect();

      navIndicator.style.width = `${itemRect.width}px`;
      navIndicator.style.left = `${itemRect.left - trackRect.left}px`;
    }
  }

  window.addEventListener(
    "resize",
    debounce(() => {
      const activeItem = document.querySelector(".nav-item.active");
      if (activeItem) {
        updateNavIndicator(activeItem);
      }
    }, 100),
  );

  function debounce(func, wait) {
    let timeout;
    return function () {
      const context = this;
      const args = arguments;
      clearTimeout(timeout);
      timeout = setTimeout(() => func.apply(context, args), wait);
    };
  }

  if (window.location.hash) {
    const targetItem = document.querySelector(
      `.nav-item[href="${window.location.hash}"]`,
    );
    if (targetItem) {
      setTimeout(() => {
        targetItem.click();
      }, 100);
    }
  }
});
