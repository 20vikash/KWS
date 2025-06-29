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
    
    // Modal handling
    const deleteModal = document.getElementById('deleteModal');
    const dbNameToDelete = document.getElementById('dbNameToDelete');
    const closeModal = document.getElementById('closeModal');
    const cancelDelete = document.getElementById('cancelDelete');
    const confirmDelete = document.getElementById('confirmDelete');
    const passwordInput = document.getElementById('passwordInput');
    const passwordError = document.getElementById('passwordError');
    
    // State variables
    let currentDbToDelete = '';
    
    // Add event listeners to all remove buttons
    document.querySelectorAll('.remove-btn').forEach(button => {
        button.addEventListener('click', function() {
        // Get database name from table row
        const dbRow = this.closest('tr');
        const dbName = dbRow.querySelector('td:first-child').textContent;
        
        // Set modal content
        currentDbToDelete = dbName;
        dbNameToDelete.textContent = dbName;
        
        // Reset password field and error
        passwordInput.value = '';
        passwordError.classList.add('hidden');
        
        // Show modal
        deleteModal.classList.remove('hidden');
        });
    });
    
    // Close modal handlers
    document.getElementById('closeModal').addEventListener('click', closeModal);
    document.getElementById('cancelDelete').addEventListener('click', closeModal);
    
    // Confirm delete handler
    document.getElementById('confirmDelete').addEventListener('click', function() {
        const password = passwordInput.value.trim();
        
        if (!password) {
        showPasswordError('Please enter your password');
        return;
        }
        
        // Here you would normally make a request to verify password
        // For DOM-only implementation, we'll just show an example
        showPasswordError('Password verification would happen here');
        
        // In a real implementation, you would:
        // 1. Send password to server for verification
        // 2. If valid: proceed with deletion
        // 3. If invalid: show error message
    });
    
    function closeModal() {
        deleteModal.classList.add('hidden');
        passwordInput.value = '';
        passwordError.classList.add('hidden');
        currentDbToDelete = '';
    }
    
    function showPasswordError(message) {
        passwordError.textContent = message;
        passwordError.classList.remove('hidden');
    }
    
    function showDeleteModal(dbName) {
    dbNameToDelete.textContent = dbName;
    deleteModal.classList.remove('hidden');
    }
    
    function closeDeleteModal() {
    deleteModal.classList.add('hidden');
    }
    
    if (closeModal) closeModal.addEventListener('click', closeDeleteModal);
    if (cancelDelete) cancelDelete.addEventListener('click', closeDeleteModal);
    
    // Create database button
    const createDbBtn = document.getElementById('createDbBtn');
    if (createDbBtn) {
    createDbBtn.addEventListener('click', function() {
        const dbName = document.getElementById('dbName').value;
        const dbOwner = document.getElementById('dbOwner').value;
        
        if (!dbName) {
        alert('Please enter a database name');
        return;
        }
        
        if (!dbOwner) {
        alert('Please select an owner for the database');
        return;
        }
        
        alert(`Creating database: ${dbName}\nOwner: ${dbOwner}`);
        // In a real application, you would make an API call here
    });
    }
    
    // Delete confirmation
    if (confirmDelete) {
    confirmDelete.addEventListener('click', function() {
        const dbName = dbNameToDelete.textContent;
        alert(`Database "${dbName}" deleted successfully!`);
        closeDeleteModal();
    });
    }
});
