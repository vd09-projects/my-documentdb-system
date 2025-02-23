// js/main.js
document.addEventListener("DOMContentLoaded", () => {
  // Cached DOM elements
  const navElements = {
    login: document.getElementById("navLogin"),
    register: document.getElementById("navRegister"),
    upload: document.getElementById("navUpload"),
    dashboard: document.getElementById("navDashboard"), // DASHBOARD CHANGES
    aggregator: document.getElementById("navAggregator"),
    logout: document.getElementById("navLogout"),
  };

  const sections = {
    login: document.getElementById("loginSection"),
    register: document.getElementById("registerSection"),
    upload: document.getElementById("uploadSection"),
    aggregator: document.getElementById("aggregatorSection"),
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

      // Aggregator elements
      const aggregatorElements = {
        recordType: document.getElementById("aggRecordType"),
        field: document.getElementById("aggField"),
        operation: document.getElementById("aggOperation"),
        calculateButton: document.getElementById("calculateButton"),
        resultDiv: document.getElementById("aggResult"),
      };

  // Utility function to update UI based on login state
  function updateUI() {
    const token = localStorage.getItem("jwtToken");

    if (token) {
      navElements.logout.style.display = "inline-block";
      navElements.upload.style.display = "inline-block";
      navElements.dashboard.style.display = "inline-block"; // Show Dashboard if logged in
      navElements.aggregator.style.display = "inline-block";
      navElements.login.style.display = "none";
      navElements.register.style.display = "none";
    } else {
      navElements.login.style.display = "inline-block";
      navElements.register.style.display = "inline-block";
      navElements.upload.style.display = "none";
      navElements.dashboard.style.display = "none"; // Hide Dashboard if not logged in
      navElements.aggregator.style.display = "none";
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
      case "aggregator": // NEW aggregator case
        handleAggregatorNav();
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

  function handleAggregatorNav() {
    const token = localStorage.getItem("jwtToken");
    if (!token) {
      alert("Please log in first!");
      return;
    }
    showSection("aggregator");
    aggregatorElements.resultDiv.textContent = "";

    // Load record types so user can pick them
    loadRecordTypes();
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
  
    // 1. Get the chosen record type
    const recordTypeSelect = document.getElementById("recordType");
    const selectedRecordType = recordTypeSelect.value;
  
    // 2. Prepare the form data
    const formData = new FormData(forms.upload);
    formData.append("recordType", selectedRecordType); // Append extra field
  
    try {
      // 3. Upload with record type
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

  async function loadRecordTypes() {
    console.log("loadRecordTypes")
    const token = localStorage.getItem("jwtToken"); // Get the JWT token from localStorage
  
    try {
      // Fetch record types from the API
      const res = await fetch("/listRecordTypes", {
        headers: { Authorization: `Bearer ${token}` },
      });
  
      if (!res.ok) {
        throw new Error(`Failed to fetch record types: ${res.status}`);
      }
  
      const types = await res.json(); // Assuming the response is an array of record types
  
      // Clear existing options in the dropdown
      aggregatorElements.recordType.innerHTML = '<option value="" disabled selected>Select a record type</option>';
  
      // Populate the dropdown with the received record types
      types.forEach((type) => {
        const option = document.createElement("option");
        option.value = type.recordType; // Use `type` as the value
        option.textContent = type.recordType; // Display `type` as the text
        aggregatorElements.recordType.appendChild(option);
      });
  
      console.log("Record types loaded successfully.");
    } catch (error) {
      console.error("Error loading record types:", error);
      alert("Failed to load record types. Please try again.");
    }
  }

  // Load fields for the selected record type
  async function loadFieldsForType(type) {
    const token = localStorage.getItem("jwtToken");
    if (!token) return;

    try {
      const res = await fetch(`/listFields?recordType=${type}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) throw new Error(await res.text());

      const fields = await res.json(); // e.g. ["price", "quantity"]
      aggregatorElements.field.innerHTML = "";
      console.log(fields)
      fields.fields.forEach((f) => {
        const opt = document.createElement("option");
        opt.value = f;
        opt.textContent = f;
        aggregatorElements.field.appendChild(opt);
      });
    } catch (err) {
      console.error("Error loading fields:", err.message);
    }
  }

  // Perform sum/average on chosen field
  async function handleCalculate() {
    const token = localStorage.getItem("jwtToken");
    if (!token) {
      alert("You are not logged in!");
      return;
    }

    const recordType = aggregatorElements.recordType.value;
    const field = aggregatorElements.field.value;
    const op = aggregatorElements.operation.value;

    const url = `/aggregate?recordType=${recordType}&field=${field}&op=${op}`;
    try {
      const res = await fetch(url, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) throw new Error(await res.text());

      const data = await res.json(); // { result: someNumber }
      aggregatorElements.resultDiv.textContent = `Result: ${data.result}`;
    } catch (err) {
      aggregatorElements.resultDiv.textContent = `Error: ${err.message}`;
    }
  }

  // Attach event listeners
  function attachEventListeners() {
    navElements.login.addEventListener("click", () => handleNavClick("login"));
    navElements.register.addEventListener("click", () => handleNavClick("register"));
    navElements.upload.addEventListener("click", () => handleNavClick("upload"));
    navElements.dashboard.addEventListener("click", () => handleNavClick("dashboard")); // DASHBOARD CHANGES
    navElements.aggregator.addEventListener("click", () => handleNavClick("aggregator"));
    navElements.logout.addEventListener("click", () => handleNavClick("logout"));

    forms.login.addEventListener("submit", handleLoginSubmit);
    forms.register.addEventListener("submit", handleRegisterSubmit);
    forms.upload.addEventListener("submit", handleUploadSubmit);

    // DASHBOARD CHANGES: button to fetch data
    dashboardElements.getDataButton.addEventListener("click", handleGetDataClick);

    // Aggregator
    if (aggregatorElements.recordType) {
      aggregatorElements.recordType.addEventListener("change", () => {
        loadFieldsForType(aggregatorElements.recordType.value);
      });
    }
    aggregatorElements.calculateButton.addEventListener("click", handleCalculate);
  }

  // Initialize application
  function initializeApp() {
    updateUI();
    attachEventListeners();
  }

  initializeApp();
});