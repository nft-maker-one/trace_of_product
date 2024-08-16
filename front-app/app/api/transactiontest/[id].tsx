import Head from "next/head";
import { useState, useEffect } from "react";
import axios, { AxiosResponse } from "axios";
import { useRouter } from "next/router";
import TransactionListDetail from "@/components/TransactionList";

// 定义交易数据的类型
interface TransactionData {
  transactionHash: string;
  /*
  type Eggplant struct {
	MetaData
	PublickKey crypto.PublicKey  //the Validator of this Eggplant
	Signature  *crypto.Signature //the Signature of the Validator
	hash       types.Hash        //the digest for eggplant's metadata
	firstSeen  int64             //the creation time of the eggplant
}

  */
}

interface ApiResponse {
  transaction: TransactionData;
}

export default function TransactionDetail() {
  const [loading, setLoading] = useState(false);
  const [transactionData, setTransactionData] = useState<TransactionData | null>(null);
  const [errorMessage, setErrorMessage] = useState("");
  const router = useRouter();

  useEffect(() => {
    const getTransaction = async () => {
      try {
        setLoading(true);
        setErrorMessage("");

        const response: AxiosResponse<ApiResponse> = await axios.post("/api/transaction", {
          transactionHash: router.query?.id,
        });

        if (response.status === 200) {
          setTransactionData(response.data.transaction);
        }
      } catch (error) {
        setErrorMessage(
          (error as any)?.response?.data?.message ||
            "无法获取交易信息，请稍后再试。"
        );
      } finally {
        setLoading(false);
      }
    };

    getTransaction();
  }, [router.query?.id]);

  return (
    <>
      <Head>
        <title>Solana 区块链浏览器：交易详情</title>
      </Head>
      <main className="w-full h-full p-6 flex flex-col items-center justify-between gap-6 mx-auto relative">
        <h1 className="text-2xl">交易详情</h1>
        {errorMessage && (
          <p className="text-red-600 text-base text-center my-1">
            {errorMessage}
          </p>
        )}

        <TransactionListDetail
          loading={loading}
          transactionData={transactionData}
        />

        {loading && (
          <div className="absolute inset-0 bg-white/70 flex items-center justify-center">
            加载中...
          </div>
        )}
      </main>
    </>
  );
}