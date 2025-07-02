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

    let domainToRemove = null;

    function isValidDomain(domain) {
        domain = domain.trim().toLowerCase();
        if (domain.length < 3 || domain.length > 63) return false;
        const regex = /^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/;
        return regex.test(domain);
    }

    function updateDomainsDisplay() {
        const cards = domainsContainer.querySelectorAll('.domain-card');
        if (cards.length === 0) {
            domainsContainer.classList.add('hidden');
            if (emptyDomains) emptyDomains.classList.remove('hidden');
        } else {
            domainsContainer.classList.remove('hidden');
            if (emptyDomains) emptyDomains.classList.add('hidden');
        }
    }

    function attachButtonListeners(scope = document) {
        scope.querySelectorAll('.remove-btn').forEach(btn => {
            btn.addEventListener('click', function () {
                domainToRemove = this.getAttribute('data-domain');
                confirmationModal.style.display = 'flex';
            });
        });

        scope.querySelectorAll('.copy-domain-btn').forEach(btn => {
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

    addBtn.addEventListener('click', function () {
        const domainName = domainInput.value.trim().toLowerCase();
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

        // Check if domain already exists
        const existing = domainsContainer.querySelector(`[data-domain="${domainName}"]`);
        if (existing) {
            errorMsg.classList.remove('hidden');
            return;
        }

        if (domainsContainer.querySelectorAll('.domain-card').length >= 3) {
            alert('You have reached the maximum of 3 domains per instance');
            return;
        }

        const formData = new URLSearchParams();
        formData.append('domain_name', domainName);
        formData.append('port', port.toString());

        fetch('/adddomain', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: formData
        })
        .then(res => {
            if (!res.ok) throw new Error('Domain creation failed');
            return res.json();
        })
        .then(data => {
            const domainCard = document.createElement('div');
            domainCard.className = 'domain-card';
            domainCard.innerHTML = `
                <div class="flex justify-between items-start">
                    <div class="w-full">
                        <div class="flex justify-between items-center">
                            <h3 class="text-lg font-bold text-white">${data.Domain}</h3>
                            <button class="remove-btn" data-domain="${data.Domain}">
                                <i class="fas fa-times"></i>
                            </button>
                        </div>
                        <div class="domain-url-container">
                            <div class="domain-url">https://${data.Domain}.kwscloud.in</div>
                            <button class="copy-domain-btn" data-url="https://${data.Domain}.kwscloud.in">
                                <i class="fas fa-copy"></i>
                            </button>
                        </div>
                        <div class="mt-4 grid grid-cols-2 gap-2">
                            <div>
                                <span class="text-gray-500">Port:</span>
                                <span class="ml-2 text-white">${data.Port}</span>
                            </div>
                            <div>
                                <span class="text-gray-500">Status:</span>
                                <span class="ml-2 text-green-400">${data.Status}</span>
                            </div>
                        </div>
                    </div>
                </div>
            `;
            domainsContainer.appendChild(domainCard);
            attachButtonListeners(domainCard);
            updateDomainsDisplay();
            domainInput.value = '';
            portInput.value = '';
        })
        .catch(err => {
            console.error(err);
            alert('Failed to add domain. It might already exist or an error occurred.');
        });
    });

    function removeDomain(domainName) {
        const formData = new URLSearchParams();
        formData.append('domain_name', domainName);

        fetch('/removedomain', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: formData
        })
        .then(res => {
            if (!res.ok) throw new Error('Failed to remove domain');

            const card = domainsContainer.querySelector(`[data-domain="${domainName}"]`)?.closest('.domain-card');
            if (card) card.remove();

            updateDomainsDisplay();
            confirmationModal.style.display = 'none';
            domainToRemove = null;
        })
        .catch(err => {
            console.error(err);
            alert('Error removing domain.');
        });
    }

    confirmRemoveBtn.addEventListener('click', function () {
        if (domainToRemove) {
            removeDomain(domainToRemove);
        }
    });

    confirmCancelBtn.addEventListener('click', function () {
        confirmationModal.style.display = 'none';
        domainToRemove = null;
    });

    attachButtonListeners();
    updateDomainsDisplay();
});
