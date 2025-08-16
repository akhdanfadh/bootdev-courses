import { describe, assert, it, withSubmit } from "./unit_test.js";
import { getMonthlyPrice } from "./main.js";

describe("getMonthlyPrice", () => {
  const runCases = [
    { tier: "basic", expected: 10000 },
    { tier: "premium", expected: 15000 },
    { tier: "", expected: 0 },
  ];

  const submitCases = runCases.concat([
    { tier: "invalid", expected: 0 },
    { tier: "enterprise", expected: 50000 },
  ]);

  let testCases = runCases;
  if (withSubmit) {
    testCases = submitCases;
  }

  testCases.forEach(({ tier, expected }) => {
    it(`should return ${expected} for tier "${tier}"`, () => {
      assert.strictEqual(getMonthlyPrice(tier), expected);
    });
  });

  const numSkipped = submitCases.length - testCases.length;
  if (numSkipped > 0) {
    console.log(`- Skip: ${numSkipped} test case(s) for submit`);
  }
});

