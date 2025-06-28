document.addEventListener("DOMContentLoaded", () => {
    const form = document.querySelector("form");
    const errorBox = document.getElementById("error-box");
    const errorMessage = document.getElementById("error-message");
    const passwordInput = document.getElementById("password");
    const confirmPasswordInput = document.querySelector('input[name="confirm_password"]');
    const passwordStrength = document.getElementById("password-strength");

    const emailRegex = /^[a-zA-Z0-9._%+\-]+@gmail\.com$/i;
    const usernameRegex = /^[a-zA-Z][a-zA-Z0-9._]{4,29}$/;
    const nameRegex = /^[A-Za-z'.-]{2,30}$/;
    const passwordRegex = {
        length: /.{8,}/,
        lower: /[a-z]/,
        upper: /[A-Z]/,
        digit: /[0-9]/,
        special: /[!@#$%^&*()\-_=+\[\]{}|;:'",.<>/?]/
    };

    function showError(message) {
        errorMessage.textContent = message;
        errorBox.classList.remove("hidden");
    }

    function hideError() {
        errorBox.classList.add("hidden");
        errorMessage.textContent = "";
    }

    // Clear error box on input change
    document.querySelectorAll("input").forEach(input => {
        input.addEventListener("input", () => {
            hideError();
            
            // Update password strength indicator
            if (input === passwordInput) {
                updatePasswordStrength();
            }
        });
    });

    // Password strength indicator
    function updatePasswordStrength() {
        const password = passwordInput.value;
        let strength = 0;
        
        if (password.length >= 8) strength += 1;
        if (/[A-Z]/.test(password)) strength += 1;
        if (/[0-9]/.test(password)) strength += 1;
        if (/[^A-Za-z0-9]/.test(password)) strength += 1;
        
        passwordStrength.className = 'password-strength password-strength-' + strength;
    }

    form.addEventListener("submit", async (e) => {
        e.preventDefault(); // prevent default form submission

        const email = form.email.value.trim();
        const username = form.user_name.value.trim();
        const first = form.first_name.value.trim();
        const last = form.last_name.value.trim();
        const password = form.password.value;
        const confirmPassword = form.confirm_password.value;

        const errors = [];

        if (!emailRegex.test(email)) {
            errors.push("• Email must be a valid Gmail address.");
        }

        if (!usernameRegex.test(username)) {
            errors.push("• Username must start with a letter, 5–30 chars, using letters/numbers/._ only.");
        }

        if (!nameRegex.test(first)) {
            errors.push("• First name must be 2–30 chars with valid characters.");
        }

        if (!nameRegex.test(last)) {
            errors.push("• Last name must be 2–30 chars with valid characters.");
        }

        if (
            !passwordRegex.length.test(password) ||
            !passwordRegex.lower.test(password) ||
            !passwordRegex.upper.test(password) ||
            !passwordRegex.digit.test(password) ||
            !passwordRegex.special.test(password)
        ) {
            errors.push("• Password must be 8+ characters with uppercase, lowercase, number, and special character.");
        }

        if (password !== confirmPassword) {
            errors.push("• Passwords do not match.");
        }

        // Frontend validation failed
        if (errors.length > 0) {
            showError(errors.join("\n"));
            return;
        }

        // Prepare and send POST request
        const formData = new FormData(form);
        const data = new URLSearchParams(formData);

        try {
            const res = await fetch("/create_user", {
                method: "POST",
                body: data,
            });

            const text = await res.text();

            if (res.ok) {
                // Redirect to login
                window.location.href = "/kws_signin";
            } else {
                // Server-side error (like "username or email already exists")
                showError(`• ${text}`);
            }
        } catch (err) {
            showError("• Network error. Please try again.");
        }
    });
});
