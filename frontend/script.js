document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('uploadForm');
    const fileInput = document.getElementById('resumeFile');
    const loading = document.getElementById('loading');
    const resultDiv = document.getElementById('result');

    form.addEventListener('submit', async function (e) { // Prevent default form submission
        e.preventDefault(); // Prevent page reload

        // Show loading spinner
        loading.classList.remove('hidden');
        resultDiv.classList.add('hidden');
        resultDiv.innerHTML = '';

        const formData = new FormData();
        formData.append('resume', fileInput.files[0]);

        try {
            const response = await fetch('http://localhost:8080/upload', {
                method: 'POST',
                body: formData
            });

            const data = await response.json();

            // Hide loading spinner
            loading.classList.add('hidden');
            resultDiv.classList.remove('hidden');

            if (data.error) {
                resultDiv.innerHTML = `<div class="text-red-600">Error: ${data.error}</div>`;
            } else {
                resultDiv.innerHTML = `
                    <div><strong>Name:</strong> ${data.name}</div>
                    <div><strong>Email:</strong> ${data.email}</div>
                    <div><strong>Skills:</strong> ${data.skills.join(', ')}</div>
                    <div class="mt-2"><strong>Raw Text (first 500 chars):</strong><br><pre class="whitespace-pre-wrap">${data.raw_text}</pre></div>
                `;
            }
        } catch (err) {
            loading.classList.add('hidden');
            resultDiv.classList.remove('hidden');
            resultDiv.innerHTML = `<div class="text-red-600">Error: Could not connect to backend.</div>`;
        }
    });
});