document.addEventListener("DOMContentLoaded", () => {
    const MAX_DEVICES = 3;
    const deviceList = document.querySelector(".card-grid");
    const addForm = document.getElementById("add-device-form");
    const warning = document.getElementById("limit-warning");
    const registerForm = document.getElementById("register-form");
    
    // Changed to event delegation for remove forms
    deviceList.addEventListener("submit", handleRemoveSubmit);

    function updateUI() {
        const deviceCount = deviceList.querySelectorAll(".device-card").length;
        if (deviceCount >= MAX_DEVICES) {
            addForm.classList.add("hidden");
            warning.classList.remove("hidden");
        } else {
            addForm.classList.remove("hidden");
            warning.classList.add("hidden");
        }
    }

    updateUI();

    registerForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const keyInput = document.getElementById("user_public_key");
        const publicKey = keyInput.value.trim();
        if (!publicKey) return;

        const res = await fetch("/register", {
            method: "POST",
            headers: { "Content-Type": "application/x-www-form-urlencoded" },
            body: new URLSearchParams({ public_key: publicKey })
        });

        if (res.ok) {
            const newDevice = await res.json();

            // Updated to match HTML template structure
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
                        <form class="remove-form">
                            <input type="hidden" name="public_key" value="${newDevice.publicKey}" />
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
        if (e.target.classList.contains('remove-form')) {
            e.preventDefault();
            const form = e.target;
            const formData = new FormData(form);
            
            const res = await fetch("/remove", {
                method: "POST",
                body: formData,
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
