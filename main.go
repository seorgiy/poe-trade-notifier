package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/sqweek/dialog"

	read_file_last_line "github.com/abusizhishen/read-file-last-line"
	"github.com/spf13/viper"
)

var dateRegexp = regexp.MustCompile(`\d{4}\/\d{1,2}\/\d{1,2}\s\d{2}:\d{2}:\d{2}`)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
}

func main() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.WriteConfigAs("config.json")
			viper.ReadInConfig()
		} else {
			dialog.Message("%s", err).Title("Something is wrong").Error()
			os.Exit(0)
		}
	}

	locale := getLocale()
	fmt.Printf("Locale is %s\n", locale)
	whisperStartings := map[string]string{
		"ru":    "От кого",
		"en":    "From",
		"de":    "Von",
		"fr":    "De",
		"pt-Br": "De",
		"es":    "De",
	}

	if whisperStartings[locale] == "" {
		dialog.Message("%s", "Your locale is not supported yet, please create an issue here https://github.com/seorgiy/poe-trade-notifier").Title("Something is wrong").Error()
		os.Exit(0)
	}
	regexString := strings.Replace("] (@From.+?:.*)", "From", whisperStartings[locale], 1)
	whisperRegexp := regexp.MustCompile(regexString)

	logFilePath := getLogFile()
	fmt.Printf("Logs path is %s\n", logFilePath)

	botToken := getBotToken()
	fmt.Printf("BotToken is %s\n", botToken)

	chatId := getUserTelegramId()
	fmt.Printf("TelegramID is %d\n", chatId)

	bot := getBot(botToken)
	isChatAvailable(bot, chatId)

	byt, err := read_file_last_line.ReadLastLine(logFilePath)
	if err != nil {
		dialog.Message("%s", err).Title("Can't find the logs").Error()
		os.Exit(0)
	}

	lastScannedTime, err := getTime(*dateRegexp, string(byt))
	if err != nil {
		fmt.Println(err)
	}

	//observe logs loop
	fmt.Println("Listening...")
	for {
		logFile, err := os.Open(logFilePath)
		if err != nil {
			dialog.Message("%s", err).Title("Can't find the logs").Error()
			os.Exit(0)
		}

		scanner := bufio.NewScanner(logFile)
		line := 0
		for scanner.Scan() {
			res1 := whisperRegexp.Match(scanner.Bytes())

			if res1 {
				date, err := getTime(*dateRegexp, scanner.Text())
				if err != nil {
					dialog.Message("%s", err).Title("Cant parse time from logs").Error()
					os.Exit(0)
				}

				if lastScannedTime.Unix() < date.Unix() {
					text := whisperRegexp.FindStringSubmatch(scanner.Text())[1]
					fmt.Println(text)
					bot.SendMessage(chatId, text, &gotgbot.SendMessageOpts{})
					lastScannedTime = date
				}
			}
			line++
		}
		err = scanner.Err()
		if err != nil {
			dialog.Message("%s", err).Title("Can't read the logs").Error()
			os.Exit(0)
		}

		time.Sleep(2 * time.Second)
	}
}

func getLogFile() string {
	logFilePath := viper.GetString("client_logs_path")
	if logFilePath == "" {
		dialog.Message("%s", "Select a Client.txt log file in the POE logs directory").Title("Greetings, exile!").Info()
		newFilePath, err := dialog.File().Filter("Client", "txt").Load()
		if err != nil {
			log.Fatalf("Can't find the logs :( %s", err)
		}
		logFilePath = newFilePath
		viper.Set("client_logs_path", newFilePath)
		viper.WriteConfig()
	}
	return logFilePath
}

func getBotToken() string {
	botToken := viper.GetString("bot_token")
	if botToken == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the bot token. If you don't have one, use @BotFather to create a new bot and get its token, it's free:")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		viper.Set("bot_token", text)
		botToken = text
		viper.WriteConfig()
	}

	return botToken
}

func getUserTelegramId() int64 {
	id := viper.GetInt64("chat_id")
	if id == 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your telegram_id. Use bot @getmyid_bot to obtain it:")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		int_id, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			dialog.Message("%s", err).Title("Can't form chat id").Error()
			os.Exit(0)
		}
		id = int_id
		viper.Set("chat_id", int_id)
		viper.WriteConfig()
	}

	return id
}

func getLocale() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		dialog.Message("%s", err).Title("Can't find the config file").Error()
		os.Exit(0)
	}
	poeConfigDir := filepath.Join(userHomeDir, "Documents", "My Games", "Path of Exile 2", "poe2_production_Config.ini")

	poeConfigFile, err := os.Open(poeConfigDir)
	if err != nil {
		dialog.Message("%s", err).Title("Can't find the config file").Error()
		os.Exit(0)
	}

	scanner := bufio.NewScanner(poeConfigFile)
	locale := ""
	line := 0
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "language=") {
			locale = strings.Replace(scanner.Text(), "language=", "", 1)
			return locale
		}
		line++
	}

	if locale == "" {
		dialog.Message("%s", err).Title("Can't find the config file").Error()
		os.Exit(0)
	}
	return locale
}

func getTime(regexp regexp.Regexp, rawString string) (time.Time, error) {
	layout := "2006/01/02 15:04:05"
	date := regexp.FindStringSubmatch(rawString)[0]
	return time.Parse(layout, date)
}

func getBot(botToken string) *gotgbot.Bot {
	b, err := gotgbot.NewBot(botToken, nil)
	if err != nil {
		dialog.Message("%s", err).Title("Can't initialize telegram bot").Error()
		os.Exit(0)
	}
	return b
}

func isChatAvailable(bot *gotgbot.Bot, chatId int64) {
	_, errChat := bot.GetChat(chatId, &gotgbot.GetChatOpts{})
	if errChat != nil {
		dialog.Message("%s", errChat).Title("The bot can't text you. You need to send him «/start» for the first time").Error()
		os.Exit(0)
	}
}
