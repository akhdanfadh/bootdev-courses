function reformat(message, formatter) {
  message = formatter(message);
  message = formatter(message);
  message = formatter(message);
  return "TEXTIO: " + message;
}

// don't touch below this line

export { reformat };

