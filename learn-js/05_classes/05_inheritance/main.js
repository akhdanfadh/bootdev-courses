class Sender {
  constructor(recipient) {
    this.recipient = recipient;
  }

  sendMessage(message) {
    throw new Error("sendMessage method must be implemented by subclasses");
  }
}

// don't touch above this line

class SMSSender extends Sender {
  sendMessage(message) {
    console.log(`Sending SMS to ${this.recipient}: ${message}`)
  }
}
class EmailSender extends Sender {
  sendMessage(message) {
    console.log(`Sending email to ${this.recipient}: ${message}`)
  }
}

// don't touch below this line

export { Sender, SMSSender, EmailSender };

