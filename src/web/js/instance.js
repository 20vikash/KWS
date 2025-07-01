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

  const deployBtn = document.getElementById('deploy-btn');
  const killBtn = document.getElementById('kill-btn');
  const stopBtn = document.getElementById('stop-btn');
  const codeBtn = document.getElementById('code-btn');
  const deployModal = document.getElementById('deploy-modal');
  const closeModal = document.getElementById('close-modal');
  const cancelDeploy = document.getElementById('cancel-deploy');
  const confirmDeploy = document.getElementById('confirm-deploy');
  const credentialsRequired = document.getElementById('credentials-required');

  updateCodeButtonState();

  deployBtn.addEventListener('click', function () {
    if (credentialsRequired.value === 'no') {
      deployModal.classList.remove('hidden');
    } else {
      startDeployDirectly();
    }
  });

  closeModal.addEventListener('click', () => deployModal.classList.add('hidden'));
  cancelDeploy.addEventListener('click', () => deployModal.classList.add('hidden'));

  confirmDeploy.addEventListener('click', function () {
    deployModal.classList.add('hidden');
    const username = document.getElementById('deploy-username').value.trim();
    const password = document.getElementById('deploy-password').value.trim();
    const confirm = document.getElementById('deploy-confirm').value.trim();

    if (password !== confirm) {
      alert('Passwords do not match!');
      return;
    }

    lockActionButtons();

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
        unlockActionButtons();
      });
  });

  killBtn.addEventListener('click', function () {
    startAction('kill');
  });

  stopBtn.addEventListener('click', function () {
    startAction('stop');
  });

  codeBtn.addEventListener('click', function () {
    alert('Opening VS Code in browser...');
  });

  function startDeployDirectly() {
    lockActionButtons();

    fetch("/deploy", {
      method: "POST",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      body: new URLSearchParams(),
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
        unlockActionButtons();
      });
  }

  function startAction(type) {
    lockActionButtons();

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
        unlockActionButtons();
      });
  }

  function pollDeployResult(jobID, attempts = 0) {
    if (attempts > 20) {
      alert("Deployment timed out.");
      unlockActionButtons();
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

        updateInstanceDetails(data.Instance);
        unlockActionButtons();
      })
      .catch(err => {
        console.error("Polling failed:", err);
        unlockActionButtons();
      });
  }

  function pollSKResult(type, jobID, attempts = 0) {
    if (attempts > 20) {
      alert(`${type} action timed out.`);
      unlockActionButtons();
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
          unlockActionButtons();
          return;
        }

        if (type === 'stop') {
          document.getElementById('status-badge').className = 'status-badge status-stopped';
          document.getElementById('status-text').textContent = 'Instance Stopped';
          document.getElementById("credentials-required").value = "exists";

          deployBtn.disabled = false;
          killBtn.disabled = false;
          stopBtn.disabled = true;
        }

        if (type === 'kill') {
          document.getElementById('instance-details').classList.add('hidden');
          document.getElementById('empty-state').classList.remove('hidden');

          document.getElementById('status-badge').className = 'status-badge status-inactive';
          document.getElementById('status-text').textContent = 'Instance not deployed';
          document.getElementById("credentials-required").value = "no";

          deployBtn.disabled = false;
          killBtn.disabled = true;
          stopBtn.disabled = true;
        }

        updateCodeButtonState();
        removeBlinking();
      })
      .catch(err => {
        console.error(`${type} polling failed:`, err);
        unlockActionButtons();
      });
  }

  function updateInstanceDetails(instance) {
    const instanceDetails = document.getElementById('instance-details');
    const emptyState = document.getElementById('empty-state');
    const statusBadge = document.getElementById('status-badge');
    const statusText = document.getElementById('status-text');

    instanceDetails.classList.remove('hidden');
    emptyState.classList.add('hidden');

    document.getElementById('instance-username').textContent = instance.Username;
    document.getElementById('instance-password').textContent = '••••••••';
    document.getElementById('instance-ip').textContent = instance.IP;
    document.getElementById('instance-ssh').textContent = `${instance.Username}@${instance.IP}`;

    statusBadge.className = 'status-badge status-active';
    statusText.textContent = 'Instance Active';

    document.querySelectorAll('.copy-btn').forEach(btn => {
      const copyType = btn.previousElementSibling.id;
      switch (copyType) {
        case "instance-username":
          btn.dataset.copy = instance.Username;
          break;
        case "instance-password":
          btn.dataset.copy = instance.Password;
          break;
        case "instance-ip":
          btn.dataset.copy = instance.IP;
          break;
        case "instance-ssh":
          btn.dataset.copy = `ssh ${instance.Username}@${instance.IP}`;
          break;
      }
    });

    document.getElementById("credentials-required").value = "exists";
    updateCodeButtonState();
  }

  function updateCodeButtonState() {
    const status = document.getElementById('status-badge').className;
    codeBtn.disabled = !status.includes('status-active');
  }

  function lockActionButtons() {
    deployBtn.disabled = true;
    killBtn.disabled = true;
    stopBtn.disabled = true;

    deployBtn.classList.add("action-blinking");
    killBtn.classList.add("action-blinking");
    stopBtn.classList.add("action-blinking");
  }

  function unlockActionButtons() {
    deployBtn.disabled = false;
    killBtn.disabled = false;
    stopBtn.disabled = false;
    removeBlinking();
  }

  function removeBlinking() {
    deployBtn.classList.remove("action-blinking");
    killBtn.classList.remove("action-blinking");
    stopBtn.classList.remove("action-blinking");
  }

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
});
