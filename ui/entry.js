import { server as app } from "./index.js";

function shutdownGracefully() {
  app.server.close();
}

process.on("SIGINT", shutdownGracefully);
process.on("SIGTERM", shutdownGracefully);
