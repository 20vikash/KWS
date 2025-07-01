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
    
    // DOM elements
    const deployBtn = document.getElementById('deploy-btn');
    const killBtn = document.getElementById('kill-btn');
    const stopBtn = document.getElementById('stop-btn');
    const codeBtn = document.getElementById('code-btn');
    const deployModal = document.getElementById('deploy-modal');
    const closeModal = document.getElementById('close-modal');
    const cancelDeploy = document.getElementById('cancel-deploy');
    const confirmDeploy = document.getElementById('confirm-deploy');
    const instanceDetails = document.getElementById('instance-details');
    const emptyState = document.getElementById('empty-state');
    const statusBadge = document.getElementById('status-badge');
    const statusText = document.getElementById('status-text');
    
    // Hidden field to track credential requirement
    const credentialsRequired = document.createElement('input');
    credentialsRequired.type = 'hidden';
    credentialsRequired.id = 'credentials-required';
    credentialsRequired.value = 'no';
    document.body.appendChild(credentialsRequired);
    
    // Open deploy modal
    deployBtn.addEventListener('click', function() {
    if (credentialsRequired.value === 'no') {
        deployModal.classList.remove('hidden');
    } else {
        // Simulate direct deployment
        startAction('deploy');
    }
    });
    
    // Close modal
    closeModal.addEventListener('click', function() {
    deployModal.classList.add('hidden');
    });
    
    cancelDeploy.addEventListener('click', function() {
    deployModal.classList.add('hidden');
    });
    
    // Confirm deployment
    confirmDeploy.addEventListener('click', function() {
    deployModal.classList.add('hidden');
    startAction('deploy');
    });
    
    // Kill instance
    killBtn.addEventListener('click', function() {
    startAction('kill');
    });
    
    // Stop instance
    stopBtn.addEventListener('click', function() {
    startAction('stop');
    });
    
    // Code button
    codeBtn.addEventListener('click', function() {
    alert('Opening VS Code in browser...');
    // In a real app, this would open the code editor
    });
    
    // Start action (deploy, kill, stop)
    function startAction(action) {
    // Add blinking effect to action buttons
    deployBtn.classList.add('action-blinking');
    killBtn.classList.add('action-blinking');
    stopBtn.classList.add('action-blinking');
    
    // Disable buttons during action
    deployBtn.disabled = true;
    killBtn.disabled = true;
    stopBtn.disabled = true;
    
    // Simulate API call delay
    setTimeout(function() {
        // Remove blinking effect
        deployBtn.classList.remove('action-blinking');
        killBtn.classList.remove('action-blinking');
        stopBtn.classList.remove('action-blinking');
        
        // Update UI based on action
        if (action === 'deploy') {
        // Show instance details
        instanceDetails.classList.remove('hidden');
        emptyState.classList.add('hidden');
        
        // Update status
        statusBadge.className = 'status-badge status-active';
        statusText.textContent = 'Instance Active';
        
        // Enable/disable buttons
        deployBtn.disabled = true;
        killBtn.disabled = false;
        stopBtn.disabled = false;
        
        // Set credentials required to "exists" for next time
        credentialsRequired.value = 'exists';
        } else if (action === 'kill') {
        // Hide instance details
        instanceDetails.classList.add('hidden');
        emptyState.classList.remove('hidden');
        
        // Update status
        statusBadge.className = 'status-badge status-inactive';
        statusText.textContent = 'Instance Terminated';
        
        // Enable/disable buttons
        deployBtn.disabled = false;
        killBtn.disabled = true;
        stopBtn.disabled = true;
        
        // Reset credentials required
        credentialsRequired.value = 'no';
        } else if (action === 'stop') {
        // Update status
        statusBadge.className = 'status-badge status-stopped';
        statusText.textContent = 'Instance Stopped';
        
        // Enable/disable buttons
        deployBtn.disabled = false;
        killBtn.disabled = false;
        stopBtn.disabled = true;
        }
    }, 2000); // Simulate 2 second delay for action
    }
    
    // Copy functionality
    document.querySelectorAll('.copy-btn').forEach(button => {
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
});