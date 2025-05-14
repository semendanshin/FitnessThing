import kafka from '../config/kafka';
import { BaseMessage } from '../types';
import logger from '../utils/logger';
import WelcomeEmailHandler from '../handlers/welcomeEmail';

const handlerMap: Record<string, any> = {
  welcome: WelcomeEmailHandler,
};

import { context, propagation, SpanStatusCode, trace } from '@opentelemetry/api';
import { IHeaders } from 'kafkajs';

function extractContextFromKafkaHeaders(headers: IHeaders) {
  const carrier: Record<string, string> = {};
  
  for (const key in headers) {
    if (headers[key]) {
      const value = headers[key]!.toString();
  
      carrier[key.toLowerCase()] = value;
    }
  }
  
  return propagation.extract(context.active(), carrier);
}

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
        
        const ctx = extractContextFromKafkaHeaders(message.headers || {});
        
        const span = trace.getTracer('email-service').startSpan('handle_message', {}, ctx);
        
        span.setAttribute('message.type', event.type);
        
        try {
          await context.with(trace.setSpan(ctx, span), async () => {
            await handler.handle(event);
          });
        } catch (err: any) {
          span.recordException(err);
          span.setStatus({
            code: SpanStatusCode.ERROR,
            message: err.message
          });
          throw err;
        } finally {
          span.end();
        }

      } catch (err: any) {
        logger.error(`Error processing message: ${err}`);
      }
    },
  });
}
