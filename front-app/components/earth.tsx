// components/PointSphere.tsx
import React, { useRef, useEffect } from 'react';
import * as THREE from 'three';

const PointSphere: React.FC = () => {
  // 为 mountRef 指定类型为 HTMLDivElement | null
  const mountRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    // 确保 mountRef.current 不为 null
    if (mountRef.current === null) return;

    // 创建 Three.js 场景、相机和渲染器
    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
    const renderer = new THREE.WebGLRenderer();
    
    // 设置渲染器的大小并将其添加到 DOM 中
    renderer.setSize(window.innerWidth, window.innerHeight);
    mountRef.current.appendChild(renderer.domElement);

    // 创建球体几何体和点材质
    const sphereGeometry = new THREE.SphereGeometry(5, 32, 32);
    const pointMaterial = new THREE.PointsMaterial({
      color: 0x00FF00,
      size: 0.1
    });

    // 创建点云并添加到场景中
    const pointSphere = new THREE.Points(sphereGeometry, pointMaterial);
    scene.add(pointSphere);

    // 设置相机的位置
    camera.position.z = 10;

    // 处理窗口大小调整
    const handleResize = () => {
      if (mountRef.current) {
        camera.aspect = window.innerWidth / window.innerHeight;
        camera.updateProjectionMatrix();
        renderer.setSize(window.innerWidth, window.innerHeight);
      }
    };

    window.addEventListener('resize', handleResize);

    // 动画循环函数
    const animate = () => {
      requestAnimationFrame(animate);
      pointSphere.rotation.y += 0.01; // 旋转点状地球
      renderer.render(scene, camera); // 渲染场景和相机
    };

    animate();

    // 清理函数，用于组件卸载时清理 Three.js 的资源
    return () => {
      if (mountRef.current) {
        mountRef.current.removeChild(renderer.domElement); // 从 DOM 中移除渲染器
      }
      window.removeEventListener('resize', handleResize); // 移除窗口大小调整事件监听器
      renderer.dispose(); // 释放 Three.js 渲染器的资源
    };
  }, []); // 空数组作为依赖，确保 useEffect 仅在组件挂载和卸载时运行

  return <div ref={mountRef} />; // 将 mountRef 绑定到 div 元素
};

export default PointSphere;