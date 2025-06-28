document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("signin-form");
    const errorBox = document.getElementById("error-box");
    const errorMessage = document.getElementById("error-message");

    function showError(message) {
        errorMessage.textContent = message;
        errorBox.classList.remove("hidden");
    }

    function hideError() {
        errorBox.classList.add("hidden");
        errorMessage.textContent = "";
    }

    // Hide error when any input is modified
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
                // Successful login - redirect to home
                window.location.href = "/";
            } else {
                // Show error message from server
                showError(text);
            }
        } catch (error) {
            // Handle network errors
            showError("Network error. Please check your connection and try again.");
            console.error("Login error:", error);
        }
    });
});