import React from "react";
import { fromUnixTime, format, formatDistanceToNow } from "date-fns";
import Link from "next/link";

interface Transaction {
  signature: string;
  slot: number;
  blockTime: number;
  confirmationStatus: "finalized" | "processed" | "confirmed";
}

interface TransactionListProps {
  transactionList: Transaction[];
  balance?: number;
}

const TransactionList: React.FC<TransactionListProps> = ({ transactionList, balance }) => {
  return (
    <div className="first-line:overflow-hidden transition-all duration-300 max-h-fit w-full h-full">
      {balance && (
        <h2 className="flex justify-between text-lg mb-4">
          Balance: <span>◎{balance}</span>
        </h2>
      )}
      {transactionList?.length > 0 ? (
        <div className="overflow-x-auto">
          <table className="w-full border-spacing-x-4 -ml-4 border-separate">
            <thead className="text-left">
              <tr>
                <th className="font-medium">Signature</th>
                <th className="font-medium">Block</th>
                <th className="font-medium">Age</th>
                <th className="font-medium">Status</th>
              </tr>
            </thead>
            <tbody>
              {transactionList.map((transaction) => (
                <tr key={transaction?.signature}>
                  <td className="truncate max-w-[230px] text-blue-600 hover:underline">
                    <Link href={`/transaction/${transaction?.signature}`}>
                      {transaction?.signature}
                    </Link>
                  </td>
                  <td>{transaction?.slot}</td>
                  <td
                    className="whitespace-nowrap"
                    title={format(
                      fromUnixTime(transaction?.blockTime),
                      "MMMM d, yyyy 'at' HH:mm:ss OOOO"
                    )}
                  >
                    {formatDistanceToNow(fromUnixTime(transaction?.blockTime), {
                      includeSeconds: true,
                    })}
                  </td>
                  <td>
                    <span
                      className={`inline-block px-2 py-1 rounded-full text-xs font-bold leading-none text-white ${
                        transaction?.confirmationStatus === "finalized"
                          ? "bg-green-500"
                          : "bg-yellow-400"
                      }`}
                    >
                      {transaction?.confirmationStatus}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <div className="text-center">No transaction to display</div>
      )}
    </div>
  );
};

export default TransactionList;