document.addEventListener('DOMContentLoaded', function () {
  // Sparkle effect
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

  // State management
  const stateInput = document.getElementById('instance-state');
  const deployBtn = document.getElementById('deploy-btn');
  const killBtn = document.getElementById('kill-btn');
  const stopBtn = document.getElementById('stop-btn');
  const codeBtn = document.getElementById('code-btn');
  const deployModal = document.getElementById('deploy-modal');
  const closeModal = document.getElementById('close-modal');
  const cancelDeploy = document.getElementById('cancel-deploy');
  const confirmDeploy = document.getElementById('confirm-deploy');
  const statusBadge = document.getElementById('status-badge');
  const statusText = document.getElementById('status-text');
  const instanceDetails = document.getElementById('instance-details');
  const emptyState = document.getElementById('empty-state');
  const hiddenUsername = document.getElementById('hidden-username');
  const hiddenPassword = document.getElementById('hidden-password');
  const containerName = document.getElementById('container-name').value;

  // Initialize button states
  function updateButtonStates() {
    const state = stateInput.value;
    
    deployBtn.disabled = (state === 'active');
    killBtn.disabled = (state !== 'active' && state !== 'stopped');
    stopBtn.disabled = (state !== 'active');
    codeBtn.disabled = (state !== 'active');
  }

  // Update UI based on state
  function updateUIFromState() {
    const state = stateInput.value;
    
    // Update status badge
    statusBadge.className = 'status-badge';
    if (state === 'active') {
      statusBadge.classList.add('status-active');
      statusText.textContent = 'Instance Active';
      instanceDetails.classList.remove('hidden');
      emptyState.classList.add('hidden');
    } else if (state === 'stopped') {
      statusBadge.classList.add('status-stopped');
      statusText.textContent = 'Instance Stopped';
      instanceDetails.classList.remove('hidden');
      emptyState.classList.add('hidden');
    } else {
      statusBadge.classList.add('status-inactive');
      statusText.textContent = 'Instance not deployed';
      instanceDetails.classList.add('hidden');
      emptyState.classList.remove('hidden');
    }
    
    updateButtonStates();
  }
  
  // Update visible instance details
  function updateInstanceDetails(instance) {
    document.getElementById('instance-username').textContent = instance.Username;
    document.getElementById('instance-password').textContent = '••••••••';
    document.getElementById('instance-ip').textContent = instance.IP;
    document.getElementById('instance-ssh').textContent = `${instance.Username}@${instance.IP}`;
    
    // Update copy buttons
    document.querySelector('[data-copy]').setAttribute('data-copy', instance.Username);
    document.querySelectorAll('.copy-btn')[1].setAttribute('data-copy', instance.Password);
    document.querySelectorAll('.copy-btn')[2].setAttribute('data-copy', instance.IP);
    document.querySelectorAll('.copy-btn')[3].setAttribute('data-copy', `ssh ${instance.Username}@${instance.IP}`);
  }
  
  // Set initial state
  updateUIFromState();

  // Disable all buttons and add blinking to the active button
  function lockButton(activeButton) {
    // Disable all buttons
    deployBtn.disabled = true;
    killBtn.disabled = true;
    stopBtn.disabled = true;
    codeBtn.disabled = true;
    
    // Remove blinking from all buttons
    deployBtn.classList.remove("action-blinking");
    killBtn.classList.remove("action-blinking");
    stopBtn.classList.remove("action-blinking");
    codeBtn.classList.remove("action-blinking");
    
    // Add blinking to the active button (and keep it disabled)
    activeButton.classList.add("action-blinking");
  }
  
  // Re-enable all buttons and remove blinking
  function unlockAllButtons() {
    // Remove blinking from all buttons
    deployBtn.classList.remove("action-blinking");
    killBtn.classList.remove("action-blinking");
    stopBtn.classList.remove("action-blinking");
    codeBtn.classList.remove("action-blinking");
    
    // Re-enable buttons based on current state
    updateButtonStates();
  }

  // Button event handlers
  deployBtn.addEventListener('click', function () {
    lockButton(this);
    
    if (stateInput.value === 'inactive') {
      deployModal.classList.remove('hidden');
    } else if (stateInput.value === 'stopped') {
      // Reuse existing credentials
      const username = hiddenUsername.value;
      const password = hiddenPassword.value;
      startDeployWithCredentials(username, password);
    }
  });

  closeModal.addEventListener('click', () => {
    deployModal.classList.add('hidden');
    unlockAllButtons();
  });
  
  cancelDeploy.addEventListener('click', () => {
    deployModal.classList.add('hidden');
    unlockAllButtons();
  });

  confirmDeploy.addEventListener('click', function () {
    const username = document.getElementById('deploy-username').value.trim();
    const password = document.getElementById('deploy-password').value.trim();
    const confirm = document.getElementById('deploy-confirm').value.trim();

    if (!username || !password) {
      alert('Username and password are required!');
      unlockAllButtons();
      return;
    }
    
    if (password !== confirm) {
      alert('Passwords do not match!');
      unlockAllButtons();
      return;
    }

    deployModal.classList.add('hidden');
    startDeployWithCredentials(username, password);
  });

  killBtn.addEventListener('click', function () {
    lockButton(this);
    startAction('kill');
  });

  stopBtn.addEventListener('click', function () {
    lockButton(this);
    startAction('stop');
  });

  codeBtn.addEventListener('click', function () {
    const containerName = document.getElementById('container-name').value;
    // Open VS Code in a new tab
    if (containerName) {
      window.open(`http://${containerName}.kwscloud.in`, '_blank');
    } else {
      alert('Container name not found!');
    }
  });

  const publishBtn = document.getElementById('publish-btn');
  if (publishBtn) {
    publishBtn.addEventListener('click', function () {
      window.location.href = '/kws_publish';
    });
  }

  // Copy buttons
  document.querySelectorAll('.copy-btn').forEach(button => {
    button.addEventListener('click', function () {
      const text = this.getAttribute('data-copy');
      navigator.clipboard.writeText(text);
      const icon = this.querySelector('i');
      icon.className = 'fas fa-check';
      this.classList.add('copied');
      setTimeout(() => {
        icon.className = 'fas fa-copy';
        this.classList.remove('copied');
      }, 2000);
    });
  });

  function startDeployWithCredentials(username, password) {
    const formData = new URLSearchParams();
    formData.append("insUser", username);
    formData.append("insPassword", password);

    fetch("/deploy", {
      method: "POST",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      body: formData,
    })
      .then(res => {
        if (!res.ok) throw new Error("Deployment failed");
        return res.json();
      })
      .then(data => {
        if (!data.jobID) throw new Error("No jobID received");
        pollDeployResult(data.jobID);
      })
      .catch(err => {
        alert(err.message);
        unlockAllButtons();
      });
  }

  function startAction(type) {
    fetch(`/${type}`, {
      method: "GET"
    })
      .then(res => {
        if (!res.ok) throw new Error(`${type} failed`);
        return res.json();
      })
      .then(data => {
        if (!data.jobID) throw new Error("No jobID received");
        pollSKResult(type, data.jobID);
      })
      .catch(err => {
        alert(err.message);
        unlockAllButtons();
      });
  }

  function pollDeployResult(jobID, attempts = 0) {
    if (attempts > 20) {
      alert("Deployment timed out.");
      unlockAllButtons();
      return;
    }

    fetch(`/deployresult?jobID=${encodeURIComponent(jobID)}`, {
      method: "POST"
    })
      .then(res => res.json())
      .then(data => {
        if (!data.Done) {
          setTimeout(() => pollDeployResult(jobID, attempts + 1), 2000);
          return;
        }

        // Update credentials in DOM
        document.getElementById('hidden-username').value = data.Instance.Username;
        document.getElementById('hidden-password').value = data.Instance.Password;
        document.getElementById("container-name").value = data.Instance.ContainerID;
        
        // Update visible instance details
        updateInstanceDetails(data.Instance);
        
        stateInput.value = 'active';
        updateUIFromState();
        // Show the publish section again
        const publishSection = document.getElementById('publish-section');
        if (publishSection) {
          publishSection.classList.remove('hidden');
        }
        unlockAllButtons();
      })
      .catch(err => {
        console.error("Polling failed:", err);
        unlockAllButtons();
      });
  }

  function pollSKResult(type, jobID, attempts = 0) {
    if (attempts > 20) {
      alert(`${type} action timed out.`);
      unlockAllButtons();
      return;
    }

    const endpoint = type === 'kill' ? '/killresut' : '/stopresult';

    fetch(`${endpoint}?jobID=${encodeURIComponent(jobID)}`, {
      method: "POST"
    })
      .then(res => res.json())
      .then(data => {
        if (!data.Done) {
          setTimeout(() => pollSKResult(type, jobID, attempts + 1), 2000);
          return;
        }

        if (!data.Success) {
          alert(`${type.charAt(0).toUpperCase() + type.slice(1)} failed.`);
          unlockAllButtons();
          return;
        }

        if (type === 'kill') {
          stateInput.value = 'inactive';
          updateUIFromState();

          // Remove/hide the publish section if it exists
          const publishSection = document.getElementById('publish-section');
          if (publishSection) {
            publishSection.classList.add('hidden');
          }

          unlockAllButtons();
          return;
        }

        stateInput.value = type === 'kill' ? 'inactive' : 'stopped';
        updateUIFromState();
        unlockAllButtons();
      })
      .catch(err => {
        console.error(`${type} polling failed:`, err);
        unlockAllButtons();
      });
  }
});