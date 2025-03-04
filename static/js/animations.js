document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('button').forEach(button => {
        button.addEventListener('mouseover', () => {
            button.classList.add('transition', 'duration-300');
        });
    });
});