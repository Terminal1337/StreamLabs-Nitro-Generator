import asyncio
from aiosmtpd.controller import Controller
import email
from fastapi import FastAPI, HTTPException
import uvicorn
from colorama import Fore, init
import re

init(convert=True)
app = FastAPI()

__emails__ = {}

class MyHandler:

    async def handle_DATA(self, server, session, envelope):
        try:
            email_message = email.message_from_bytes(envelope.content)
            subject = email_message["Subject"]
            sender = email_message["From"]

            if sender and sender.strip().lower() == "support@streamlabs.com":
                print(f"{Fore.RED}INFO: {Fore.GREEN}{Fore.RESET}{Fore.MAGENTA} [{envelope.rcpt_tos[0]}] ------------>  {Fore.RESET}{Fore.BLUE}[{subject}]{Fore.RESET}")
                __emails__[envelope.rcpt_tos[0]] = subject
                return '250 OK'
            else:
                print(f"{Fore.RED}INFO: {Fore.GREEN}Ignored email from unauthorized sender: {sender}{Fore.RESET}")
        except Exception as e:
            print(f"{Fore.RED}ERROR: {Fore.RESET}Error processing email: {str(e)}")

@app.get("/api/streamlabs.com")
async def get_email(email: str, type: int = 1):
    try:
        if type == 1:
            subject = __emails__.get(email, "Email not found")
            if subject != "Email not found":
                match = re.match(r"(\d{6,}) - Streamlabs ID verification code", subject)
                if match:
                    code = match.group(1)
                else:
                    code = "Invalid format"

                del __emails__[email]
                return {"success": True, "email": email, "subject": subject, "code": code}

            return {"success": False, "email": email, "subject": subject}
    except Exception as e:
        return {"success": False, "email": email, "subject": "", "error": str(e)}

if __name__ == "__main__":
    controller = Controller(MyHandler(), hostname='0.0.0.0', port=25)
    controller.start()

    try:
        print(f"{Fore.RED}INFO: {Fore.GREEN}{Fore.RESET}{Fore.MAGENTA} [SMTP] {Fore.RESET}{Fore.BLUE} Server Started on Port 25{Fore.RESET}")
        print(f"{Fore.RED}INFO: {Fore.GREEN}{Fore.RESET}{Fore.MAGENTA} [HTTP] {Fore.RESET}{Fore.BLUE} Server Started on Port 8000{Fore.RESET}")
        uvicorn.run(app, host="0.0.0.0", port=8000, log_level="error", access_log=False)
    except KeyboardInterrupt:
        controller.stop()
