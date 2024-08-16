import React from "react";

interface SearchTransactionFormProps {
  handleFormSubmit: (event: React.FormEvent<HTMLFormElement>) => void;
  address: string;
  loading: boolean;
  setAddress: (address: string) => void;
  errorMessage?: string;
}

const SearchTransactionForm: React.FC<SearchTransactionFormProps> = ({
  handleFormSubmit,
  address,
  loading,
  setAddress,
  errorMessage,
}) => {
  return (
    <form onSubmit={handleFormSubmit} className="flex flex-wrap w-full">
      <label htmlFor="address" className="w-full shrink-0 text-lg mb-2">
        Transaction address
      </label>
      <input
        type="text"
        name="address"
        value={address}
        onChange={(event) => setAddress(event.target.value)}
        className="w-3/4 border-2 border-r-0 border-gray-500 h-12 rounded-l-lg px-4 focus:outline-none focus:border-blue-600 disabled:bg-gray-500 transition-colors duration-150"
        placeholder="输入交易地址"
        disabled={loading}
        required
      />
      <button
        type="submit"
        disabled={loading}
        className="flex-grow bg-sky-600 flex items-center justify-center rounded-r-lg text-white text-sm hover:bg-blue-900 disabled:bg-gray-500 transition-colors duration-150"
      >
        查询
      </button>
      {errorMessage && (
        <p className="text-red-600 text-base my-1">{errorMessage}</p>
      )}
    </form>
  );
};

export default SearchTransactionForm;