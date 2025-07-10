document.addEventListener("DOMContentLoaded", function () {
  createEntropyParticles();

  animateFooterEntropy();

  generateBinaryPattern();
});

function createEntropyParticles() {
  const container = document.querySelector(".entropy-particles");
  if (!container) return;

  for (let i = 0; i < 50; i++) {
    const particle = document.createElement("div");

    const size = Math.random() * 3 + 1;
    const posX = Math.random() * 100;
    const posY = Math.random() * 100;
    const duration = Math.random() * 10 + 10;
    const delay = Math.random() * 5;

    particle.style.position = "absolute";
    particle.style.width = `${size}px`;
    particle.style.height = `${size}px`;
    particle.style.backgroundColor = getRandomColor();
    particle.style.left = `${posX}%`;
    particle.style.top = `${posY}%`;
    particle.style.opacity = `${Math.random() * 0.5 + 0.2}`;
    particle.style.borderRadius = "50%";

    particle.style.animation = `float ${duration}s linear infinite`;
    particle.style.animationDelay = `-${delay}s`;

    container.appendChild(particle);
  }

  addKeyframeAnimation();
}

function getRandomColor() {
  const colors = [
    "var(--accent-blue)",
    "var(--accent-purple)",
    "var(--accent-cyan)",
    "var(--accent-green)",
  ];
  return colors[Math.floor(Math.random() * colors.length)];
}

function addKeyframeAnimation() {
  const styleSheet = document.createElement("style");
  styleSheet.textContent = `
    @keyframes float {
      0% {
        transform: translateY(0) translateX(0);
        opacity: 0;
      }
      10% {
        opacity: 0.7;
      }
      90% {
        opacity: 0.5;
      }
      100% {
        transform: translateY(-100vh) translateX(20px);
        opacity: 0;
      }
    }
  `;
  document.head.appendChild(styleSheet);
}

function animateFooterEntropy() {}

function generateBinaryPattern() {
  const binaryLayer = document.querySelector(".binary-layer");
  if (!binaryLayer) return;

  setInterval(() => {
    if (Math.random() > 0.7) {
      binaryLayer.style.opacity = "0.05";
      setTimeout(() => {
        binaryLayer.style.opacity = "0.03";
      }, 200);
    }
  }, 3000);
}
