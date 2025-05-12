import dotenv from 'dotenv';
dotenv.config();

import { runConsumer } from './consumers/kafkaConsumer';
import logger from './utils/logger';

runConsumer()
  .then(() => logger.info('Kafka consumer started'))
  .catch(err => {
    logger.error('Startup error', err);
    process.exit(1);
  });
