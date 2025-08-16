import { describe, assert, it, withSubmit } from "../unit_test.js";
import { createContact } from "./main.js";

describe("createContact", () => {
  const runCases = [
    {
      phoneNumber: "555-1234",
      name: "Alice",
      expected:
        "Contact saved! Name: Alice, Phone number: 555-1234, Avatar: /public/pictures/default.jpg",
    },
    {
      phoneNumber: "555-5678",
      avatar: "bob.jpg",
      expected:
        "Contact saved! Name: Anonymous, Phone number: 555-5678, Avatar: /public/pictures/bob.jpg",
    },
    {
      phoneNumber: "",
      name: "John",
      avatar: "john.jpg",
      expected: "Invalid phone number",
    },
  ];

  const submitCases = runCases.concat([
    {
      phoneNumber: "555-1234",
      name: "Lane",
      avatar: "lane.png",
      expected:
        "Contact saved! Name: Lane, Phone number: 555-1234, Avatar: /public/pictures/lane.png",
    },
    {
      phoneNumber: "555-9876",
      name: "Jesse",
      expected:
        "Contact saved! Name: Jesse, Phone number: 555-9876, Avatar: /public/pictures/default.jpg",
    },
  ]);

  let testCases = runCases;
  if (withSubmit) {
    testCases = submitCases;
  }

  testCases.forEach(({ phoneNumber, name, avatar, expected }) => {
    it(`should return "${expected}" for phone number "${phoneNumber}", name "${name}", and avatar "${avatar}"`, () => {
      assert.strictEqual(createContact(phoneNumber, name, avatar), expected);
    });
  });

  const numSkipped = submitCases.length - testCases.length;
  if (numSkipped > 0) {
    console.log(`- Skip: ${numSkipped} test case(s) for submit`);
  }
});

