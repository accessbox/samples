const app = require('./main');
const request = require('supertest');

const johnToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3MTMyNjM5MDAsImV4cCI6MTc0NDc5OTkwMCwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiam9obiJ9._Mdu2Gvz6QsApNpACSZfwIJTOP1ZoKJADmXAHGqHJMc";
const wadeToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3MTMyNjM5MDAsImV4cCI6MTc0NDc5OTkwMCwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoid2FkZSJ9.EYx6AXb8dOLIYfDLhhrVEb0VJxeDM_RVTZp82tZ934w";

describe('john should be able to view /projects/test', () => {
  test("200 OK", async () => {

    const response = await request(app)
      .get("/projects/test")
      .set("Authorization", `Bearer ${johnToken}`);
    
    expect(response.statusCode).toBe(200);
  })
});