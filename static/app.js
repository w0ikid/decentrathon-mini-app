document.addEventListener("DOMContentLoaded", () => {
    if (window.Telegram && window.Telegram.WebApp) {
        console.log("Telegram WebApp initialized"); // Логирование инициализации Telegram Web App

        const user = window.Telegram?.WebApp?.initDataUnsafe?.user || {};
        console.log("User data:", user); // Логирование данных пользователя

        // Отображаем информацию о пользователе
        const userInfoContainer = document.getElementById('user-info');
        if (userInfoContainer) {
            userInfoContainer.innerHTML = `
                <h2>Welcome, ${user.first_name || 'Guest'}!</h2>
                <p>Username: ${user.username || 'N/A'}</p>
            `;
        } else {
            console.error("Element with id 'user-info' not found");
        }

        // Настройки Telegram Web App
        Telegram.WebApp.setBackgroundColor('#f0f0f0');
        Telegram.WebApp.setHeaderColor('#4CAF50');
    } else {
        console.error("Telegram WebApp not available");
    }

    const coursesContainer = document.getElementById('courses');
    const loadingElement = document.getElementById('loading');

    if (!coursesContainer || !loadingElement) {
        console.error("Courses or loading element not found");
        return;
    }

    // Функция для загрузки курсов с фильтром
    const loadCourses = (category = 'all') => {
        coursesContainer.innerHTML = ''; // Очищаем контейнер курсов
        loadingElement.style.display = 'block'; // Показываем спиннер

        fetch('/api/courses')
            .then(response => {
                console.log("API Response:", response); // Логирование ответа API
                return response.json();
            })
            .then(data => {
                loadingElement.style.display = 'none'; // Скрываем спиннер
                if (!data.length) {
                    coursesContainer.innerHTML = '<p>No courses available.</p>';
                } else {
                    data.forEach(course => {
                        if (category === 'all' || course.category === category) {
                            const courseElement = document.createElement('div');
                            courseElement.className = 'course-item';
                            courseElement.innerHTML = `
                                <div class="course-title">${course.title}</div>
                                <div class="course-description">${course.description}</div>
                            `;
                            coursesContainer.appendChild(courseElement);
                        }
                    });
                }
                Telegram.WebApp.showAlert('Courses have been successfully updated!');
            })
            .catch(error => {
                loadingElement.innerText = 'Error loading courses';
                Telegram.WebApp.showAlert('Failed to load courses!');
                console.error('Error fetching courses:', error);
            });
    };

    // Первоначальная загрузка курсов
    loadCourses();

    // Обработчик для фильтрации курсов
    const courseFilter = document.getElementById('course-filter');
    if (courseFilter) {
        courseFilter.addEventListener('change', () => {
            const selectedCategory = courseFilter.value;
            loadCourses(selectedCategory); // Загружаем курсы с фильтром
        });
    }

    // Обработчик для кнопки "Send a Message"
    const sendMessageButton = document.getElementById('send-message');
    if (sendMessageButton) {
        sendMessageButton.addEventListener('click', () => {
            Telegram.WebApp.sendData('Hello from Web App!'); // Отправляем сообщение обратно в Telegram
        });
    }

    // Обработчик для кнопки "Refresh Courses"
    const refreshCoursesButton = document.getElementById('refresh-courses');
    if (refreshCoursesButton) {
        refreshCoursesButton.addEventListener('click', () => {
            loadCourses(courseFilter ? courseFilter.value : 'all'); // Обновляем курсы с учетом выбранного фильтра
        });
    }

    // Обработчик для кнопки "Close App"
    const closeAppButton = document.getElementById('close-app');
    if (closeAppButton) {
        closeAppButton.addEventListener('click', () => {
            Telegram.WebApp.close(); // Закрываем Web App
        });
    }
    const aiResponse = document.getElementById('ai-response');
    const aiInput = document.getElementById('ai-input');
    const askAiButton = document.getElementById('ask-ai');

    if (askAiButton) {
        askAiButton.addEventListener('click', () => {
            const userInput = aiInput.value.toLowerCase();

            // Имитация ИИ-ответа
            if (userInput.includes('course') || userInput.includes('lesson')) {
                aiResponse.innerText = 'Sure! I recommend starting with the "Go Basics" course to improve your skills.';
            } else if (userInput.includes('progress')) {
                aiResponse.innerText = 'You have completed 50% of your current learning goals. Keep going!';
            } else if (userInput === '') {
                aiResponse.innerText = 'Please ask something!';
            } else {
                aiResponse.innerText = 'I am here to help with your learning journey! Ask me about your courses or progress.';
            }

            aiInput.value = ''; // Очищаем поле ввода
        });
    }
});
