
// ===============================
// Scan Repository
// Called when user clicks Scan Repository
// ===============================

async function scanRepo() {

    // Get repository URL
    const repoUrl =
        document
            .getElementById("repoUrl")
            .value
            .trim();

    // Validate input
    if (!repoUrl) {

        alert("Please enter a GitHub Repository URL.");

        return;
    }

    // Show loading spinner
    document
        .getElementById("loading")
        .style
        .display = "block";

    // Remove previous result
    document
        .getElementById("result")
        .innerHTML = "";

    try {

        // Call backend
        const response = await fetch("/scan", {

            method: "POST",

            headers: {

                "Content-Type": "application/json"

            },

            body: JSON.stringify({

                repo_url: repoUrl

            })

        });

        if (!response.ok) {

            throw new Error("Scan failed.");

        }

        const data = await response.json();

        // Hide spinner
        document
            .getElementById("loading")
            .style
            .display = "none";

        // Build complete dashboard
        buildDashboard(data);

        // Draw chart
        createChart(data);

    }

    catch (error) {

        document
            .getElementById("loading")
            .style
            .display = "none";

        document
            .getElementById("result")
            .innerHTML =

        `
        <div class="alert alert-danger">

            ${error.message}

        </div>
        `;

    }

}

// =========================================
// Build the complete dashboard
// =========================================

function buildDashboard(data) {

    // Decide risk alert color
    alert("buildDashboard called");
     console.log("Dashboard Loaded");
    let riskClass = "success";

    if (data.risk_level === "MEDIUM")
        riskClass = "warning";

    else if (data.risk_level === "HIGH")
        riskClass = "danger";

    else if (data.risk_level === "CRITICAL")
        riskClass = "danger";


    let html = `

<div class="card shadow p-4">

    <div class="d-flex justify-content-between align-items-center mb-3">

    <h2>Scan Results</h2>

    <div>

        <a
            href="/report"
            class="btn btn-success me-2">
            Download HTML
        </a>

        <a
            href="/report/pdf"
            class="btn btn-danger">
            Download PDF
        </a>

    </div>

</div>

<hr>

<div class="row mb-4">

    <div class="col-md-8">

        <input
            type="email"
            id="email"
            class="form-control"
            placeholder="Enter recipient email">

    </div>

    <div class="col-md-4">

        <button
            class="btn btn-primary w-100"
            onclick="sendReport()">

            Send PDF

        </button>

    </div>

</div>


    <div class="alert alert-${riskClass}">

        <strong>Risk Level:</strong>

        ${data.risk_level}

    </div>


    <div class="row text-center mb-4">

        <div class="col-md-3">

            <div class="card shadow-sm p-3">

                <h5>Risk Score</h5>

                <h2>${data.risk_score}</h2>

            </div>

        </div>

        <div class="col-md-3">

            <div class="card shadow-sm p-3">

                <h5>Errors</h5>

                <h2>${data.error_count}</h2>

            </div>

        </div>

        <div class="col-md-3">

            <div class="card shadow-sm p-3">

                <h5>Warnings</h5>

                <h2>${data.warning_count}</h2>

            </div>

        </div>

        <div class="col-md-3">

            <div class="card shadow-sm p-3">

                <h5>Info</h5>

                <h2>${data.info_count}</h2>

            </div>

        </div>

    </div>


    <hr>


    <h4 class="mb-3">

        Severity Distribution

    </h4>

    <div class="text-center mb-4">

        <div style="max-width:350px; margin:auto;">

            <canvas id="severityChart"></canvas>

        </div>

    </div>


    <hr>


    <h4>

        Vulnerabilities

    </h4>


    <div class="row mb-3">

        <div class="col-md-8">

            <input

                type="text"

                id="searchInput"

                class="form-control"

                placeholder="Search by vulnerability type or file..."

                onkeyup="filterTable()">

        </div>


        <div class="col-md-4">

            <select

                id="severityFilter"

                class="form-select"

                onchange="filterTable()">

                <option value="">All Severities</option>

                <option value="ERROR">ERROR</option>

                <option value="WARNING">WARNING</option>

                <option value="INFO">INFO</option>

            </select>

        </div>

    </div>


    ${buildTable(data.vulnerabilities)}

</div>

`;

    document.getElementById("result").innerHTML = html;

}

