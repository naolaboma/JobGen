"use client";
import { useSession } from "next-auth/react";
import Link from "next/link";
import ChatBotAfterWeb from "../chat-bot-after-web/chat-bot-after-web";
import NavBar from "../components/NavBar";

export default function ChatPage() {
    const { data: session, status } = useSession();

    // Temporarily commented out for UI inspection
    // if (status === "loading") {
    //     return <div className="flex items-center justify-center h-screen">Loading...</div>;
    // }

    // if (!session) {
    //     return (
    //         <div className="flex flex-col items-center justify-center h-screen">
    //             <p>You must be logged in to access this page.</p>
    //             <Link href="/login" className="text-blue-500 hover:underline">
    //                 Go to login
    //             </Link>
    //         </div>
    //     );
    // }

    return (
        <div>
            <NavBar />
            <ChatBotAfterWeb />
        </div>
    );
}
