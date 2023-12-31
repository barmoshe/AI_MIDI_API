"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const worker_1 = require("@temporalio/worker");
const activities_1 = require("./activities");
function run() {
    return __awaiter(this, void 0, void 0, function* () {
        // Step 1: Establish a connection with Temporal server.
        //
        // Worker code uses `@temporalio/worker.NativeConnection`.
        // (But in your application code it's `@temporalio/client.Connection`.)
        const connection = yield worker_1.NativeConnection.connect({
            address: 'localhost:7233',
            // TLS and gRPC metadata configuration goes here.
        });
        // Step 2: Register Workflows and Activities with the Worker.
        const worker = yield worker_1.Worker.create({
            connection,
            namespace: 'default',
            taskQueue: 'ts-worker',
            // Workflows are registered using a path as they run in a separate JS context.
            activities: { SendMIDITextRequest: activities_1.SendMIDITextRequest },
        });
        yield worker.run();
    });
}
run();
