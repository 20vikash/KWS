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

    // DOM elements
    const domainInput = document.getElementById('domain-name');
    const portInput = document.getElementById('domain-port');
    const addBtn = document.getElementById('add-domain');
    const errorMsg = document.getElementById('domain-error');
    const domainsContainer = document.getElementById('domains-container');
    const emptyDomains = document.getElementById('empty-domains');
    const confirmationModal = document.getElementById('confirmation-modal');
    const confirmRemoveBtn = document.getElementById('confirm-remove');
    const confirmCancelBtn = document.getElementById('confirm-cancel');

    // Empty domains array (starts clean)
    let domains = [];

    // Domain to be removed
    let domainToRemove = null;

    // Update domains display
    function updateDomainsDisplay() {
        if (domains.length === 0) {
            domainsContainer.classList.add('hidden');
            emptyDomains.classList.remove('hidden');
            domainsContainer.innerHTML = '';
        } else {
            domainsContainer.classList.remove('hidden');
            emptyDomains.classList.add('hidden');
            renderDomains();
        }
    }

    // Render domains
    function renderDomains() {
        domainsContainer.innerHTML = '';

        domains.forEach(domain => {
            const domainCard = document.createElement('div');
            domainCard.className = 'domain-card';
            domainCard.innerHTML = `
                <div class="flex justify-between items-start">
                    <div class="w-full">
                        <div class="flex justify-between items-center">
                            <h3 class="text-lg font-bold text-white">${domain.name}</h3>
                            <button class="remove-btn" data-domain="${domain.name}">
                                <i class="fas fa-times"></i>
                            </button>
                        </div>
                        <div class="domain-url-container">
                            <div class="domain-url">https://${domain.name}.kwscloud.in</div>
                            <button class="copy-domain-btn" data-url="https://${domain.name}.kwscloud.in">
                                <i class="fas fa-copy"></i>
                            </button>
                        </div>
                        <div class="mt-4 grid grid-cols-2 gap-2">
                            <div>
                                <span class="text-gray-500">Port:</span>
                                <span class="ml-2 text-white">${domain.port}</span>
                            </div>
                            <div>
                                <span class="text-gray-500">Status:</span>
                                <span class="ml-2 text-green-400">Active</span>
                            </div>
                        </div>
                    </div>
                </div>
            `;
            domainsContainer.appendChild(domainCard);
        });

        // Add event listeners to remove buttons
        document.querySelectorAll('.remove-btn').forEach(btn => {
            btn.addEventListener('click', function () {
                domainToRemove = this.getAttribute('data-domain');
                confirmationModal.style.display = 'flex';
            });
        });

        // Add event listeners to copy buttons
        document.querySelectorAll('.copy-domain-btn').forEach(btn => {
            btn.addEventListener('click', function () {
                const url = this.getAttribute('data-url');
                navigator.clipboard.writeText(url);
                const icon = this.querySelector('i');
                icon.className = 'fas fa-check';
                setTimeout(() => {
                    icon.className = 'fas fa-copy';
                }, 2000);
            });
        });
    }

    // Remove domain
    function removeDomain(domainName) {
        domains = domains.filter(d => d.name !== domainName);
        updateDomainsDisplay();
        confirmationModal.style.display = 'none';
    }

    // Validate domain name
    function isValidDomain(domain) {
        domain = domain.trim().toLowerCase();
        if (domain.length < 3 || domain.length > 63) return false;
        const regex = /^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/;
        return regex.test(domain);
    }

    // Add new domain
    addBtn.addEventListener('click', function () {
        const domainName = domainInput.value.trim();
        const port = parseInt(portInput.value);

        errorMsg.classList.add('hidden');

        if (!domainName || isNaN(port)) {
            alert('Please fill in both fields');
            return;
        }

        if (port < 1 || port > 65535) {
            alert('Please enter a valid port number (1-65535)');
            return;
        }

        if (!isValidDomain(domainName)) {
            alert('Domain name must be 3-63 characters, start and end with letter/number, and contain only lowercase letters, numbers, and hyphens');
            return;
        }

        const normalizedDomain = domainName.toLowerCase();

        if (domains.some(d => d.name === normalizedDomain)) {
            errorMsg.classList.remove('hidden');
            return;
        }

        if (domains.length >= 3) {
            alert('You have reached the maximum of 3 domains per instance');
            return;
        }

        domains.push({ name: normalizedDomain, port: port });
        updateDomainsDisplay();

        domainInput.value = '';
        portInput.value = '';
    });

    // Confirmation modal buttons
    confirmRemoveBtn.addEventListener('click', function () {
        if (domainToRemove) {
            removeDomain(domainToRemove);
        }
    });

    confirmCancelBtn.addEventListener('click', function () {
        confirmationModal.style.display = 'none';
        domainToRemove = null;
    });

    // Initial render
    updateDomainsDisplay();
});
