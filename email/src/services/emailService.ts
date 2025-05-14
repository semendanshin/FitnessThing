import nodemailer from 'nodemailer';
import emailConfig from '../config/email';
import { renderTemplate } from './templateService';
import { context, SpanStatusCode, trace } from '@opentelemetry/api';
import logger from '../utils/logger';

const transporter = nodemailer.createTransport(emailConfig);

interface SendMailOptions {
  to: string;
  subject: string;
  templateName: string;
  context: Record<string, any>;
}

export async function sendMail({ to, subject, templateName, context }: SendMailOptions): Promise<void> {
  const tracer = trace.getTracer('email-service');
  
  // Создаем спан для рендеринга шаблона
  const renderSpan = tracer.startSpan('render_email_template');
  renderSpan.setAttribute('template.name', templateName);
  
  let html: string;
  try {
    html = renderTemplate(templateName, context);
    renderSpan.setStatus({ code: SpanStatusCode.OK });
  } catch (error: any) {
    renderSpan.recordException(error);
    renderSpan.setStatus({
      code: SpanStatusCode.ERROR,
      message: error.message
    });
    throw error;
  } finally {
    renderSpan.end();
  }

  // Создаем спан для отправки письма
  const sendSpan = tracer.startSpan('send_email');
  sendSpan.setAttribute('email.to', to);
  sendSpan.setAttribute('email.subject', subject);
  
  try {
    const result = await transporter.sendMail({
      from: emailConfig.from,
      to,
      subject,
      html,
    });
    
    sendSpan.setAttribute('email.messageId', result.messageId || '');
    sendSpan.setStatus({ code: SpanStatusCode.OK });

    logger.info(`Email sent to ${to}, messageId: ${result.messageId}`);
  } catch (error: any) {
    sendSpan.recordException(error);
    sendSpan.setStatus({
      code: SpanStatusCode.ERROR,
      message: error.message
    });
    
    logger.error(`Failed to send email to ${to}: ${error.message}`);
    throw error;
  } finally {
    sendSpan.end();
  }
}
