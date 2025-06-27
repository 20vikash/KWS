const sparkleContainer = document.getElementById('sparkle-container');

for (let i = 0; i < 160; i++) {
    const sparkle = document.createElement('div');
    sparkle.classList.add('sparkle');
    sparkle.style.left = Math.random() * 100 + 'vw';
    sparkle.style.top = Math.random() * 100 + 'vh';
    sparkle.style.animationDuration = 6 + Math.random() * 12 + 's';
    sparkle.style.animationDelay = Math.random() * 10 + 's';
    sparkleContainer.appendChild(sparkle);
}
