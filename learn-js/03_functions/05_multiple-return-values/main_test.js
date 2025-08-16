import { describe, assert, it, withSubmit } from "../unit_test.js";
import { isClean } from "./main.js";

describe("isClean", () => {
  const runCases = [
    { review: "avril lavigne has best dang tour", expected: false },
    { review: "what a bad film", expected: true },
    { review: "oh my heck I hated it", expected: false },
  ];

  const submitCases = runCases.concat([
    { review: "ripoff", expected: true },
    { review: "That was a pleasure", expected: true },
    { review: "shoot! I cant say I liked it", expected: false },
  ]);

  let testCases = runCases;
  if (withSubmit) {
    testCases = submitCases;
  }

  testCases.forEach(({ review, expected }) => {
    it(`should return ${expected} for review "${review}"`, () => {
      assert.strictEqual(isClean(review), expected);
    });
  });

  const numSkipped = submitCases.length - testCases.length;
  if (numSkipped > 0) {
    console.log(`- Skip: ${numSkipped} test case(s) for submit`);
  }
});

