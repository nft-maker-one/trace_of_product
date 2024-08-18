import { NextApiRequest, NextApiResponse } from 'next';
import * as solanaWeb3 from "@solana/web3.js";

const MAINNET_BETA = solanaWeb3.clusterApiUrl("mainnet-beta");
const solanaConnection = new solanaWeb3.Connection(MAINNET_BETA);

const getAddressInfo = async (address: string, numTx: number = 3) => {
  const pubKey = new solanaWeb3.PublicKey(address);
  const transactionList = await solanaConnection.getSignaturesForAddress(
    pubKey,
    { limit: numTx }
  );
  const accountBalance = await solanaConnection.getBalance(pubKey);

  return { transactionList, accountBalance };
};

const handler = async (req: NextApiRequest, res: NextApiResponse) => {
  const queryAddress = req.query?.address as string;
  if (!queryAddress) {
    return res.status(401).json({
      message: "无效的地址",
    });
  }
  try {
    const { accountBalance, transactionList } = await getAddressInfo(
      queryAddress
    );
    return res.status(200).json({ transactionList, accountBalance });
  } catch (error) {
    console.log(error);
    return res.status(500).json({
      message: "出了点问题，请稍后再试",
    });
  }
};

export default handler;