const user = {
  name: "Default User",
  type: "user",
};

const adminUser = Object.create(user);
adminUser.type = "admin";

function isAdmin(user) {
  return Object.getPrototypeOf(user) === adminUser;
}

// don't touch below this line

export { user, adminUser, isAdmin };

