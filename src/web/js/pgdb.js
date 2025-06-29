document.addEventListener('DOMContentLoaded', function () {
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

  // 🌟 CREATE DATABASE HANDLER
  const createDbBtn = document.getElementById('createDbBtn');
  if (createDbBtn) {
    createDbBtn.addEventListener('click', async function () {
      const dbName = document.getElementById('dbName').value.trim();
      const dbOwner = document.getElementById('dbOwner').value;
      const dbEncoding = document.getElementById('dbEncoding').value;
      const dbPassword = document.getElementById('pg-user-password').value;

      const urlParams = new URLSearchParams(window.location.search);
      const ownerFromURL = urlParams.get("owner");

      if (!dbName) {
        alert("Please enter a database name");
        return;
      }

      if (!ownerFromURL || !dbPassword) {
        alert("Missing owner or password info");
        return;
      }

      try {
        const payload = new URLSearchParams({
          db_name: dbName,
          user_name: ownerFromURL,
          password: dbPassword,
          encoding: dbEncoding
        });

        const response = await fetch('/createpgdb', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
          },
          body: payload
        });

        const result = await response.json();
        if (result.success) {
          const tbody = document.querySelector('.database-table tbody');
          const newRow = document.createElement('tr');
          newRow.innerHTML = `
            <td class="font-mono">${dbName}</td>
            <td><span class="owner-badge">${ownerFromURL}</span></td>
            <td>
              <button class="action-btn remove-btn" onclick="showDeleteModal('${dbName}')">
                <i class="fas fa-trash mr-1"></i> Remove
              </button>
            </td>`;
          tbody.appendChild(newRow);
        } else {
          alert(`Error: ${result.error || "Unknown error occurred"}`);
        }
      } catch (err) {
        alert("Failed to create database: " + err.message);
      }
    });
  }

  // DELETE Confirm
  if (confirmDelete) {
    confirmDelete.addEventListener('click', function () {
      const dbName = dbNameToDelete.textContent;
      alert(`Database "${dbName}" deleted successfully!`);
      closeDeleteModal();
    });
  }
});
