// Sparkle background animation generator
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
    const form = document.querySelector("form");
    const errorBox = document.getElementById("error-box");

    // Regex patterns
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

    function showErrors(messages) {
    errorBox.innerHTML = messages.join("<br>");
    errorBox.classList.remove("hidden", "opacity-0");
    errorBox.classList.add("animate-fadeIn");
    }

    function hideErrors() {
    errorBox.classList.add("hidden");
    errorBox.innerHTML = "";
    }

    // Clear error box on input change
    document.querySelectorAll("input").forEach(input => {
    input.addEventListener("input", hideErrors);
    });

    form.addEventListener("submit", async (e) => {
    e.preventDefault(); // prevent default form submission

    const email = form.email.value.trim();
    const username = form.user_name.value.trim();
    const first = form.first_name.value.trim();
    const last = form.last_name.value.trim();
    const password = form.password.value;

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

    // Frontend validation failed
    if (errors.length > 0) {
        showErrors(errors);
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
        // Redirect to login or dashboard
        window.location.href = "/kws_signin";
        } else {
        // Server-side error (like "username or email already exists")
        showErrors([`• ${text}`]);
        }
    } catch (err) {
        showErrors(["• Network error. Please try again."]);
    }
    });
});
