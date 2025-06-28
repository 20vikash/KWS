document.addEventListener("DOMContentLoaded", () => {
    const MAX_DEVICES = 3;
    const deviceList = document.querySelector(".card-grid");
    const addForm = document.querySelector(".add-device-form"); // Changed to class selector
    const registerForm = document.getElementById("register-form");
    
    // Event delegation for remove forms
    if (deviceList) {
        deviceList.addEventListener("submit", handleRemoveSubmit);
    }

    function updateUI() {
        const deviceCount = deviceList.querySelectorAll(".device-card").length;
        
        // Find or create the warning element
        let warning = document.querySelector(".device-limit-warning");
        if (!warning) {
            warning = document.createElement("div");
            warning.className = "bg-yellow-900/50 text-yellow-200 p-4 rounded-lg mb-6 border border-yellow-700/50 flex items-start device-limit-warning";
            warning.innerHTML = `
                <i class="fas fa-exclamation-circle mt-1 mr-3"></i>
                <div>
                    <p class="font-medium">Device Limit Reached</p>
                    <p class="text-sm mt-1">You've reached the maximum device limit (${MAX_DEVICES}). Remove a device to add a new one.</p>
                </div>
            `;
        }

        const mainContent = document.querySelector('main > .max-w-6xl');
        
        if (deviceCount >= MAX_DEVICES) {
            // Hide add form and show warning
            if (addForm) addForm.classList.add("hidden");
            
            // Insert warning if not already present
            if (!document.querySelector(".device-limit-warning")) {
                const header = document.querySelector('.flex.items-center.mb-8');
                header.parentNode.insertBefore(warning, header.nextElementSibling);
            }
        } else {
            // Show add form and hide warning
            if (addForm) addForm.classList.remove("hidden");
            
            // Remove warning if present
            const existingWarning = document.querySelector(".device-limit-warning");
            if (existingWarning) existingWarning.remove();
        }
    }

    updateUI();

    registerForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const keyInput = document.getElementById("user_public_key"); // Corrected ID
        const publicKey = keyInput.value.trim();
        if (!publicKey) return;

        const res = await fetch("/register", {
            method: "POST",
            headers: { "Content-Type": "application/x-www-form-urlencoded" },
            body: new URLSearchParams({ public_key: publicKey })
        });

        if (res.ok) {
            const newDevice = await res.json();

            const deviceCard = document.createElement("div");
            deviceCard.className = "device-card p-5 glow-effect";
            deviceCard.innerHTML = `
                <div class="flex justify-between items-start">
                    <div>
                        <div class="device-icon mb-3">
                            ${newDevice.active 
                                ? '<i class="fas fa-laptop-code text-blue-400"></i>' 
                                : '<i class="fas fa-laptop text-gray-500"></i>'}
                        </div>
                        <h3 class="text-lg font-bold text-white">Device</h3>
                        <p class="text-gray-400 text-sm mt-1">Allocated IP: ${newDevice.ip}</p>
                    </div>
                    <span class="device-status ${newDevice.active ? 'status-active' : 'status-inactive'}">
                        <span class="pulse ${newDevice.active ? 'pulse-active' : 'pulse-inactive'}"></span>
                        ${newDevice.active ? 'Active' : 'Disconnected'}
                    </span>
                </div>
                <div class="mt-5">
                    <div>
                        <p class="text-gray-500 text-sm mb-1">Public Key</p>
                        <div class="key-display text-gray-300">
                            ${newDevice.publicKey}
                        </div>
                    </div>
                    <div class="mt-4">
                        <form class="remove-device">
                            <input id="${newDevice.publicKey}" type="hidden" name="public_key" value="${newDevice.publicKey}" />
                            <button class="w-full bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded-md transition flex items-center justify-center">
                                <i class="fas fa-trash mr-2"></i>Remove Device
                            </button>
                        </form>
                    </div>
                </div>
            `;
            deviceList.appendChild(deviceCard);
            keyInput.value = "";
            updateUI();
        } else {
            alert("Failed to add device");
        }
    });

    async function handleRemoveSubmit(e) {
        const form = e.target.closest('form.remove-device');
        if (form) {
            e.preventDefault();

            const publicKey = form.id.trim(); // get public key from form id

            const res = await fetch("/remove", {
                method: "POST",
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
                body: new URLSearchParams({ public_key: publicKey })
            });

            if (res.ok) {
                form.closest(".device-card").remove();
                updateUI();
            } else {
                alert("Failed to remove device");
            }
        }
    }
});
