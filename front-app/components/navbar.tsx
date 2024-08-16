'use client';

import Link from 'next/link';
import { useSession, signOut } from 'next-auth/react';
import React from "react";

export default function Navbar () {
    const { data: session, status } = useSession();

    return (
        <div className="navbar bg-base-100">
            <div className="navbar-start">
                <div className="dropdown">
                    <div tabIndex={0} role="button" className="btn btn-ghost lg:hidden">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            className="h-5 w-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor">
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth="2"
                                d="M4 6h16M4 12h8m-8 6h16" />
                        </svg>
                    </div>
                </div>
                <a className="btn btn-ghost text-xl">
                    多链农产溯源通     
                </a>
            </div>
            <div className="navbar-center hidden lg:flex">
                <ul className="menu menu-horizontal px-1">
                    <li>
                        <Link href={"/mainchain"}>
                        主链
                        </Link>
                    </li>
                    <li>
                        <details>
                            <summary>侧链1</summary>
                            <ul className="p-2">
                                <li><a>子菜单 1</a></li>
                                <li><a>子菜单 2</a></li>
                            </ul>
                        </details>
                    </li>
                    <li>
                        <Link href={"/sidechain2"}>
                        侧链2
                        </Link>
                    </li>
                </ul>
            </div>
            <div className="navbar-end">
                {status === "authenticated" ? (
                    <div className="flex items-center gap-2">
                        {session.user?.image ? (
                            <img
                                src={session.user.image}
                                alt="User Avatar"
                                className="w-10 h-10 rounded-full"
                            />
                        ) : (
                            <div className="w-10 h-10 rounded-full bg-gray-300 flex items-center justify-center">
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    className="w-6 h-6 text-gray-500"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor">
                                    <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth="2"
                                        d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zM12 14c-4.41 0-8 2.59-8 6v1h16v-1c0-3.41-3.59-6-8-6z" />
                                </svg>
                            </div>
                        )}
                        <button className="btn" onClick={() => signOut()}>Sign out</button>
                    </div>
                ) : (
                    <div className="flex items-center">
                        <Link className="btn" href={"/sign-in"}>
                          Sign in
                        </Link>
                    </div>
                )}
            </div>
        </div>
    );
};