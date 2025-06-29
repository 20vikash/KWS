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