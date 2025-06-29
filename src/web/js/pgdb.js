window.showDeleteModal = function(dbName) {
    console.log("dbname", dbName)
  const dbNameToDelete = document.getElementById('dbNameToDelete');
  const deleteModal = document.getElementById('deleteModal');
  if (dbNameToDelete && deleteModal) {
    dbNameToDelete.textContent = dbName;
    window.pendingDeleteName = dbName;
    deleteModal.classList.remove('hidden');
  }
};


document.addEventListener('DOMContentLoaded', function () {
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

  let pendingDeleteName = '';

  function closeDeleteModal() {
    deleteModal.classList.add('hidden');
    pendingDeleteName = '';
  }

  if (closeModal) closeModal.addEventListener('click', closeDeleteModal);
  if (cancelDelete) cancelDelete.addEventListener('click', closeDeleteModal);

  const updateStats = (total, limit) => {
    document.querySelectorAll('.text-white').forEach(el => {
      if (el.innerText === `${total - 1}` || el.innerText === `${total + 1}` || el.innerText === `${total}`) {
        el.innerText = `${total}`;
      }
    });

    const available = limit - total;
    const percent = Math.round((total / limit) * 100);

    document.querySelectorAll('.text-white').forEach(el => {
      if (el.innerText.includes('/')) {
        el.innerText = `${total}/${limit}`;
      }
    });

    document.querySelectorAll('.text-white').forEach(el => {
      if (el.innerText.includes('%')) {
        el.innerText = `${percent}% used`;
      }
    });

    const progress = document.querySelector('.progress-fill');
    if (progress) progress.style.width = `${percent}%`;
  };

  const createDbBtn = document.getElementById('createDbBtn');
  if (createDbBtn) {
    createDbBtn.addEventListener('click', async function () {
      const dbName = document.getElementById('dbName').value.trim();
      const dbOwner = document.getElementById('dbOwner').value;
      const dbEncoding = document.getElementById('dbEncoding').value;
      const dbPassword = document.getElementById('pg-user-password').value;
      const urlParams = new URLSearchParams(window.location.search);
      const ownerFromURL = urlParams.get("owner");

      if (!dbName) return alert("Please enter a database name");
      if (!ownerFromURL || !dbPassword) return alert("Missing owner or password info");

      try {
        const payload = new URLSearchParams({
          db_name: dbName,
          user_name: ownerFromURL,
          password: dbPassword,
          encoding: dbEncoding
        });

        const response = await fetch('/createpgdb', {
          method: 'POST',
          headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
          body: payload
        });

        if (response.status === 201) {
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

          // Update stats
          const totalNow = tbody.children.length;
          const limit = parseInt(document.querySelector('.db-limit').innerText.match(/\d+/)[0]);
          updateStats(totalNow, limit);

          document.getElementById('dbName').value = '';
        } else {
          alert(`Error: ${result.error || "Unknown error occurred"}`);
        }
      } catch (err) {
        alert("Failed to create database: " + err.message);
      }
    });
  }

  if (confirmDelete) {
    confirmDelete.addEventListener('click', async function () {
      const dbName = window.pendingDeleteName;
      const urlParams = new URLSearchParams(window.location.search);
      const userName = urlParams.get("owner");
      const password = document.getElementById('pg-user-password').value;

      if (!dbName || !userName || !password) return alert("Missing data for delete");

      try {
        const payload = new URLSearchParams({
          db_name: dbName,
          user_name: userName,
          password: password
        });

        const response = await fetch('/deletepgdb', {
          method: 'POST',
          headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
          body: payload
        });

        if (response.ok) {
          const tbody = document.querySelector('.database-table tbody');
          const rows = Array.from(tbody.querySelectorAll('tr'));
          for (const row of rows) {
            if (row.innerText.includes(dbName)) {
              row.remove();
              break;
            }
          }

          // Update stats
          const totalNow = tbody.children.length;
          const limit = parseInt(document.querySelector('.db-limit').innerText.match(/\d+/)[0]);
          updateStats(totalNow, limit);

          closeDeleteModal();
        } else {
          alert("Failed to delete database: " + (result.error || "Unknown error"));
        }
      } catch (err) {
        alert("Error deleting database: " + err.message);
      }
    });
  }
});
