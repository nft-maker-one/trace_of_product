import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Navbar from "@/components/navbar";
import React, { ReactNode } from "react";
import { NextAuthProvider } from "@/components/Providers";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "test",
  description: "",
};

interface LayoutProps {
  children: ReactNode;
}

// 修正 Layout 组件的定义
const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="container mx-auto">
      <Navbar />
      {children}
    </div>
  );
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <NextAuthProvider>
        <Layout>{children}</Layout> {/* 使用修正后的 Layout */}
        </NextAuthProvider>
      </body>
    </html>
  );
}