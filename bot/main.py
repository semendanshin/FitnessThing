from telegram import Application

app = Application

@app.command("start")
async def start(update: Update):
    await update.message.reply_text("Hello, world!")

async def main():
    await app.run()

if __name__ == "__main__":
    import asyncio
    asyncio.run(main())
