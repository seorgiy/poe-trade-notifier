
POE/POE2 telegram whispering notificator

## How to use
1. [Download](https://github.com/seorgiy/poe-trade-notifier/releases) last release, unzip somewhere and run poe-trade-notifier.exe, follow the instructions
2. Show a path to the logfile (for the steam version it's something like C:\Steam\steamapps\common\Path of Exile 2\logs\Client.txt)
3. Enter the telegram bot token. If you don't have one, use [@BotFather](https://t.me/botfather) to create a new bot and get its token, it's free and simple
4. Enter your telegram_id. Use [GetMyId](https://t.me/getmyid_bot) to obtain it
5. You have to be online in POE on the same PC in order to get new messages

## How it works
Application constantly observes log file with all chat messages. If there is Whisper message, it instantly forwards it to Telegram.
There is no any types of injections, so it's not violating any POE rules what so ever