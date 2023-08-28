document.addEventListener("DOMContentLoaded", function () {
    const getProtectedResourceButton = document.getElementById("get-protected-resource");
    const protectedResourceDetails = document.getElementById("protected-resource-details");
    const errorDisplay = document.getElementById("error-display");

    getProtectedResourceButton.addEventListener("click", async () => {
        try {
            const response = await fetch("/services"); // Change the URL if needed
            const responseData = await response.clone().json();

            if (response.ok) {
                // Update protected resource details
                const servicesList = document.createElement("ul");
                for (const service of responseData.services) {
                    const serviceItem = document.createElement("li");
                    serviceItem.textContent = service;
                    servicesList.appendChild(serviceItem);
                }
                protectedResourceDetails.appendChild(servicesList);
            } else {
                const errorData = responseData; // Use the cloned response data
                const errorMessage = errorData.error || "An error occurred";
                errorDisplay.textContent = errorMessage;
                protectedResourceDetails.innerHTML = ""; // Clear previous data
            }
        } catch (error) {
            console.error("An error occurred:", error);
            errorDisplay.textContent = "An error occurred while processing your request.";
            protectedResourceDetails.innerHTML = ""; // Clear previous data
        }
    });
});
