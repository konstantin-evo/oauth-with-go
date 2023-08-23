document.addEventListener("DOMContentLoaded", function () {
    const getProtectedResourceButton = document.getElementById("get-protected-resource");
    const protectedResourceDetails = document.getElementById("protected-resource-details");

    getProtectedResourceButton.addEventListener("click", async () => {
        try {
            const response = await fetch("/services"); // Change the URL if needed
            const data = await response.json();

            if (response.ok) {
                // Update protected resource details
                const servicesList = document.createElement("ul");
                for (const service of data.services) {
                    const serviceItem = document.createElement("li");
                    serviceItem.textContent = service;
                    servicesList.appendChild(serviceItem);
                }
                protectedResourceDetails.appendChild(servicesList);
            } else {
                console.error("Error fetching protected resource data");
            }
        } catch (error) {
            console.error("An error occurred:", error);
        }
    });
});
