import { describe, assert, it, withSubmit } from "./unit_test.js";
import { getMessageStatus } from "./main.js";

describe("getMessageStatus", () => {
  const runCases = [
    { message: "Tread lightly.", expected: "valid" },
    { message: "Jesse, we need to cook!", expected: "valid" },
    { message: "", expected: "invalid" },
  ];

  const submitCases = runCases.concat([
    { message: "I am the danger", expected: "valid" },
    { message: " Say my name  ", expected: "valid" },
  ]);

  let testCases = runCases;
  if (withSubmit) {
    testCases = submitCases;
  }

  testCases.forEach(({ message, expected }) => {
    it(`should return "${expected}" for message "${message}"`, () => {
      assert.strictEqual(getMessageStatus(message), expected);
    });
  });

  const numSkipped = submitCases.length - testCases.length;
  if (numSkipped > 0) {
    console.log(`- Skip: ${numSkipped} test case(s) for submit`);
  }
});

