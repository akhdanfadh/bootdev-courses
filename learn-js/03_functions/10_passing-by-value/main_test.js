import { describe, assert, it, withSubmit } from "./unit_test.js";
import { monthlyBillIncrease } from "./main.js";

describe("monthlyBillIncrease", () => {
  const runCases = [
    { costPerSend: 0.05, numLastMonth: 1000, numThisMonth: 1500, expected: 25 },
    { costPerSend: 0.1, numLastMonth: 500, numThisMonth: 1000, expected: 50 },
    { costPerSend: 0.05, numLastMonth: 0, numThisMonth: 1000, expected: 50 },
  ];

  const submitCases = runCases.concat([
    { costPerSend: 0.25, numLastMonth: 200, numThisMonth: 400, expected: 50 },
    { costPerSend: 0.02, numLastMonth: 1000, numThisMonth: 1500, expected: 10 },
  ]);

  let testCases = runCases;
  if (withSubmit) {
    testCases = submitCases;
  }

  testCases.forEach(({ costPerSend, numLastMonth, numThisMonth, expected }) => {
    it(`should return ${expected} for costPerSend ${costPerSend}, numLastMonth ${numLastMonth}, and numThisMonth ${numThisMonth}`, () => {
      assert.strictEqual(
        monthlyBillIncrease(costPerSend, numLastMonth, numThisMonth),
        expected,
      );
    });
  });

  const numSkipped = submitCases.length - testCases.length;
  if (numSkipped > 0) {
    console.log(`- Skip: ${numSkipped} test case(s) for submit`);
  }
});

