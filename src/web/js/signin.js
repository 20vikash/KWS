const sparkleContainer = document.getElementById('sparkle-container');
for (let i = 0; i < 160; i++) {
    const sparkle = document.createElement('div');
    sparkle.classList.add('sparkle');
    sparkle.style.left = Math.random() * 100 + 'vw';
    sparkle.style.top = Math.random() * 100 + 'vh';
    sparkle.style.animationDuration = 6 + Math.random() * 12 + 's';
    sparkle.style.animationDelay = Math.random() * 10 + 's';
    sparkleContainer.appendChild(sparkle);
}

document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("signin-form");
    const errorBox = document.getElementById("error-box");

    function showError(message) {
    errorBox.innerHTML = `• ${message}`;
    errorBox.classList.remove("hidden");
    }

    function hideError() {
    errorBox.classList.add("hidden");
    errorBox.innerHTML = "";
    }

    document.querySelectorAll("input").forEach(input => {
    input.addEventListener("input", hideError);
    });

    form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const formData = new FormData(form);
    const data = new URLSearchParams(formData);

    try {
        const res = await fetch("/login", {
        method: "POST",
        body: data
        });

        const text = await res.text();

        if (res.ok) {
        window.location.href = "/";
        } else {
        showError(text);
        }
    } catch {
        showError("Network error. Please try again.");
    }
    });
});
