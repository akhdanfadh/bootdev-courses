function createContact(phoneNumber, name = "Anonymous", avatar = "default.jpg") {
  if (!phoneNumber) {
    return "Invalid phone number";
  }
  return `Contact saved! Name: ${name}, Phone number: ${phoneNumber}, Avatar: /public/pictures/${avatar}`;

}

// don't touch below this line

export { createContact };


