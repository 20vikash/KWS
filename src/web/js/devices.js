document.addEventListener("DOMContentLoaded", () => {
const MAX_DEVICES = 3;
const deviceList = document.getElementById("device-list");
const addForm = document.getElementById("add-device-form");
const warning = document.getElementById("limit-warning");
const registerForm = document.getElementById("register-form");

function updateUI() {
    const deviceCount = deviceList.querySelectorAll(".device-box").length;
    if (deviceCount >= MAX_DEVICES) {
    addForm.classList.add("hidden");
    warning.classList.remove("hidden");
    } else {
    addForm.classList.remove("hidden");
    warning.classList.add("hidden");
    }
}

updateUI(); // Run once at page load

registerForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const keyInput = document.getElementById("public-key-input");
    const publicKey = keyInput.value.trim();
    if (!publicKey) return;

    const res = await fetch("/register", {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams({ public_key: publicKey })
    });

    if (res.ok) {
    const newDevice = await res.json(); // expect JSON { id, public_key, ip, active }

    const div = document.createElement("div");
    div.className = "device-box bg-gray-900 border border-gray-800 rounded-xl p-6 flex justify-between items-center";
    div.innerHTML = `
        <div>
        <p class="text-sm text-gray-500 mb-1">Allocated IP</p>
        <p class="text-lg font-mono text-blue-300">${newDevice.ip}</p>
        <p class="mt-3 text-sm">
            <span class="font-semibold text-gray-400">Status:</span>
            ${newDevice.active ? '<span class="text-green-400 font-semibold">Active</span>' : '<span class="text-gray-500 font-semibold">Disconnected</span>'}
        </p>
        </div>
        <form class="remove-form" method="POST" action="/remove">
        <input type="hidden" name="device_id" value="${newDevice.id}" />
        <button class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium">
            Remove
        </button>
        </form>
    `;
    deviceList.appendChild(div);
    keyInput.value = "";
    updateUI();
    } else {
    alert("Failed to add device");
    }
});

deviceList.addEventListener("submit", async (e) => {
    if (!e.target.classList.contains("remove-form")) return;
        e.preventDefault();

        const form = e.target;
        const formData = new FormData(form);
        const res = await fetch("/remove", {
        method: "POST",
        body: formData,
        });

        if (res.ok) {
            form.closest(".device-box").remove();
            updateUI();
        } else {
            alert("Failed to remove device");
        }
    });
});
