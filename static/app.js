async function scanRepo() {

    const repoUrl =
        document.getElementById("repoUrl").value;

    if (!repoUrl) {
        alert("Please enter repository URL");
        return;
    }

    document.getElementById("loading").style.display =
        "block";

    document.getElementById("result").innerHTML = "";

    try {

        const response = await fetch("/scan", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                repo_url: repoUrl
            })
        });

        const data = await response.json();

        document.getElementById("loading").style.display =
            "none";

        let riskClass = "success";

        if (data.risk_level === "MEDIUM")
            riskClass = "warning";

        if (data.risk_level === "HIGH")
            riskClass = "danger";

        if (data.risk_level === "CRITICAL")
            riskClass = "danger";

        let html = `

        <div class="card shadow p-4">

            <h2 class="mb-3">
                Scan Results
            </h2>

            <div class="alert alert-${riskClass}">
                <strong>Risk Level:</strong>
                ${data.risk_level}
            </div>

            <div class="row text-center mb-4">

                <div class="col-md-3">
                    <div class="card p-3 shadow-sm">
                        <h5>Risk Score</h5>
                        <h2>${data.risk_score}</h2>
                    </div>
                </div>

                <div class="col-md-3">
                    <div class="card p-3 shadow-sm">
                        <h5>Errors</h5>
                        <h2>${data.error_count}</h2>
                    </div>
                </div>

                <div class="col-md-3">
                    <div class="card p-3 shadow-sm">
                        <h5>Warnings</h5>
                        <h2>${data.warning_count}</h2>
                    </div>
                </div>

                <div class="col-md-3">
                    <div class="card p-3 shadow-sm">
                        <h5>Info</h5>
                        <h2>${data.info_count}</h2>
                    </div>
                </div>

            </div>

            <hr>

            <h4 class="mb-3">
                Severity Distribution
            </h4>

           <div class="mb-4 text-center">
    <div style="max-width:350px; margin:auto;">
        <canvas id="severityChart"></canvas>
    </div>
</div>

            <hr>

            <h4 class="mb-3">
                Vulnerabilities
            </h4>

            <div class="table-container">

                <table class="table table-striped table-hover">

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

    data.vulnerabilities.forEach((v, index) => {

    let severityClass = "";

    if (v.severity === "ERROR")
        severityClass = "severity-error";

    if (v.severity === "WARNING")
        severityClass = "severity-warning";

    if (v.severity === "INFO")
        severityClass = "severity-info";

    html += `
    <tr>
        <td>
            <span class="badge ${severityClass}">
                ${v.severity}
            </span>
        </td>

        <td>${v.type}</td>
        <td>${v.file}</td>
        <td>${v.line}</td>

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

        html += `
                    </tbody>
                </table>

            </div>

        </div>
        `;

        document.getElementById("result").innerHTML =
            html;

        // Create Pie Chart
        const ctx = document.getElementById("severityChart");

    new Chart(ctx, {
    type: "pie",
    data: {
        labels: ["Errors", "Warnings", "Info"],
        datasets: [{
            data: [
                data.error_count,
                data.warning_count,
                data.info_count
            ],
            backgroundColor: [
                "#dc3545", // Red (Error)
                "#fd7e14", // Orange (Warning)
                "#0d6efd"  // Blue (Info)
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
    } catch (error) {

        document.getElementById("loading").style.display =
            "none";

        document.getElementById("result").innerHTML =
            `
            <div class="alert alert-danger">
                ${error}
            </div>
            `;
    }
}