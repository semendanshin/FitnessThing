import nodemailer from 'nodemailer';
import emailConfig from '../config/email';
import { renderTemplate } from './templateService';

const transporter = nodemailer.createTransport(emailConfig);

interface SendMailOptions {
  to: string;
  subject: string;
  templateName: string;
  context: Record<string, any>;
}

export async function sendMail({ to, subject, templateName, context }: SendMailOptions): Promise<void> {
  const html = renderTemplate(templateName, context);
  await transporter.sendMail({
    from: emailConfig.from,
    to,
    subject,
    html,
  });
}
