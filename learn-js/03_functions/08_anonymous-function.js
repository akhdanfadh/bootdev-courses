function printReports(intro, body, outro) {
  printCostReport(function (a) { return a.length * 2; }, intro);
  printCostReport(function (a) { return a.length * 3; }, body);
  printCostReport(function (a) { return a.length * 4; }, outro);
}

// don't touch below this line

function printCostReport(costCalculator, message) {
  const cost = costCalculator(message);
  console.log(`Message: "${message}" Cost: ${cost} cents`);
}

printReports(
  "Welcome to the Hotel California",
  "Such a lovely place",
  "Plenty of room at the Hotel California",
);

