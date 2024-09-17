# Goods Monitoring

Goods Monitoring is a tool that checks the availability of specific products in various online venues, based on a predefined configuration. If there is a change in the availability status of a product, the app will send a notification to a configured Telegram chat.

## Key Features

- **Market Monitoring**: Continuously monitors the availability of products listed in the `config.json`.
- **Telegram Notifications**: Alerts you via Telegram whenever a product's availability status changes.
- **Easy Configuration**: Add or modify products and venues in the `config.json` file to customize what the app monitors.

## Why Goods Monitoring is Useful

- **Automated Monitoring**: Stop manually checking multiple marketplaces for your favorite products. Goods Monitoring does the job for you.
- **Real-Time Updates**: Be instantly notified when an out-of-stock item becomes available or when a product goes out of stock.
- **Convenience**: Easily monitor multiple venues and products from one place.

## Installation

1. Clone the repository.
2. Ensure you have the following files in the same directory:
   - `goodinitoring` (the binary executable)
   - `config.json` (the configuration file for venues and products)
   - `.env` (containing your Telegram bot token and chat ID)
   
   Example `.env` file:
   ```env
   TELEGRAM_BOT_TOKEN=your_telegram_bot_token
   TELEGRAM_CHAT_ID=your_telegram_chat_id
   ```

   Example `config.json` file:
   ```json
   [
     {
       "venue": "Bubble Tea",
       "endpoint": "https://consumer-api.wolt.com/consumer-api/consumer-assortment/v1/venues/slug/mao-bubble-tea/assortment",
       "names": ["buldak (pink)"]
     },
     {
       "venue": "Best Friends",
       "endpoint": "https://consumer-api.wolt.com/consumer-api/consumer-assortment/v1/venues/slug/best-friends1/assortment/categories/slug/sauces-and-pates-for-cat-2?language=en",
       "names": ["Wellfed (Pate for Sterelized cat, with beef & salmon)", "Wellfed (Pate for sterilized cat, with chicken & Turkey)"]
     }
   ]
   ```

## How to Run

1. Open a terminal and navigate to the folder containing the binary file.
2. Run the Goods Monitoring app:
   ```bash
   ./goodinitoring
   ```
3. The app will start checking the availability of the products listed in `config.json` every 30 minutes. If a product's status changes, you will receive a notification on Telegram.

## Configuration

### Products and Venues

You can modify `config.json` to monitor different products and venues. Each entry requires:
- `venue`: The name of the venue (e.g., "Bubble Tea").
- `endpoint`: The API endpoint that provides product data.
- `names`: A list of product names you want to monitor for availability.

### Telegram Setup

Ensure that your `.env` file contains the correct `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID`. You can create a Telegram bot using [BotFather](https://core.telegram.org/bots#botfather) and get your chat ID by messaging the bot.

## Example Notification

When the status of a product changes, you will receive a message like this on Telegram:

```
**Product**: Wellfed (Pate for Sterelized cat, with beef & salmon)
**Venue**: Best Friends
**Status**: The product is available!
```

## Logs and Errors

The app logs important events, such as errors and status changes, to the console. You can also view the results of the last check in the `store.json` file.