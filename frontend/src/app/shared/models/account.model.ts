export type AccountStatus = 'active' | 'blocked' | 'closed';

export interface Account {
  id: string;
  number: string;
  customer_id: string;
  balance_cents: number;
  status: AccountStatus;
}
