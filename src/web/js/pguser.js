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
  
  // Copy functionality for connection info
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

    // Add event listener to remove buttons
    document.querySelectorAll('.remove-btn').forEach(button => {
    button.addEventListener('click', async function () {
        const icon = this.querySelector('i');
        const username = icon.getAttribute('data-username');
        const password = icon.getAttribute('data-password');

        if (!confirm(`Are you sure you want to delete user '${username}'?`)) return;

        const formData = new URLSearchParams();
        formData.append('user_name', username);
        formData.append('password', password);

        try {
        const response = await fetch('/deletepguser', {
            method: 'POST',
            headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
            },
            body: formData.toString()
        });

        if (!response.ok) throw new Error(`Failed to delete user. Status: ${response.status}`);
        
        // Remove the row
        const row = this.closest('tr');
        row.remove();

        // Update the count
        const userCountElement = document.querySelector('.text-white.ml-2');
        if (userCountElement) {
            const parts = userCountElement.textContent.split('/');
            const newCount = parseInt(parts[0]) - 1;
            userCountElement.textContent = `${newCount}/${parts[1]}`;
        }

        } catch (err) {
        console.error(err);
        alert('Error deleting user.');
        }
    });
    });

    // Add event listener to manage buttons
    document.querySelectorAll('.manage-btn').forEach(button => {
        button.addEventListener('click', function () {
            const pid = this.getAttribute('data-id');
            const owner = this.getAttribute("data-username");
            if (!pid) return;
            window.location.href = `/kws_services/postgres/db?pid=${encodeURIComponent(pid)}&owner=${encodeURIComponent(owner)}`;
        });
    });
  
  // Toggle password visibility in form
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
  
  // Password toggle for table rows
  document.querySelectorAll('.password-toggle').forEach(button => {
    button.addEventListener('click', function() {
      const username = this.getAttribute('data-username');
      const passwordSpan = document.getElementById(`password-${username}`);
      const hiddenPassword = document.getElementById(`real-password-${username}`);
      const icon = this.querySelector('i');
      
      if (passwordSpan.textContent === '••••••••') {
        // Show the actual password
        passwordSpan.textContent = hiddenPassword.value;
        icon.classList.replace('fa-eye-slash', 'fa-eye');
      } else {
        // Hide the password
        passwordSpan.textContent = '••••••••';
        icon.classList.replace('fa-eye', 'fa-eye-slash');
      }
    });
  });

  // Password copy for table rows
  document.querySelectorAll('.password-copy').forEach(button => {
    button.addEventListener('click', function() {
      const username = this.getAttribute('data-username');
      const hiddenPassword = document.getElementById(`real-password-${username}`);
      const icon = this.querySelector('i');

      navigator.clipboard.writeText(hiddenPassword.value);
      
      // Visual feedback
      icon.className = 'fas fa-check';
      setTimeout(() => {
        icon.className = 'fas fa-copy';
      }, 2000);
    });
  });
  
  // Add event listener to Create User button
  const createUserBtn = document.querySelector('.btn-primary');
  if (createUserBtn) {
    createUserBtn.addEventListener('click', createNewUser);
  }

  async function createNewUser() {
    const usernameInput = document.querySelector('.form-input[placeholder="Enter username"]');
    const passwordInput = document.getElementById('password');
    const confirmPasswordInput = document.querySelector('.form-input[placeholder="Confirm your password"]');
    
    // Basic validation
    if (!usernameInput.value || !passwordInput.value) {
      alert('Please fill in all fields');
      return;
    }
    
    if (passwordInput.value !== confirmPasswordInput.value) {
      alert('Passwords do not match!');
      return;
    }

    try {
      // Send POST request to create user
      const formData = new URLSearchParams();
      formData.append('user_name', usernameInput.value);
      formData.append('password', passwordInput.value);

      const response = await fetch('/createpguser', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: formData.toString()
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const newUser = await response.json();
      console.log(newUser)

      // Update the table
      addUserToTable(newUser);

      // Update user count
      const userCountElement = document.querySelector('.text-white.ml-2');
      if (userCountElement) {
        const currentCount = parseInt(userCountElement.textContent.split('/')[0]);
        userCountElement.textContent = `${currentCount + 1}/${newUser.UserLimit}`;
      }

      // Reset form
      usernameInput.value = '';
      passwordInput.value = '';
      confirmPasswordInput.value = '';
      document.getElementById('passwordStrength').style.width = '0%';

    } catch (error) {
      console.error('Error creating user:', error);
      alert('Error creating user. Please try again.');
    }
  }

  function addUserToTable(user) {
    const tbody = document.querySelector('.user-table tbody');
    
    // Create new table row
    const newRow = document.createElement('tr');
    
    // Username cell
    const usernameCell = document.createElement('td');
    usernameCell.className = 'font-mono';
    usernameCell.textContent = user.Username;
    
    // Password cell
    const passwordCell = document.createElement('td');
    passwordCell.innerHTML = `
      <div class="password-field-wrapper">
        <div class="password-display">
          <input type="hidden" id="real-password-${user.Username}" value="${user.Password}">
          <input type="hidden" id="user-id-${user.Username}" value="${user.ID}">
          <span class="password-text" id="password-${user.Username}">••••••••</span>
          <button class="password-copy" data-username="${user.Username}" title="Copy password">
            <i class="fas fa-copy text-sm"></i>
          </button>
          <button class="password-toggle" data-username="${user.Username}" title="Show password">
            <i class="fas fa-eye-slash text-sm"></i>
          </button>
        </div>
      </div>
    `;
    
    // Permissions cell
    const permissionsCell = document.createElement('td');
    permissionsCell.innerHTML = `
      <span class="status-badge status-active">${user.Permissions}</span>
    `;
    
    // Actions cell
    const actionsCell = document.createElement('td');
    actionsCell.innerHTML = `
      <button class="action-btn remove-btn" data-username="${user.Username}" data-password="${user.Password}">
        <i class="fas fa-trash mr-1"></i> Remove
      </button>
      <button class="action-btn manage-btn" data-username="${user.Username}" data-password="${user.Password}" data-id="${user.ID}">
        <i class="fas fa-database mr-1"></i> Manage
      </button>
    `;
    
    // Assemble row
    newRow.appendChild(usernameCell);
    newRow.appendChild(passwordCell);
    newRow.appendChild(permissionsCell);
    newRow.appendChild(actionsCell);
    
    // Add to table
    tbody.appendChild(newRow);

    // Add remove event to the dynamically created button
    actionsCell.querySelector('.remove-btn')?.addEventListener('click', async function () {
    const username = this.getAttribute('data-username');
    const password = this.getAttribute('data-password');

    if (!confirm(`Are you sure you want to delete user '${username}'?`)) return;

    const formData = new URLSearchParams();
    formData.append('user_name', username);
    formData.append('password', password);

    try {
        const response = await fetch('/deletepguser', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: formData.toString()
        });

        if (!response.ok) throw new Error(`Failed to delete user. Status: ${response.status}`);
        
        // Remove the row
        const row = this.closest('tr');
        row.remove();

        // Update the count
        const userCountElement = document.querySelector('.text-white.ml-2');
        if (userCountElement) {
        const parts = userCountElement.textContent.split('/');
        const newCount = parseInt(parts[0]) - 1;
        userCountElement.textContent = `${newCount}/${parts[1]}`;
        }

    } catch (err) {
        console.error(err);
        alert('Error deleting user.');
    }
    });

    // Add manage event to the dynamically created button
    actionsCell.querySelector('.manage-btn')?.addEventListener('click', function () {
        const pid = this.getAttribute('data-id');
        const owner = this.getAttribute("data-username");
        if (!pid) return;
        window.location.href = `/kws_services/postgres/db?pid=${encodeURIComponent(pid)}&owner=${encodeURIComponent(owner)}`;
    });

    
    // Attach event listeners to new password controls
    attachPasswordListeners(newRow);
  }

  function attachPasswordListeners(row) {
    // Toggle password visibility
    row.querySelector('.password-toggle')?.addEventListener('click', function() {
      const username = this.getAttribute('data-username');
      const passwordSpan = document.getElementById(`password-${username}`);
      const hiddenPassword = document.getElementById(`real-password-${username}`);
      const icon = this.querySelector('i');
      
      if (passwordSpan.textContent === '••••••••') {
        passwordSpan.textContent = hiddenPassword.value;
        icon.classList.replace('fa-eye-slash', 'fa-eye');
      } else {
        passwordSpan.textContent = '••••••••';
        icon.classList.replace('fa-eye', 'fa-eye-slash');
      }
    });

    // Copy password
    row.querySelector('.password-copy')?.addEventListener('click', function() {
      const username = this.getAttribute('data-username');
      const hiddenPassword = document.getElementById(`real-password-${username}`);
      const icon = this.querySelector('i');

      navigator.clipboard.writeText(hiddenPassword.value);
      
      // Visual feedback
      icon.className = 'fas fa-check';
      setTimeout(() => {
        icon.className = 'fas fa-copy';
      }, 2000);
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
