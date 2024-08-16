import * as solanaWeb3 from "@solana/web3.js";

const MAINNET_BETA = solanaWeb3.clusterApiUrl("mainnet-beta");
const solanaConnection = new solanaWeb3.Connection(MAINNET_BETA);

interface Request {
  body: {
    transactionHash: string;
  };
}

interface Response {
  status: (code: number) => {
    json: (data: any) => void;
  };
}

const handler = async (req: Request, res: Response) => {
  const { transactionHash } = req.body;
  if (!transactionHash) {
    return res.status(401).json({
      error: "Invalid transaction hash",
    });
  }

  try {
    const transaction = await solanaConnection.getParsedTransaction(transactionHash);
    
    if (transaction) {
      return res.status(200).json(transaction);
    } else {
      return res.status(404).json({
        error: "Transaction not found",
      });
    }
  } catch (error) {
    console.log("Error:", error);
    return res.status(500).json({
      error: "Server error",
    });
  }
};

export default handler;