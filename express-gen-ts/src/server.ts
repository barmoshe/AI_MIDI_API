/**
 * Setup express server.
 */

import cookieParser from 'cookie-parser';
import morgan from 'morgan';
import path from 'path';
import helmet from 'helmet';
import express, { Request, Response, NextFunction } from 'express';
import logger from 'jet-logger';

import 'express-async-errors';

import BaseRouter from '@src/routes/api';
import Paths from '@src/constants/Paths';

import EnvVars from '@src/constants/EnvVars';
import HttpStatusCodes from '@src/constants/HttpStatusCodes';

import { NodeEnvs } from '@src/constants/misc';
import { RouteError } from '@src/other/classes';
import { Connection, Client } from '@temporalio/client';
import { nanoid } from 'nanoid';



// **** Variables **** //

const app = express();
const port = 8000;


// **** Setup **** //

// Basic middleware
app.use(express.json());
app.use(express.urlencoded({extended: true}));
app.use(cookieParser(EnvVars.CookieProps.Secret));

// Show routes called in console during development
if (EnvVars.NodeEnv === NodeEnvs.Dev.valueOf()) {
  app.use(morgan('dev'));
}

// Security
if (EnvVars.NodeEnv === NodeEnvs.Production.valueOf()) {
  app.use(helmet());
}

// Add APIs, must be after middleware
app.use(Paths.Base, BaseRouter);

// Add error handler
app.use((
  err: Error,
  _: Request,
  res: Response,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  next: NextFunction,
) => {
  if (EnvVars.NodeEnv !== NodeEnvs.Test.valueOf()) {
    logger.err(err, true);
  }
  let status = HttpStatusCodes.BAD_REQUEST;
  if (err instanceof RouteError) {
    status = err.status;
  }
  return res.status(status).json({ error: err.message });
});

// Handle POST requests to "/api/prompt"
app.post('/api/prompt', async (req: Request, res: Response) => {
  //log request body
  console.log(req.body);
  const { prompt, type} = req.body; // Extract the "prompt", "type", and "bpm" properties from the request body
  const bpm = req.body.bpm || 120; // Set the "bpm" variable to the value of the "bpm" property, or null if it is not present
  //init temporal workflow  with try catch
  

  const handle = await run(prompt , type);
  //if handle error return error message with s 400
  if (handle instanceof Error) {
    console.error('Error:', handle.message);
    res.status(400).send(handle.message);
  }
  else{
    var result = await handle.result();
    console.log(result);
   //send handle result to client
   res.send(result); 
  }

});

async function run( prompt: string, type: string) {
  // Connect to the default Server location
  const connection = await Connection.connect({ address: 'localhost:7233' });
  // In production, pass options to configure TLS and other settings:
  // {
  //   address: 'foo.bar.tmprl.cloud',
  //   tls: {}
  // }

  const client = new Client({
    connection,
    // namespace: 'foo.bar', // connects to 'default' namespace if not specified
  });

  const handle = await client.workflow.start("GenerateMIDIWorkflow", {
    taskQueue: 'GENERATE_MIDI_TASK_QUEUE',
    // type inference works! args: [name: string]
    args: [prompt, type],
    // in practice, use a meaningful business ID, like customerId or transactionId
    workflowId: 'workflow-' + nanoid(),
    retry: {
      maximumAttempts: 3,
    }
  });

  console.log(`Started workflow ${handle.workflowId}`);

 
  return  handle;
}




// **** Export default **** //

export default app;
