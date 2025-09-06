'use client';

import React, { useState, useRef } from 'react';
import { useSession } from "next-auth/react";

// --- UI Components ---

function ChatBubble({ text, byUser = false }: { text: string; byUser?: boolean }) {
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
}: {    title: string;    company: string;    location: string;    salary: string;    posted: string;    match: number;}) {
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
                        <circle cx="24" cy="24" r="20" stroke="#e5e7eb" strokeWidth="4" fill="none" />
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
                    <span className="absolute inset-0 flex items-center justify-center text-sm font-bold text-black">
                        {match}%
                    </span>
                </div>
            </div>
        </div>
    );
}

// --- Main Chat Component ---

export default function ChatBot() {
    const { data: session } = useSession();
    const [messages, setMessages] = useState([
        { type: 'bubble', text: "Hello! How can I help you today? Upload your resume (PDF) to get started.", byUser: false },
    ]);
    const [inputValue, setInputValue] = useState("");
    const [file, setFile] = useState<File | null>(null);
    const [uploading, setUploading] = useState(false);
    const [uploadStatus, setUploadStatus] = useState("");
    const fileInputRef = useRef<HTMLInputElement>(null);

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files) {
            setFile(e.target.files[0]);
            setUploadStatus("");
        }
    };

    const handleUpload = async () => {
        if (!file) return;
        if (!session) {
            setUploadStatus("You must be logged in to upload files.");
            return;
        }

        setUploading(true);
        setUploadStatus("Uploading...");

        const formData = new FormData();
        formData.append("file", file);

        try {
            const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/files/upload/document`, {
                method: "POST",
                headers: { "Authorization": `Bearer ${(session as any).accessToken}` },
                body: formData,
            });

            const result = await res.json();

            if (!res.ok) {
                setUploadStatus(`Error: ${result.message || "Upload failed"}`);
            } else {
                setUploadStatus(`Success: ${result.message || "File uploaded successfully!"}`);
                setFile(null);
            }
        } catch (err) {
            console.error("Upload error:", err);
            setUploadStatus("An unexpected error occurred during upload.");
        } finally {
            setUploading(false);
        }
    };

    const handleSendMessage = () => {
        if (inputValue.trim() === "") return;

        const newMessage = { type: 'bubble', text: inputValue, byUser: true };
        setMessages(prev => [...prev, newMessage]);
        setInputValue("");

        // --- TODO: Send message to backend and get response ---
    };

    return (
        <div className="flex h-screen bg-gray-100">
            {/* Sidebar */}
            <aside className="w-64 bg-white shadow-md p-4 flex flex-col">
                <h2 className="text-xl font-bold text-black mb-4">Chat History</h2>
                <div className="flex-grow space-y-2">
                    {/* Chat history will be populated here */}
                </div>
            </aside>

            {/* Main Content */}
            <main className="flex-1 flex flex-col">
                <div className="flex-grow p-4 overflow-y-auto space-y-4">
                    {messages.map((msg, index) => {
                        if (msg.type === 'bubble') {
                            return <ChatBubble key={index} text={msg.text} byUser={msg.byUser} />;
                        }
                        // Job cards would be rendered here
                        return null;
                    })}
                </div>

                {/* Input bar */}
                <div className="border-t p-4 flex items-center space-x-2 bg-white shadow-sm">
                    <input type="file" accept="application/pdf" onChange={handleFileChange} ref={fileInputRef} className="hidden" />
                    <button onClick={() => fileInputRef.current?.click()} className="bg-gray-200 text-black rounded-full w-10 h-10 flex items-center justify-center hover:bg-gray-300">
                        +
                    </button>
                    <input
                        type="text"
                        placeholder="ask ai for help"
                        value={inputValue}
                        onChange={(e) => setInputValue(e.target.value)}
                        onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
                        className="flex-1 px-4 py-2 rounded-full border-2 border-gray-300 focus:border-[#44C3BB] focus:outline-none text-black"
                    />
                    <button onClick={handleSendMessage} className="bg-[#44C3BB] text-white rounded-full w-10 h-10 flex items-center justify-center">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 12h14M12 5l7 7-7 7" /></svg>
                    </button>
                </div>
                {file && (
                    <div className="flex items-center justify-between p-4 border-t bg-white">
                        <p className="text-sm text-gray-600">Selected file: {file.name}</p>
                        <button
                            onClick={handleUpload}
                            className="bg-green-500 text-white rounded-lg px-4 py-1 text-sm hover:bg-green-600 disabled:opacity-50"
                            disabled={uploading}
                        >
                            {uploading ? "Uploading..." : "Upload"}
                        </button>
                    </div>
                )}
            </main>
        </div>
    );
}