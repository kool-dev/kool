const hostname = "0.0.0.0";
const port = 3000;

Bun.serve({
  hostname,
  port,
  fetch() {
    return new Response("Hello World");
  },
});

console.log("Server running at http://localhost:" + port + "/");
