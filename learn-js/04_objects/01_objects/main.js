function createMessage(phoneNumber, message) {
  return {
    phoneNumber: phoneNumber,
    message: message,
    messageLength: message.length
  };
}

export { createMessage };

