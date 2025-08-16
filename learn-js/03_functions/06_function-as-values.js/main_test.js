import { describe, assert, it, withSubmit } from "../unit_test.js";
import { reformat } from "./main.js";

describe("reformat", () => {
  const runCases = [
    {
      message: "hello",
      formatter: (msg) => msg.toUpperCase(),
      expected: "TEXTIO: HELLO",
    },
    {
      message: "world",
      formatter: (msg) => `*${msg}*`,
      expected: "TEXTIO: ***world***",
    },
  ];

  const submitCases = runCases.concat([
    {
      message: "test",
      formatter: (msg) => msg.split("").reverse().join(""),
      expected: "TEXTIO: tset",
    },
    {
      message: "repeat",
      formatter: (msg) => msg + "!",
      expected: "TEXTIO: repeat!!!",
    },
  ]);

  let testCases = runCases;
  if (withSubmit) {
    testCases = submitCases;
  }

  testCases.forEach(({ message, formatter, expected }) => {
    it(`should return "${expected}" for message "${message}" with given formatter`, () => {
      assert.strictEqual(reformat(message, formatter), expected);
    });
  });

  const numSkipped = submitCases.length - testCases.length;
  if (numSkipped > 0) {
    console.log(`- Skip: ${numSkipped} test case(s) for submit`);
  }
});

