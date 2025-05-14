import { BaseHandler } from './baseHandler';
import { sendMail } from '../services/emailService';
import { BaseMessage, WelcomePayload } from '../types';

class WelcomeEmailHandler extends BaseHandler {
  async handle({ payload }: BaseMessage): Promise<void> {
    const { email, name } = payload as WelcomePayload;
    
    await this.executeWithTracing(
      'welcome_email_handler',
      async () => {
        await sendMail({
          to: email,
          subject: 'Твои тренировки только начинаются — привет от FitnessThing 💪',
          templateName: 'welcome',
          context: { name: name, year: new Date().getFullYear() },
        });
      },
      {
        'user.email': email,
        'user.name': name,
        'template': 'welcome'
      }
    );
  }
}
export default new WelcomeEmailHandler();
