import { Kafka } from 'kafkajs';

const kafka = new Kafka({
  clientId: process.env.KAFKA_CLIENT_ID || 'email-service',
  brokers: (process.env.KAFKA_BROKERS || 'localhost:9092').split(','),
});

export default kafka;
