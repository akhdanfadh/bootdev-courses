class Message {
  static numMessages = 0;
  static totalMessagesLength = 0;

  constructor(recipient, sender, body) {
    this.recipient = recipient;
    this.sender = sender;
    this.body = body;
    Message.numMessages++;
    Message.totalMessagesLength += body.length;
  }

  static getAverageMessageLength() {
    let average = Message.totalMessagesLength / Message.numMessages;
    return Math.round(average * 100) / 100; // Round to 2 decimal places }
  }
}

// don't touch below this line

export { Message };

