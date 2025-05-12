import { BaseMessage } from '../types';

export abstract class BaseHandler {
  abstract handle(event: BaseMessage): Promise<void>;
}
