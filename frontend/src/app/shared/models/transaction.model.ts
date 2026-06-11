export enum TransactionType {
  Debit = 0,
  Credit = 1,
  Pix = 2,
}

export interface Transaction {
  type: TransactionType;
  amount_cents: number;
}
