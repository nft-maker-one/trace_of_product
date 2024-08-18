'use client';
import Head from "next/head";
import { useState } from "react";
import axios from "axios";
import TransactionList from "@/components/TransactionList";
import SearchTransactionForm from "@/components/SearchTransactionForm";

// 定义接口以明确响应数据的结构
interface TransactionResponse {
  transactionList: any[];
  accountBalance: string;
}

export default function Home() {
  const [loading, setLoading] = useState(false);
  const [transactionList, setTransactionList] = useState<any[]>([]);
  const [balance, setBalance] = useState<number | null>(null);
  const [address, setAddress] = useState<string>("");
  const [errorMessage, setErrorMessage] = useState<string>("");

  const handleFormSubmit = async (event: React.FormEvent) => {
    try {
      event.preventDefault();
      setLoading(true);
      setErrorMessage("");

      const response = await axios.get<TransactionResponse>(`/api/transactions/?address=${address}`);
      if (response.status === 200) {
        setTransactionList(response.data.transactionList);

        const accountBalanceText = response.data.accountBalance;
        const accountBalance = parseInt(accountBalanceText) / 1_000_000_000;

        if (!isNaN(accountBalance)) {
          setBalance(accountBalance);
        } else {
          setBalance(null);
        }
      }
    } catch (error: any) {
      console.log("client", error);
      setErrorMessage(
        error?.response?.data?.message ||
          "Unable to fetch transactions. Please try again later."
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Head>
        <title>Solana 区块链浏览器</title>
      </Head>
      <main className="w-full h-full max-w-2xl p-6 flex flex-col items-center justify-between gap-6 mx-auto relative">
        <h1 className="text-2xl">Solana 区块链浏览器</h1>
        <SearchTransactionForm
          handleFormSubmit={handleFormSubmit}
          address={address}
          setAddress={setAddress}
          loading={loading}
          errorMessage={errorMessage}
        />

        {/* 将 balance ?? 0 传递给 TransactionList */}
        <TransactionList transactionList={transactionList} balance={balance ?? 0} />

        {loading && (
          <div className="absolute inset-0 bg-white/70 flex items-center justify-center">
            Loading
          </div>
        )}
      </main>
    </>
  );
}