const total = (function calculateTotal(numMessages, bytesPerMessage) {
  return numMessages * bytesPerMessage;
})(100, 24);

// don't touch below this line

console.log("Total message cost:", total);

