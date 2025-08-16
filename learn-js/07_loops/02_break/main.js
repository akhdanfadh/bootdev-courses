function maxMessagesWithinBudget(budget) {
  let totalCost = 0;
  let i = 0
  for (i = 0; i >= 0; i++) {
    let cost = 1.0 + 0.01 * i;
    if (totalCost + cost > budget) {
      return i;
    }
    totalCost += cost;
  }
  return i;
}

export { maxMessagesWithinBudget };

