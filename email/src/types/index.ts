export interface BaseMessage {
  type: string;
  payload: any;
}

export interface WelcomePayload {
  email: string;
  name: string;
}

export interface InvoicePayload {
  email: string;
  name: string;
  invoiceId: string;
  amount: number;
}
