const express = require("express");
const jose = require("jose");
const axios = require("axios");

// This sample shows how express.js can be used with axios and
// jose (JWT decoding library) to create a middleware validating
// the incoming requests for the users.
//
// This sample requires that you have these roles in your tenant:
// - projects.viewer: projects.read
// - projects.owner: projects.read | projects.write
//
// It also requires these bindings to be present:
// IDENTITY |      ROLES      | RESOURCE
// ---------------------------------------------
//   john   | projects.viewer | /projects/test
//   wade   | projects.owner  | /projects
//
// To understand how these work, please refer to the documentation at:
// https://docs.accessbox.dev/basics/designing_authorization

const app = express();

const tenant = process.env.ACCESSBOX_TENANT;
const apiKey = process.env.ACCESSBOX_API_KEY;

// This middleware validates the incoming requests for the users.
// It checks if the request has a valid JWT in the Authorization header.
// If the JWT is valid, it will set the claims in the request object.
app.use(async (req, res, next) => {
  const authHeader = req.headers.authorization;
  if (!authHeader) {
    return res.status(401).send("Missing Authorization header");
  }

  const token = authHeader.split(" ")[1];
  if (!token) {
    return res.status(401).send("Missing token");
  }

  try {
    // REPLACE THIS WITH VERIFY IN PRODUCTION
    const claims = jose.decodeJwt(token);

    req.claims = claims;

    next();
  } catch (err) {
    console.log(err);
    return res.status(401).send("Invalid token");
  }
});

// protect is a middleware that checks if the user has the required permission
// to access the resource. It sends a request to the AccessBox API to check
// if the user has the required permission.
function protect(permission) {
  return async (req, res, next) => {
    console.log(req.claims.sub, req.path, permission)

    const { data } = await axios.post(
      `https://api.accessbox.dev/v1/authorize?tenant=${tenant}`,
      {
        identity: req.claims.sub,
        resource: req.path,
        permission: permission,
      },
      {
        headers: {
          Authorization: `Bearer ${apiKey}`,
        },
      }
    );

    if (!data.allow) {
      return res.status(403).send("Forbidden");
    }

    next();
  };
}

// This route is protected by the middleware above.
// Only users with the projects.viewer role can access it.
app.get("/projects/:projectId?", protect('projects.read'), async (req, res) => {
  res.status(200).send("OK");
});

// This route is protected by the middleware above.
// Only users with the projects.owner role can access it.
app.post("/projects", protect('projects.write'), async (req, res) => {
  res.status(200).send("OK");
});

// export the app so it can be used in tests
module.exports = app;

// this makes the app testable in a test environment and only runs
// the server if the environment variable RUN_SERVER is set
if (process.env.RUN_SERVER) {
  app.listen(3000, () => {
    console.log("Server listening on port 3000");
  });
}