// =========================================
// Build Vulnerability Table
// =========================================

function buildTable(vulnerabilities) {

    let table = `

<div class="table-responsive">

<table class="table table-hover table-striped">

    <thead class="table-dark">

        <tr>

            <th>Severity</th>
            <th>Type</th>
            <th>File</th>
            <th>Line</th>
            <th>Description</th>

        </tr>

    </thead>

    <tbody>

`;

    vulnerabilities.forEach((v, index) => {

        table += `

<tr

    class="vuln-row"

    data-type="${v.type}"

    data-file="${v.file}"

    data-severity="${v.severity}">

    <td>

        <span class="badge ${getSeverityClass(v.severity)}">

            ${v.severity}

        </span>

    </td>

    <td>

        ${v.type}

    </td>

    <td>

        ${v.file}

    </td>

    <td>

        ${v.line}

    </td>

    <td>

        <button

            class="btn btn-sm btn-outline-primary"

            data-bs-toggle="collapse"

            data-bs-target="#desc${index}">

            View Details

        </button>

        <div

            id="desc${index}"

            class="collapse mt-2">

            <div class="card card-body">

                ${v.description}

            </div>

        </div>

    </td>

</tr>

`;

    });

    table += `

    </tbody>

</table>

</div>

`;

    return table;

}
// =========================================
// Return CSS class for severity badge
// =========================================

function getSeverityClass(severity) {

    switch (severity) {

        case "ERROR":
            return "severity-error";

        case "WARNING":
            return "severity-warning";

        case "INFO":
            return "severity-info";

        default:
            return "bg-secondary";

    }

}

// =========================================
// Create Severity Pie Chart
// =========================================

function createChart(data) {

    const canvas =
        document.getElementById("severityChart");

    // If there is no canvas, do nothing
    if (!canvas)
        return;

    // Destroy old chart before creating a new one
    

    window.severityChart = new Chart(canvas, {

        type: "pie",

        data: {

            labels: [
                "Errors",
                "Warnings",
                "Info"
            ],

            datasets: [{

                data: [

                    data.error_count,
                    data.warning_count,
                    data.info_count

                ],

                backgroundColor: [

                    "#dc3545",   // Error - Red
                    "#fd7e14",   // Warning - Orange
                    "#0d6efd"    // Info - Blue

                ],

                borderWidth: 1

            }]

        },

        options: {

            responsive: true,

            maintainAspectRatio: true,

            plugins: {

                legend: {

                    position: "right"

                }

            }

        }

    });

}

// =========================================
// Search + Severity Filter
// =========================================

function filterTable() {

    const search =
        document
            .getElementById("searchInput")
            .value
            .toLowerCase();

    const severity =
        document
            .getElementById("severityFilter")
            .value;

    const rows =
        document.querySelectorAll(".vuln-row");

    rows.forEach(row => {

        const type =
            row.dataset.type.toLowerCase();

        const file =
            row.dataset.file.toLowerCase();

        const rowSeverity =
            row.dataset.severity;

        const matchesSearch =
            type.includes(search) ||
            file.includes(search);

        const matchesSeverity =
            severity === "" ||
            rowSeverity === severity;

        if (matchesSearch && matchesSeverity) {

            row.style.display = "";

        } else {

            row.style.display = "none";

        }

    });

}

async function sendReport() {

    const email =
        document.getElementById("email").value;

    if (!email) {

        alert("Please enter an email address");
        return;
    }

    const response = await fetch("/send-email", {

        method: "POST",

        headers: {
            "Content-Type": "application/json"
        },

        body: JSON.stringify({

            email: email

        })

    });

    const data = await response.json();

    alert(data.message);

}
