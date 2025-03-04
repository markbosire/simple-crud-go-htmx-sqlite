// Optional: Keep this if you want a fallback or custom function
function showSnackbar(message, type = 'success') {
    Toastify({
        text: message,
        duration: 3000,
        gravity: "bottom",
        position: "right",
        backgroundColor: type === 'success' 
            ? 'linear-gradient(to right, #00b09b, #96c93d)' 
            : 'linear-gradient(to right, #ff5f6d, #ffc371)',
        stopOnFocus: true
    }).showToast();
}

// Expose showSnackbar globally (optional)
window.showSnackbar = showSnackbar;