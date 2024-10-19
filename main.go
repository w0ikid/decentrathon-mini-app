package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "bytes"
)

var botToken string

// Проверяем наличие токена бота
func init() {
    botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
    if botToken == "" {
        log.Fatal("TELEGRAM_BOT_TOKEN is not set")
    }
}

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
                            "url": "https://decentrathon-mini-app.vercel.app/", // URL вашего приложения
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

    if resp.StatusCode != http.StatusOK {
        log.Printf("Telegram API returned non-200 status: %d", resp.StatusCode)
        return
    }

    log.Printf("Message sent to chat %d", chatID)
}

// Обработчик Webhook для получения сообщений от Telegram
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Webhook handler triggered")

    var update TelegramUpdate

    // Читаем тело запроса для отладки
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Printf("Error reading request body: %v", err)
        http.Error(w, "Cannot read body", http.StatusBadRequest)
        return
    }
    log.Printf("Received request body: %s", string(body))

    // Декодируем JSON от Telegram
    if err := json.Unmarshal(body, &update); err != nil {
        log.Printf("Error decoding update: %v", err)
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    // Логируем полученное сообщение
    log.Printf("Received message from chat %d: %s", update.Message.Chat.ID, update.Message.Text)

    // Возвращаем HTTP статус 200 OK
    w.WriteHeader(http.StatusOK)

    // Пример отправки ответа
    sendTelegramMessage(update.Message.Chat.ID, "Thanks for your message!")
}

// Обработчик для API, который возвращает курсы
func coursesHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(courses)
}

// Обработчик для статических файлов (HTML, CSS, JS)
func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
    // Проверяем существование файлов в папке static
    path := "./static" + r.URL.Path
    if _, err := os.Stat(path); os.IsNotExist(err) {
        http.NotFound(w, r)
        return
    }
    http.ServeFile(w, r, path)
}

func main() {
    // Проверка порта из переменной окружения
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Маршруты для Webhook, API курсов и статических файлов
    http.HandleFunc("/webhook", webhookHandler)
    http.HandleFunc("/api/courses", coursesHandler)
    http.HandleFunc("/", serveStaticFiles)

    // Запуск сервера
    log.Printf("Server is running on port %s...", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
