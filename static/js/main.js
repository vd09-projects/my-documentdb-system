// Production-Ready Code with Modularized Functions
document.addEventListener("DOMContentLoaded", () => {
  // Cached DOM elements
  const navElements = {
    login: document.getElementById("navLogin"),
    register: document.getElementById("navRegister"),
    upload: document.getElementById("navUpload"),
    logout: document.getElementById("navLogout"),
  };

  const sections = {
    login: document.getElementById("loginSection"),
    register: document.getElementById("registerSection"),
    upload: document.getElementById("uploadSection"),
  };

  const forms = {
    login: document.getElementById("loginForm"),
    register: document.getElementById("registerForm"),
    upload: document.getElementById("uploadForm"),
  };

  const messages = {
    login: document.getElementById("loginMessage"),
    register: document.getElementById("registerMessage"),
    uploadResponse: document.getElementById("uploadResponse"),
    uploadInfo: document.getElementById("uploadInfo"),
  };

  // Utility function to update UI based on login state
  function updateUI() {
    const token = localStorage.getItem("jwtToken");

    if (token) {
      navElements.logout.style.display = "inline-block";
      navElements.upload.style.display = "inline-block";
      navElements.login.style.display = "none";
      navElements.register.style.display = "none";
    } else {
      navElements.login.style.display = "inline-block";
      navElements.register.style.display = "inline-block";
      navElements.upload.style.display = "none";
      navElements.logout.style.display = "none";
    }
  }

  // Show the appropriate section and hide others
  function showSection(sectionName) {
    Object.keys(sections).forEach((key) => {
      sections[key].style.display = key === sectionName ? "block" : "none";
    });
  }

  // Handle navigation click events
  function handleNavClick(navType) {
    switch (navType) {
      case "login":
        showSection("login");
        break;
      case "register":
        showSection("register");
        break;
      case "upload":
        handleUploadNav();
        break;
      case "logout":
        handleLogout();
        break;
      default:
        console.error("Unknown navigation type");
    }
  }

  function handleUploadNav() {
    const token = localStorage.getItem("jwtToken");
    if (!token) {
      alert("Please log in first!");
      return;
    }
    showSection("upload");
    messages.uploadInfo.textContent = `You are logged in. Token: ${token}`;
  }

  function handleLogout() {
    localStorage.removeItem("jwtToken");
    alert("You have been logged out!");
    updateUI();
    showSection("login");
  }

  // Handle form submissions
  async function handleLoginSubmit(event) {
    event.preventDefault();
    messages.login.textContent = "";

    const username = forms.login.loginUsername.value.trim();
    const password = forms.login.loginPassword.value.trim();

    try {
      const response = await fetch("/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) throw new Error(await response.text());

      const token = await response.text();
      localStorage.setItem("jwtToken", token);
      messages.login.textContent = "Login successful!";
      updateUI();
    } catch (error) {
      messages.login.textContent = `Login failed: ${error.message}`;
    }
  }

  async function handleRegisterSubmit(event) {
    event.preventDefault();
    messages.register.textContent = "";

    const username = forms.register.registerUsername.value.trim();
    const password = forms.register.registerPassword.value.trim();

    try {
      const response = await fetch("/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) throw new Error(await response.text());

      const message = await response.text();
      messages.register.textContent = message;
    } catch (error) {
      messages.register.textContent = `Registration failed: ${error.message}`;
    }
  }

  async function handleUploadSubmit(event) {
    event.preventDefault();
    messages.uploadResponse.textContent = "";

    const token = localStorage.getItem("jwtToken");
    if (!token) {
      messages.uploadResponse.textContent = "You are not logged in!";
      return;
    }

    const formData = new FormData(forms.upload);

    try {
      const response = await fetch("/upload", {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
      });

      const text = await response.text();
      if (!response.ok) throw new Error(text);

      messages.uploadResponse.textContent = text;
    } catch (error) {
      messages.uploadResponse.textContent = `Upload failed: ${error.message}`;
    }
  }

  // Attach event listeners
  function attachEventListeners() {
    navElements.login.addEventListener("click", () => handleNavClick("login"));
    navElements.register.addEventListener("click", () => handleNavClick("register"));
    navElements.upload.addEventListener("click", () => handleNavClick("upload"));
    navElements.logout.addEventListener("click", () => handleNavClick("logout"));

    forms.login.addEventListener("submit", handleLoginSubmit);
    forms.register.addEventListener("submit", handleRegisterSubmit);
    forms.upload.addEventListener("submit", handleUploadSubmit);
  }

  // Initialize application
  function initializeApp() {
    updateUI();
    attachEventListeners();
  }

  initializeApp();
});
