'use client';
import React from 'react';
import dynamic from 'next/dynamic';

// 使用 dynamic 以禁用 SSR（服务器端渲染）
const PointSphere = dynamic(() => import('@/components/earth'), { ssr: false });

const HomePage = () => {
  return (
    <div>
      <PointSphere />
    </div>
  );
};

export default HomePage;