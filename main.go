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

// Функция для отправки простого текстового сообщения
func sendTelegramMessage(chatID int64, text string) {
    url := "https://api.telegram.org/bot" + botToken + "/sendMessage"
    message := map[string]interface{}{
        "chat_id": chatID,
        "text":    text,
    }
    jsonData, _ := json.Marshal(message)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Println("Failed to send message:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode != http.StatusOK {
        log.Printf("Telegram API returned non-200 status: %d, response: %s", resp.StatusCode, string(body))
        return
    }

    log.Printf("Message sent to chat %d", chatID)
}

// Функция для отправки инлайн-кнопки для открытия Web App
func sendWebAppButton(chatID int64) {
    url := "https://api.telegram.org/bot" + botToken + "/sendMessage"

    // Формируем сообщение с инлайн-кнопкой для запуска Web App
    message := map[string]interface{}{
        "chat_id": chatID,
        "text":    "Нажмите на кнопку ниже, чтобы открыть Web App:",
        "reply_markup": map[string]interface{}{
            "inline_keyboard": [][]map[string]interface{}{
                {
                    {
                        "text": "Открыть Web App",
                        "web_app": map[string]string{
                            "url": "https://ae08-128-127-102-7.ngrok-free.app", // URL вашего приложения
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
        log.Println("Failed to send Web App button:", err)
        return
    }
    defer resp.Body.Close()

    // Читаем и выводим ответ от Telegram API для отладки
    body, _ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode != http.StatusOK {
        log.Printf("Telegram API returned non-200 status: %d, response: %s", resp.StatusCode, string(body))
        return
    }

    log.Printf("Web App button sent to chat %d", chatID)
}

// Обработчик для вебхука Telegram
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    var update TelegramUpdate

    // Декодируем JSON от Telegram
    if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
        log.Printf("Error decoding update: %v", err)
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    chatID := update.Message.Chat.ID
    messageText := update.Message.Text

    // Если пользователь отправил команду /helloworld, отправляем "Hello, World!"
    if messageText == "/helloworld" {
        sendTelegramMessage(chatID, "Hello, World!")
    } else if messageText == "/webapp" {
        // Если команда /webapp, отправляем кнопку для открытия Web App
        sendWebAppButton(chatID)
    } else {
        sendTelegramMessage(chatID, "Unknown command.")
    }

    // Возвращаем HTTP статус 200 OK
    w.WriteHeader(http.StatusOK)
}

func main() {
    http.HandleFunc("/webhook", webhookHandler)
    
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server is running on port %s...", port)
    log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
