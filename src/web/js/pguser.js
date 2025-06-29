document.addEventListener('DOMContentLoaded', function() {
    // Generate sparkles for background
    const sparkleContainer = document.getElementById('sparkle-container');
    const sparkleCount = 60;
    
    for (let i = 0; i < sparkleCount; i++) {
    const sparkle = document.createElement('div');
    sparkle.classList.add('sparkle');
    sparkle.style.left = `${Math.random() * 100}vw`;
    sparkle.style.top = `${Math.random() * 100}vh`;
    const size = 1 + Math.random() * 2;
    sparkle.style.width = `${size}px`;
    sparkle.style.height = `${size}px`;
    sparkle.style.animationDelay = `${Math.random() * 10}s`;
    sparkleContainer.appendChild(sparkle);
    }
    
    // Copy functionality
    const copyButtons = document.querySelectorAll('.copy-btn');
    copyButtons.forEach(button => {
    button.addEventListener('click', function() {
        const text = this.getAttribute('data-copy');
        navigator.clipboard.writeText(text);
        
        // Visual feedback
        const icon = this.querySelector('i');
        icon.className = 'fas fa-check';
        this.classList.add('copied');
        
        // Reset after 2 seconds
        setTimeout(() => {
        icon.className = 'fas fa-copy';
        this.classList.remove('copied');
        }, 2000);
    });
    });
    
    // Toggle password visibility
    const togglePassword = document.getElementById('togglePassword');
    const password = document.getElementById('password');
    const passwordStrength = document.getElementById('passwordStrength');
    
    if (togglePassword && password) {
    togglePassword.addEventListener('click', function() {
        const type = password.getAttribute('type') === 'password' ? 'text' : 'password';
        password.setAttribute('type', type);
        this.querySelector('i').classList.toggle('fa-eye');
        this.querySelector('i').classList.toggle('fa-eye-slash');
    });
    }
    
    // Password strength indicator
    if (password && passwordStrength) {
    password.addEventListener('input', function() {
        const strength = calculatePasswordStrength(this.value);
        passwordStrength.style.width = strength.percentage + '%';
        passwordStrength.className = 'password-strength-fill ' + strength.class;
    });
    }
    
    function calculatePasswordStrength(password) {
    let strength = 0;
    
    // Length check
    if (password.length >= 8) strength += 25;
    if (password.length >= 12) strength += 25;
    
    // Character variety
    if (/[A-Z]/.test(password)) strength += 15;
    if (/[a-z]/.test(password)) strength += 15;
    if (/[0-9]/.test(password)) strength += 10;
    if (/[^A-Za-z0-9]/.test(password)) strength += 10;
    
    // Classify strength
    let strengthClass = '';
    if (strength < 50) {
        strengthClass = '';
    } else if (strength < 75) {
        strengthClass = 'medium';
    } else {
        strengthClass = 'strong';
    }
    
    return {
        percentage: strength,
        class: strengthClass
    };
    }
});
