document.addEventListener("DOMContentLoaded", function () {
  initSections();

  initTabs();

  initClipboard();

  generateEntropyVisualizations();

  initParanoiaMode();
});

function initSections() {
  const sections = document.querySelectorAll(".section");
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add("visible");
        }
      });
    },
    { threshold: 0.1 },
  );

  sections.forEach((section) => {
    observer.observe(section);
  });
}

function initTabs() {
  const tabs = document.querySelectorAll(".tab");
  tabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      const tabTarget = tab.getAttribute("data-tab");

      document
        .querySelectorAll(".tab")
        .forEach((t) => t.classList.remove("active"));
      document
        .querySelectorAll(".tab-content")
        .forEach((c) => c.classList.remove("active"));

      tab.classList.add("active");
      document.getElementById(`${tabTarget}-results`).classList.add("active");
    });
  });
}

function initClipboard() {
  if (typeof ClipboardJS !== "undefined") {
    const clipboard = new clipboard(".copy-btn");

    clipboard.on("success", function (e) {
      const button = e.trigger;
      button.classList.add("copied");
      button.textContent = "COPIED!";

      setTimeout(() => {
        button.classList.remove("copied");
        button.textContent = "COPY";
      }, 2000);

      e.clearSelection();
    });
  } else {
    document.querySelectorAll(".copy-btn").forEach((btn) => {
      btn.addEventListener("click", () => {
        const targetId = btn.getAttribute("data-clipboard-target");
        const text = document.querySelector(targetId).textContent;

        const textarea = document.createElement("textarea");
        textarea.value = text;
        textarea.setAttribute("readonly", "");
        textarea.style.position = "absolute";
        textarea.style.left = "-9999px";
        document.body.appendChild(textarea);
        textarea.select();
        document.execCommand("copy");
        document.body.removeChild(textarea);

        // Update button UI
        btn.classList.add("copied");
        btn.textContent = "COPIED!";

        setTimeout(() => {
          btn.classList.remove("copied");
          btn.textContent = "COPY";
        }, 2000);
      });
    });
  }
}

function generateEntropyVisualizations() {
  setInterval(() => {
    document.querySelectorAll(".binary-bit").forEach((bit) => {
      if (Math.random() > 0.7) {
        bit.textContent = bit.textContent === "0" ? "1" : "0";
      }
    });
  }, 1000);

  setInterval(() => {
    const hexChars = "0123456789ABCDEF";
    document.querySelectorAll(".hex-byte").forEach((byte) => {
      if (Math.random() > 0.8) {
        byte.textContent =
          hexChars[Math.floor(Math.random() * 16)] +
          hexChars[Math.floor(Math.random() * 16)];
      }
    });
  }, 2000);

  animateMeters();
}

function animateMeters() {
  const cpuMeter = document.querySelector(".cpu-meter");
  const ramMeter = document.querySelector(".ram-meter");
  const networkMeter = document.querySelector(".network-meter");

  setInterval(() => {
    if (cpuMeter) cpuMeter.style.width = `${75 + Math.random() * 20}%`;
    if (ramMeter) ramMeter.style.width = `${50 + Math.random() * 20}%`;
    if (networkMeter) networkMeter.style.width = `${35 + Math.random() * 20}%`;
  }, 1500);
}

function initParanoiaMode() {
  const securitySection = document.getElementById("security");
  if (securitySection) {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          document.body.classList.add("paranoia-mode");
        } else {
          document.body.classList.remove("paranoia-mode");
        }
      },
      { threshold: 0.3 },
    );

    observer.observe(securitySection);
  }
}
