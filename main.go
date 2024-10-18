package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "bytes"
)

var botToken = os.Getenv("TELEGRAM_BOT_TOKEN") // Токен бота читается из переменной окружения

// Структура для обновлений от Telegram
type TelegramUpdate struct {
    Message struct {
        Text string `json:"text"`
        Chat struct {
            ID int64 `json:"id"`
        } `json:"chat"`
    } `json:"message"`
}

// Структура для курсов
type Course struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
}

// Пример курсов
var courses = []Course{
    {ID: 1, Title: "Go Basics", Description: "Learn the basics of Go programming"},
    {ID: 2, Title: "Advanced JavaScript", Description: "Deep dive into JavaScript"},
}

// Функция для отправки сообщения пользователю с кнопкой для открытия Web App
func sendTelegramMessage(chatID int64, text string) {
    url := "https://api.telegram.org/bot" + botToken + "/sendMessage"

    // Формируем сообщение с кнопкой для запуска Web App
    message := map[string]interface{}{
        "chat_id": chatID,
        "text":    text,
        "reply_markup": map[string]interface{}{
            "inline_keyboard": [][]map[string]interface{}{
                {
                    {
                        "text": "Open LMS",
                        "web_app": map[string]string{
                            "url": "https://your-domain.com", // URL вашего приложения
                        },
                    },
                },
            },
        },
    }

    // Преобразуем сообщение в JSON
    jsonData, _ := json.Marshal(message)

    // Отправляем POST-запрос в API Telegram
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Println("Failed to send message:", err)
        return
    }
    defer resp.Body.Close()

    log.Println("Message sent to chat", chatID)
}

// Обработчик Webhook для получения сообщений от Telegram
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Webhook handler triggered") // Логирует начало работы обработчика

    var update TelegramUpdate

    // Попробуйте логировать запрос целиком
    log.Println("Received request:", r)

    // Декодируем JSON от Telegram
    if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
        log.Printf("Error decoding update: %v", err)
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    // Логируем полученное сообщение
    log.Printf("Received message from chat %d: %s", update.Message.Chat.ID, update.Message.Text)

    // Возвращаем HTTP статус 200 OK
    w.WriteHeader(http.StatusOK)
}


// Обработчик для API, который возвращает курсы
func coursesHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(courses)
}

// Обработчик для статических файлов (HTML, CSS, JS)
func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./static"+r.URL.Path)
}

func main() {
    // Маршруты для Webhook, API курсов и статических файлов
    http.HandleFunc("/webhook", webhookHandler)
    http.HandleFunc("/api/courses", coursesHandler)
    http.HandleFunc("/", serveStaticFiles)

    // Запуск сервера на порту 8080
    log.Println("Server is running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
