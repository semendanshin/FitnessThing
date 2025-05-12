import kafka from '../config/kafka';
import { BaseMessage } from '../types';
import logger from '../utils/logger';
import WelcomeEmailHandler from '../handlers/welcomeEmail';

const handlerMap: Record<string, any> = {
  welcome: WelcomeEmailHandler,
};

export async function runConsumer(): Promise<void> {
  const consumer = kafka.consumer({ groupId: 'email-group' });
  await consumer.connect();
  await consumer.subscribe({ topic: process.env.KAFKA_TOPIC!, fromBeginning: false });

  await consumer.run({
    eachMessage: async ({ message }) => {
      try {
        const event: BaseMessage = JSON.parse(message.value!.toString());
        const handler = handlerMap[event.type];
        if (!handler) throw new Error(`Unknown type ${event.type}`);
        await handler.handle(event);
      } catch (err: any) {
        logger.error(`Error processing message: ${err}`);
      }
    },
  });
}
