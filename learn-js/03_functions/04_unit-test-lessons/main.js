function getMonthlyPrice(tier) {
  switch (tier) {
    case "basic":
      return 100 * 100;
    case "premium":
      return 100 * 150;
    case "enterprise":
      return 100 * 500;
    default:
      return 0;
  }
}

// don't touch below this line

export { getMonthlyPrice };

