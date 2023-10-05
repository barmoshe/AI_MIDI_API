import './pre-start'; // Must be the first import
import logger from 'jet-logger';

import EnvVars from '@src/constants/EnvVars';
import app from './server';
import { Connection, Client } from '@temporalio/client';
import { nanoid } from 'nanoid';


// **** Run **** //

const SERVER_START_MSG = ('Express server started on port: ' + 
  EnvVars.Port.toString());

app.listen(EnvVars.Port, () => logger.info(SERVER_START_MSG));


