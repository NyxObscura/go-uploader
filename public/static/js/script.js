document.addEventListener('DOMContentLoaded', function () {
    // --- DOM Element Variables ---
    const form = document.getElementById('uploadForm');
    const fileInput = document.getElementById('fileInput');
    const dropZone = document.getElementById('dropZone');
    const resultDiv = document.getElementById('result');
    const progressBarContainer = document.getElementById('progressBarContainer');
    const progressBarFill = document.getElementById('progressBarFill');
    const uploadPlaceholder = document.getElementById('uploadPlaceholder');
    const imagePreview = document.getElementById('imagePreview');
    const submitButton = form.querySelector('button[type="submit"]');
    const darkModeToggle = document.getElementById('darkModeToggle');
    const toast = document.getElementById('toast');

    // --- Dark Mode Logic ---
    if (localStorage.getItem('color-theme') === 'dark' || (!('color-theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
        document.documentElement.classList.add('dark');
        darkModeToggle.textContent = '☀️';
    } else {
        document.documentElement.classList.remove('dark');
        darkModeToggle.textContent = '🌙';
    }
    darkModeToggle.addEventListener('click', () => {
        const isDark = document.documentElement.classList.toggle('dark');
        localStorage.setItem('color-theme', isDark ? 'dark' : 'light');
        darkModeToggle.textContent = isDark ? '☀️' : '🌙';
    });

    // --- Toast Notification Logic ---
    function showToast(message) {
        toast.textContent = message;
        toast.classList.remove('opacity-0', 'translate-y-[-20px]');
        toast.classList.add('opacity-100', 'translate-y-0');
        setTimeout(() => {
            toast.classList.remove('opacity-100', 'translate-y-0');
            toast.classList.add('opacity-0', 'translate-y-[-20px]');
        }, 3000);
    }

    // --- Drag and Drop Logic ---
    function preventDefaults(e) { e.preventDefault(); e.stopPropagation(); }
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => dropZone.addEventListener(eventName, preventDefaults, false));
    ['dragenter', 'dragover'].forEach(eventName => dropZone.addEventListener(eventName, () => dropZone.classList.add('drag-over'), false));
    ['dragleave', 'drop'].forEach(eventName => dropZone.addEventListener(eventName, () => dropZone.classList.remove('drag-over'), false));
    dropZone.addEventListener('drop', (e) => handleFiles(e.dataTransfer.files), false);
    dropZone.addEventListener('click', () => fileInput.click());
    fileInput.addEventListener('change', () => handleFiles(fileInput.files));

    function handleFiles(files) {
        if (files.length === 0) return;
        fileInput.files = files;
        const file = files[0];
        resetResult();
        uploadPlaceholder.classList.add('hidden');
        imagePreview.classList.remove('hidden');
        imagePreview.innerHTML = `<p class="font-medium break-all">${file.name}</p>`;
        if (file.type.startsWith('image/')) {
            const reader = new FileReader();
            reader.onload = (e) => {
                const img = document.createElement('img');
                img.src = e.target.result;
                img.className = 'max-h-24 mx-auto mt-2 rounded-md object-contain';
                imagePreview.appendChild(img);
            };
            reader.readAsDataURL(file);
        }
    }
    
    // --- Form Submission Logic ---
    form.addEventListener('submit', function(e) {
        e.preventDefault();
        const file = fileInput.files[0];
        if (!file) {
            showResult('error', 'Please select a file first.');
            return;
        }

        // THE FIX IS HERE!
        // This automatically creates the form data from your HTML form,
        // correctly including the file with the name "file".
        const formData = new FormData(form);
        
        const xhr = new XMLHttpRequest();
        xhr.open('POST', '/upload', true);
        resetResult();
        submitButton.disabled = true;
        submitButton.textContent = 'Uploading...';
        progressBarContainer.classList.remove('hidden');
        progressBarFill.style.width = '0%';
        xhr.upload.onprogress = (e) => {
            if (e.lengthComputable) progressBarFill.style.width = `${(e.loaded / e.total) * 100}%`;
        };
        xhr.onload = () => {
            try {
                const response = JSON.parse(xhr.responseText);
                if (xhr.status >= 200 && xhr.status < 300) {
                    showResult('success', response.message, response.data);
                    progressBarFill.style.width = '100%';
                } else {
                    showResult('error', response.message);
                }
            } catch (err) {
                showResult('error', 'Failed to process server response.');
            } finally {
                resetFormState();
            }
        };
        xhr.onerror = () => {
            showResult('error', 'A network error occurred. Please check your connection.');
            resetFormState();
        };
        xhr.send(formData);
    });

    // --- UI Helper Functions ---
    function resetFormState() {
        submitButton.disabled = false;
        submitButton.textContent = 'Upload File';
    }
    function resetResult() {
        resultDiv.innerHTML = '';
        resultDiv.className = '';
        progressBarContainer.classList.add('hidden');
        progressBarFill.style.width = '0%';
    }
    function showResult(type, message, responseData = {}) {
        resultDiv.innerHTML = '';
        let content = '';
        if (type === 'success') {
            const successPath = responseData.path;
            resultDiv.className = 'p-4 mt-4 rounded-md bg-green-100 dark:bg-green-800 border border-green-200 dark:border-green-700';
            content = `
                <p class="font-bold text-green-800 dark:text-green-100">${message}</p>
                <div class="flex items-center mt-2">
                    <input id="filePath" class="w-full bg-slate-100 dark:bg-slate-700 text-sm p-2 rounded-l-md" type="text" value="${successPath}" readonly>
                    <button id="copyButton" class="bg-blue-600 text-white p-2 rounded-r-md hover:bg-blue-700">Copy</button>
                </div>`;
        } else {
            resultDiv.className = 'p-4 mt-4 rounded-md bg-red-100 dark:bg-red-800 border border-red-200 dark:border-red-700';
            content = `<p class="font-bold text-red-800 dark:text-red-100">Error!</p><p class="text-red-700 dark:text-red-200">${message}</p>`;
        }
        resultDiv.innerHTML = content;
        if (type === 'success') {
            document.getElementById('copyButton').addEventListener('click', () => {
                const copyButton = document.getElementById('copyButton');
                const filePathInput = document.getElementById('filePath');
                navigator.clipboard.writeText(filePathInput.value).then(() => {
                    showToast('Link copied to clipboard!');
                    copyButton.textContent = 'Copied!';
                    setTimeout(() => { copyButton.textContent = 'Copy'; }, 2000);
                }).catch(err => {
                    console.error('Failed to copy:', err);
                    showToast('Failed to copy link.');
                });
            });
        }
    }
});



