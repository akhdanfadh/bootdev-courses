// Simple unit testing framework
export const withSubmit = process.argv.includes('-s');

export function describe(name, testFunction) {
  console.log(`\n--- ${name} ---`);
  testFunction();
}

export function it(description, testFunction) {
  try {
    testFunction();
    console.log(`✓ ${description}`);
  } catch (error) {
    console.log(`✗ ${description}`);
    console.log(`  Error: ${error.message}`);
  }
}

export const assert = {
  strictEqual(actual, expected) {
    if (actual !== expected) {
      throw new Error(`Expected ${expected}, but got ${actual}`);
    }
  },
  
  throws(fn, expectedMessage) {
    let threwError = false;
    let actualError = null;
    
    try {
      fn();
    } catch (error) {
      threwError = true;
      actualError = error;
    }
    
    if (!threwError) {
      throw new Error('Expected function to throw an error, but it did not');
    }
    
    if (expectedMessage && actualError.message !== expectedMessage) {
      throw new Error(`Expected error message "${expectedMessage}", but got "${actualError.message}"`);
    }
  }
};
