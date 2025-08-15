// undefined: It doesn't exist at all. In grug-speak undefined is "very nothing"
// null: It (kind of) exists, but it's empty. In grug-speak null is "kinda nothing"

let sentMessages = null; // explicit
let deliveredMessages = null;
let failedMessages = null;

// don't touch below this line

console.log("Sent is null:", sentMessages === null);
console.log("Delivered is null:", deliveredMessages === null);
console.log("Failed is null:", failedMessages === null);

