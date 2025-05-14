// filepath: /Users/semyondanshin/Programming/FitnessTrainer/email/src/handlers/baseHandler.ts
import { BaseMessage } from '../types';
import { context, Span, SpanStatusCode, trace } from '@opentelemetry/api';
import logger from '../utils/logger';

export abstract class BaseHandler {
  abstract handle(event: BaseMessage): Promise<void>;
  
  /**
   * Выполняет функцию в контексте трейсинга с указанным именем операции
   * @param operationName Имя операции для спана
   * @param fn Функция для выполнения в контексте спана
   * @param attributes Дополнительные атрибуты для спана
   * @returns Результат выполнения функции
   */
  protected async executeWithTracing<T>(
    operationName: string,
    fn: (span: Span) => Promise<T>,
    attributes: Record<string, string | number | boolean> = {}
  ): Promise<T> {
    const tracer = trace.getTracer('email-service');
    const span = tracer.startSpan(operationName);
    
    // Добавляем все переданные атрибуты к спану
    Object.entries(attributes).forEach(([key, value]) => {
      span.setAttribute(key, value);
    });
    
    try {
      // Выполняем функцию в контексте спана
      const currentContext = context.active();
      const newContext = trace.setSpan(currentContext, span);
      return await context.with(newContext, async () => {
        return await fn(span);
      });
    } catch (error: any) {
      span.recordException(error);
      span.setStatus({
        code: SpanStatusCode.ERROR,
        message: error.message
      });
      logger.error(`Error in ${operationName}: ${error.message}`);
      throw error;
    } finally {
      span.end();
    }
  }
}
