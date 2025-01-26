// js/main.js
document.addEventListener("DOMContentLoaded", () => {
  // Cached DOM elements
  const navElements = {
    login: document.getElementById("navLogin"),
    register: document.getElementById("navRegister"),
    upload: document.getElementById("navUpload"),
    dashboard: document.getElementById("navDashboard"), // DASHBOARD CHANGES
    logout: document.getElementById("navLogout"),
  };

  const sections = {
    login: document.getElementById("loginSection"),
    register: document.getElementById("registerSection"),
    upload: document.getElementById("uploadSection"),
    dashboard: document.getElementById("dashboardSection"), // DASHBOARD CHANGES
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

  // DASHBOARD CHANGES
  const dashboardElements = {
    fromDate: document.getElementById("fromDate"),
    toDate: document.getElementById("toDate"),
    getDataButton: document.getElementById("getDataButton"),
    resultsDiv: document.getElementById("dashboardResults"),
  };

  // Utility function to update UI based on login state
  function updateUI() {
    const token = localStorage.getItem("jwtToken");

    if (token) {
      navElements.logout.style.display = "inline-block";
      navElements.upload.style.display = "inline-block";
      navElements.dashboard.style.display = "inline-block"; // Show Dashboard if logged in
      navElements.login.style.display = "none";
      navElements.register.style.display = "none";
    } else {
      navElements.login.style.display = "inline-block";
      navElements.register.style.display = "inline-block";
      navElements.upload.style.display = "none";
      navElements.dashboard.style.display = "none"; // Hide Dashboard if not logged in
      navElements.logout.style.display = "none";
    }
  }

  // Show the appropriate section and hide others
  function showSection(sectionName) {
    Object.keys(sections).forEach((key) => {
      sections[key].style.display = key === sectionName ? "block" : "none";
    });
  }

  // Navigation click events
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
      case "dashboard": // DASHBOARD CHANGES
        handleDashboardNav();
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

  // DASHBOARD CHANGES: Show the Dashboard section
  function handleDashboardNav() {
    const token = localStorage.getItem("jwtToken");
    if (!token) {
      alert("Please log in first!");
      return;
    }
    showSection("dashboard");
    dashboardElements.resultsDiv.innerHTML = ""; // clear old results
  }

  function handleLogout() {
    localStorage.removeItem("jwtToken");
    alert("You have been logged out!");
    updateUI();
    showSection("login");
  }

  // --- Forms: LOGIN ---
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

  // --- Forms: REGISTER ---
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

  // --- Forms: UPLOAD ---
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

  // --- DASHBOARD CHANGES: GET DATA ---
  async function handleGetDataClick() {
    // Grab optional date strings
    const from = dashboardElements.fromDate.value.trim();
    const to = dashboardElements.toDate.value.trim();

    // Build the URL with optional query params
    let url = "/userData";
    const params = [];
    if (from) params.push(`from=${from}`);
    if (to) params.push(`to=${to}`);
    if (params.length > 0) {
      url += "?" + params.join("&");
    }

    const token = localStorage.getItem("jwtToken");
    if (!token) {
      alert("You are not logged in!");
      return;
    }

    try {
      const res = await fetch(url, {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (!res.ok) {
        throw new Error(await res.text());
      }
      const data = await res.json();
      displayDashboardData(data);
    } catch (err) {
      console.error("Error fetching user data:", err.message);
      dashboardElements.resultsDiv.innerHTML = `<p style="color:red;">${err.message}</p>`;
    }
  }

  // --- DASHBOARD CHANGES: Display results ---
  function displayDashboardData(data) {
    // For a quick demonstration, we can just JSON.stringify the data
    // or build a table, etc.
    if (!data || data.length === 0) {
      dashboardElements.resultsDiv.innerHTML = "<p>No records found.</p>";
      return;
    }

    // Example: A minimal approach
    dashboardElements.resultsDiv.innerHTML = `<pre>${JSON.stringify(data, null, 2)}</pre>`;
  }

  // Attach event listeners
  function attachEventListeners() {
    navElements.login.addEventListener("click", () => handleNavClick("login"));
    navElements.register.addEventListener("click", () => handleNavClick("register"));
    navElements.upload.addEventListener("click", () => handleNavClick("upload"));
    navElements.dashboard.addEventListener("click", () => handleNavClick("dashboard")); // DASHBOARD CHANGES
    navElements.logout.addEventListener("click", () => handleNavClick("logout"));

    forms.login.addEventListener("submit", handleLoginSubmit);
    forms.register.addEventListener("submit", handleRegisterSubmit);
    forms.upload.addEventListener("submit", handleUploadSubmit);

    // DASHBOARD CHANGES: button to fetch data
    dashboardElements.getDataButton.addEventListener("click", handleGetDataClick);
  }

  // Initialize application
  function initializeApp() {
    updateUI();
    attachEventListeners();
  }

  initializeApp();
});