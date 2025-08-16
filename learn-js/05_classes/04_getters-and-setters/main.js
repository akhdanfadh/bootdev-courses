class Contact {
  constructor(name, phoneNumber) {
    this.name = name;
    this._phoneNumber = phoneNumber;
  }

  set phoneNumber(newPhoneNumber) {
    if (typeof newPhoneNumber === 'string' && newPhoneNumber.length === 10) {
      this._phoneNumber = newPhoneNumber;
      return;
    }
    throw new Error("Invalid phone number.")
  }

  get phoneNumber() {
    let number = this._phoneNumber;
    return `(${number.slice(0, 3)}) ${number.slice(3, 6)}-${number.slice(6)}`;
  }

}

// don't touch below this line

export { Contact };

