// App.tsx
"use client";

import { useState } from "react";

interface MessageProps {
    text: string;
    byUser?: boolean;
}

function ChatBubble({ text, byUser = false }: MessageProps) {
    return (
        <div className={`flex ${byUser ? "justify-end" : "justify-start"}`}>
            <div
                className={`relative max-w-[75%] px-4 py-2 rounded-2xl ${byUser
                    ? "bg-[#44C3BB] text-white"
                    : "bg-gray-200 text-black"
                    }`}
            >
                <p className="whitespace-pre-wrap text-sm leading-relaxed">{text}</p>
            </div>
        </div>
    );
}

function JobCard({
    title,
    company,
    location,
    salary,
    posted,
    match,
}: {
    title: string;
    company: string;
    location: string;
    salary: string;
    posted: string;
    match: number;
}) {
    const circleLength = 2 * Math.PI * 20;

    return (
        <div className="flex items-center justify-between bg-gray-100 rounded-xl p-4 shadow-sm">
            <div>
                <h3 className="font-bold text-black">{title}</h3>
                <p className="text-sm text-gray-600">
                    {company} • {location}
                </p>
                <p className="text-sm text-gray-500">{salary}</p>
                <p className="text-xs text-gray-400">{posted}</p>
            </div>
            <div className="flex flex-col items-center">
                <div className="relative w-12 h-12">
                    <svg className="w-12 h-12 transform -rotate-90">
                        <circle
                            cx="24"
                            cy="24"
                            r="20"
                            stroke="#e5e7eb"
                            strokeWidth="4"
                            fill="none"
                        />
                        <circle
                            cx="24"
                            cy="24"
                            r="20"
                            stroke="#44C3BB"
                            strokeWidth="4"
                            fill="none"
                            strokeDasharray={circleLength}
                            strokeDashoffset={circleLength - (match / 100) * circleLength}
                        />
                    </svg>
                    <span className="absolute inset-0 flex items-center justify-center text-xs font-bold text-black">
                        {match}%
                    </span>
                </div>
                <p className="text-xs mt-1 text-black">Match</p>
            </div>
        </div>
    );
}

export default function App() {
    const [messages, setMessages] = useState([
        { text: "Lorem ipsum dolor sit amet", byUser: true },
        {
            text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
            byUser: false,
        },
        {
            text: "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
            byUser: true,
        },
    ]);

    return (
        <div className="h-screen flex bg-gray-50">
            {/* Sidebar */}
            <aside className="w-64 bg-gray-100 flex flex-col p-4 border-r">
                <button className="font-bold mb-4 text-left text-black">＋ New chat</button>

                <h2 className="text-sm text-gray-600 mb-2">Chats</h2>
                <div className="space-y-2 mb-4">
                    <button className="w-full bg-[#44C3BB] text-white rounded-lg px-2 py-1 text-sm text-left">
                        Excepteur sint occaecat cupidatat non proident
                    </button>
                    <button className="w-full bg-[#44C3BB] text-white rounded-lg px-2 py-1 text-sm text-left">
                        Lorem ipsum dolor sit amet
                    </button>
                </div>

                <h2 className="text-sm text-gray-600 mb-2">Previous</h2>
                <div className="space-y-2 mb-4">
                    <button className="w-full bg-[#44C3BB] text-white rounded-lg px-2 py-1 text-sm text-left">
                        Dolor in reprehenderit
                    </button>
                </div>

                <div className="mt-auto flex items-center space-x-2">
                    <div className="w-8 h-8 rounded-full bg-black text-white flex items-center justify-center">
                        J
                    </div>
                    <span className="text-sm">John@gmail.com</span>
                </div>
            </aside>

            {/* Main Chat Area */}
            <main className="flex-1 flex flex-col">
                <div className="flex-1 overflow-y-auto p-6 space-y-6">
                    {messages.map((msg, i) => (
                        <ChatBubble key={i} text={msg.text} byUser={msg.byUser} />
                    ))}

                    <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
                        <JobCard
                            title="Backend Engineer"
                            company="Acme Corp"
                            location="Remote"
                            salary="$100,000"
                            posted="2 days ago"
                            match={85}
                        />
                        <JobCard
                            title="Frontend Engineer"
                            company="Beta Ltd"
                            location="Remote"
                            salary="$90,000"
                            posted="1 day ago"
                            match={72}
                        />
                        <JobCard
                            title="Fullstack Developer"
                            company="Gamma Inc"
                            location="Remote"
                            salary="$110,000"
                            posted="3 days ago"
                            match={93}
                        />
                    </div>
                </div>

                {/* Input bar */}
                <div className="border-t p-4 flex items-center space-x-2 bg-white shadow-sm">
                    <input
                        type="text"
                        placeholder="ask ai for help"
                        className="flex-1 px-4 py-2 rounded-full border-2 border-gray-300 focus:border-[#44C3BB] focus:outline-none text-black"
                    />
                    <button className="bg-[#44C3BB] text-white rounded-full w-10 h-10 flex items-center justify-center">
                        ＋
                    </button>
                </div>
            </main>
        </div>
    );
}
