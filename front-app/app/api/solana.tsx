import { Connection, clusterApiUrl } from '@solana/web3.js';

const connection = new Connection(clusterApiUrl('mainnet-beta'), 'confirmed');

export default async function handler(req: any, res: { status: (arg0: number) => { (): any; new(): any; json: { (arg0: { blockHeight?: number; transactionCount?: number; error?: any; }): void; new(): any; }; }; }) {
  try {
    const blockHeight = await connection.getSlot();
    const transactionCount = await connection.getTransactionCount();
    res.status(200).json({ blockHeight, transactionCount });
  } catch (error) {
    // Ensure that error is an instance of Error
    const errorMessage = error instanceof Error ? error.message : 'Unknown error';
    res.status(500).json({ error: errorMessage });
  }
}