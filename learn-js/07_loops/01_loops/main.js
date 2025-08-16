function bulkSendCost(numMessages) {
  let totalCost = 0;
  for (let i = 0; i < numMessages; i++) {
    totalCost += 1.0 + 0.01 * i;
  };
  return totalCost;
}

export { bulkSendCost };

