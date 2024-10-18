document.addEventListener("DOMContentLoaded", () => {
    if (window.Telegram.WebApp) {
        const user = Telegram.WebApp.initDataUnsafe?.user || {};

        // Отображаем информацию о пользователе
        const userInfoContainer = document.getElementById('user-info');
        userInfoContainer.innerHTML = `
            <h2>Hello, ${user.first_name || 'Guest'}!</h2>
            <p>Username: ${user.username || 'N/A'}</p>
        `;

        // Настройки Telegram Web App
        Telegram.WebApp.setBackgroundColor('#f0f0f0');
        Telegram.WebApp.setHeaderColor('#4CAF50');
    }

    // Загружаем данные о курсах через API
    fetch('/api/courses')
        .then(response => response.json())
        .then(data => {
            const coursesContainer = document.getElementById('courses');
            data.forEach(course => {
                const courseElement = document.createElement('div');
                courseElement.className = 'course-item';
                courseElement.innerHTML = `
                    <div class="course-title">${course.title}</div>
                    <div class="course-description">${course.description}</div>
                `;
                coursesContainer.appendChild(courseElement);
            });
        })
        .catch(error => console.error('Error fetching courses:', error));

    // Обработчик для кнопки "Send a Message"
    const sendMessageButton = document.getElementById('send-message');
    sendMessageButton.addEventListener('click', () => {
        Telegram.WebApp.sendData('Hello from Web App!'); // Отправляем сообщение обратно в Telegram
    });

    // Обработчик для кнопки "Close App"
    const closeAppButton = document.getElementById('close-app');
    closeAppButton.addEventListener('click', () => {
        Telegram.WebApp.close(); // Закрываем Web App
    });
});
