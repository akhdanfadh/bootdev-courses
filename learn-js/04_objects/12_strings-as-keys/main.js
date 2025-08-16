function getProviderCount(provider, counts) {
  if (!counts[provider]) {
    return 0;
  }
  return counts[provider];
}

// don't touch below this line

export { getProviderCount };

