
## Description
Path of Exile II telegram whispering notifier for trade.

## How to use
1. [Download](https://github.com/seorgiy/poe-trade-notifier/releases/latest) last release, unzip somewhere and run poe-trade-notifier.exe, follow the instructions
2. Show a path to the logfile (for the steam version it's something like C:\Steam\steamapps\common\Path of Exile 2\logs\Client.txt)
3. Enter the telegram bot token. If you don't have one, use [@BotFather](https://t.me/botfather) to create a new bot and get its token, it's free and simple
4. Enter your telegram_id. Use [GetMyId](https://t.me/getmyid_bot) to obtain it
5. Send any message to your bot, just to start a conversation (bots cannot initiate a conversation due to Telegram rules). It only needs to be done once.

You have to be online in POE on the same PC in order to get new messages

## How it works
Application constantly observes log file with all chat messages. If there is Whisper message, it instantly forwards it to Telegram.
There is no any types of injections, so it's not violating any POE rules what so ever.

Windows may not like unknown .exe file, it's expected behavior, you still can run it by pressing «Run anyway» or something like that. Also any antivirus apps will not be happy i suppose.

## Known issues
Koreans messages not being logged in Client.txt so they also will not be processed by poe-trade-notifier. [more details](https://www.reddit.com/r/pathofexile/comments/c6yijd/bug_important_korean_whispers_not_being_logged_in/)

## Need help?
Please feel free to reach out to [me](https://t.me/seorgiy) or create an [issue](https://github.com/seorgiy/poe-trade-notifier/issues)
