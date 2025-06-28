// Generate sparkles for background
document.addEventListener('DOMContentLoaded', function() {
    const sparkleContainer = document.getElementById('sparkle-container');
    const sparkleCount = 100;
    
    for (let i = 0; i < sparkleCount; i++) {
    const sparkle = document.createElement('div');
    sparkle.classList.add('sparkle');
    
    // Random position
    sparkle.style.left = `${Math.random() * 100}vw`;
    sparkle.style.top = `${Math.random() * 100}vh`;
    
    // Random size
    const size = 1 + Math.random() * 3;
    sparkle.style.width = `${size}px`;
    sparkle.style.height = `${size}px`;
    
    // Random animation delay
    sparkle.style.animationDelay = `${Math.random() * 12}s`;
    
    sparkleContainer.appendChild(sparkle);
    }
});
