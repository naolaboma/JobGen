"use client";
import { useForm } from 'react-hook-form';
import Link from 'next/link';
import { useState } from 'react';
import { signIn } from "next-auth/react";
import { useSearchParams, useRouter } from "next/navigation";
import { Epilogue, Inter, Poppins } from 'next/font/google';

const epilogue = Epilogue({ subsets: ['latin'], weight: ['400', '700'] });
const inter = Inter({ subsets: ['latin'], weight: ['400', '700'] });
const poppins = Poppins({ subsets: ['latin'], weight: ['400', '700'] });

type FormValues = {
    fullName: string;
    username: string;
    email: string;
    password: string;
    confirmPassword: string;
};

export default function SignUpForm() {
    const form = useForm<FormValues>();
    const { register, handleSubmit, formState, watch } = form;
    const { errors } = formState;
    const [serverError, setServerError] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const router = useRouter();
    const searchParams = useSearchParams();
    const callbackUrl = searchParams.get("callbackUrl") || "/chat";

    const onSubmit = async (data: FormValues) => {
        setServerError("");
        setIsLoading(true);

        try {
            const res = await fetch("http://localhost:8080/api/v1/auth/register", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    email: data.email,
                    full_name: data.fullName,
                    username: data.username,
                    password: data.password,
                }),
            });

            const result = await res.json();
            console.log("Signup response:", result);

            if (!res.ok) {
                setServerError(result.message || "Signup failed");
            } else {
                router.push(callbackUrl); // Redirect directly to chat or callbackUrl
            }
        } catch (err) {
            console.error("Signup error:", err);
            setServerError("Network error");
        } finally {
            setIsLoading(false);
        }
    };

    const handleSocialSignIn = (provider: 'google' | 'github') => {
        signIn(provider, { callbackUrl });
    };

    return (
        <div className="max-w-md mx-auto mt-12 p-8 bg-white rounded shadow">
            <h1 className={`${inter.className} text-2xl font-bold mb-6 text-center text-black`}>Create an Account</h1>
            <div className="text-center mb-4 flex items-center justify-center gap-2">
                <p className={`${inter.className} text-black whitespace-nowrap`}>Sign up to continue your job search journey</p>
            </div>
            <form className={`${epilogue.className}`} onSubmit={handleSubmit(onSubmit)} noValidate>
                <button
                    type="button"
                    onClick={() => handleSocialSignIn('google')}
                    className={`${inter.className} font-semibold w-full flex items-center justify-center gap-2 bg-white border border-gray-300 rounded px-4 py-2 mb-4 hover:bg-gray-50 text-black`}
                >
                    <span className="inline-block">
                        <svg width="21" height="20" viewBox="0 0 21 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M18.6712 8.36788H18V8.33329H10.5V11.6666H15.2096C14.5225 13.607 12.6762 15 10.5 15C7.73874 15 5.49999 12.7612 5.49999 9.99996C5.49999 7.23871 7.73874 4.99996 10.5 4.99996C11.7746 4.99996 12.9342 5.48079 13.8171 6.26621L16.1742 3.90913C14.6858 2.52204 12.695 1.66663 10.5 1.66663C5.89791 1.66663 2.16666 5.39788 2.16666 9.99996C2.16666 14.602 5.89791 18.3333 10.5 18.3333C15.1021 18.3333 18.8333 14.602 18.8333 9.99996C18.8333 9.44121 18.7758 8.89579 18.6712 8.36788Z" fill="#FFC107" />
                            <path d="M3.12749 6.12121L5.8654 8.12913C6.60624 6.29496 8.4004 4.99996 10.5 4.99996C11.7746 4.99996 12.9342 5.48079 13.8171 6.26621L16.1742 3.90913C14.6858 2.52204 12.695 1.66663 10.5 1.66663C7.29915 1.66663 4.52332 3.47371 3.12749 6.12121Z" fill="#FF3D00" />
                            <path d="M10.5 18.3333C12.6525 18.3333 14.6083 17.5095 16.0871 16.17L13.5079 13.9875C12.6432 14.6451 11.5865 15.0008 10.5 15C8.33251 15 6.49209 13.6179 5.79876 11.6891L3.08126 13.7829C4.46043 16.4816 7.26126 18.3333 10.5 18.3333Z" fill="#4CAF50" />
                            <path d="M18.6713 8.36796H18V8.33337H10.5V11.6667H15.2096C14.8809 12.5902 14.2889 13.3972 13.5067 13.988L13.5079 13.9871L16.0871 16.1696C15.9046 16.3355 18.8333 14.1667 18.8333 10C18.8333 9.44129 18.7758 8.89587 18.6713 8.36796Z" fill="#1976D2" />
                        </svg>
                    </span>
                    Sign Up with Google
                </button>
                <button
                    type="button"
                    onClick={() => handleSocialSignIn('github')}
                    className={`${inter.className} font-semibold w-full flex items-center justify-center gap-2 bg-white border border-gray-300 rounded px-4 py-2 mb-4 hover:bg-gray-50 text-black`}
                >
                    <span className="inline-block">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27s1.36.09 2 .27c1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8" />
                        </svg>
                    </span>
                    Sign Up with GitHub
                </button>

                <div className="text-center mb-4 flex items-center justify-center gap-2">
                    <div className="flex-1 border-b border-gray-300"></div>
                    <p className="text-gray-500 whitespace-nowrap">Or sign up with email</p>
                    <div className="flex-1 border-b border-gray-300"></div>
                </div>

                <div className="mb-4">
                    <label className="block mb-1 text-black font-semibold" htmlFor="fullName">Full Name</label>
                    <input
                        className="w-full border rounded px-3 py-2 placeholder-gray-400 text-gray-400"
                        type="text"
                        id="fullName"
                        placeholder="Enter your full name"
                        {...register("fullName", {
                            required: "Full name is required",
                        })}
                    />
                    <p className="text-red-500 text-sm">{errors.fullName?.message}</p>
                </div>

                <div className="mb-4">
                    <label className="block mb-1 text-black font-semibold" htmlFor="username">Username</label>
                    <input
                        className="w-full border rounded px-3 py-2 placeholder-gray-400 text-gray-400"
                        type="text"
                        id="username"
                        placeholder="Enter your username"
                        {...register("username", {
                            required: "Username is required",
                        })}
                    />
                    <p className="text-red-500 text-sm">{errors.username?.message}</p>
                </div>

                <div className="mb-4">
                    <label className="block mb-1 text-black font-semibold" htmlFor="email">Email Address</label>
                    <input
                        className="w-full border rounded px-3 py-2 placeholder-gray-400 text-gray-400"
                        type="email"
                        id="email"
                        placeholder="Enter email address"
                        {...register("email", {
                            required: "Email is required",
                            pattern: {
                                value: /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/,
                                message: "Invalid email format",
                            }
                        })}
                    />
                    <p className="text-red-500 text-sm">{errors.email?.message}</p>
                </div>
                <div className="mb-4">
                    <label className="block mb-1 text-black font-semibold" htmlFor="password">Password</label>
                    <input
                        className="w-full border rounded px-3 py-2 placeholder-gray-400 text-gray-400"
                        type="password"
                        id="password"
                        placeholder="Enter password"
                        {...register("password", {
                            required: "Password is required",
                            minLength: { value: 8, message: "Password must be at least 8 characters" } // Relaxed to minLength 8
                        })}
                    />
                    <p className="text-red-500 text-sm">{errors.password?.message}</p>
                </div>
                <div className="mb-4">
                    <label className="block mb-1 text-black font-semibold" htmlFor="confirmPassword">Confirm Password</label>
                    <input
                        className="w-full border rounded px-3 py-2 placeholder-gray-400 text-gray-400"
                        type="password"
                        id="confirmPassword"
                        placeholder="Confirm password"
                        {...register("confirmPassword", {
                            required: "Confirm Password is required",
                            validate: value => value === watch("password") || "Passwords do not match"
                        })}
                    />
                    <p className="text-red-500 text-sm">{errors.confirmPassword?.message}</p>
                </div>
                {serverError && <p className="text-red-500 text-sm mb-2">{serverError}</p>}
                <button
                    className="w-full bg-[#7BBFB3] text-white py-2 rounded-xl font-semibold disabled:opacity-50"
                    type="submit"
                    disabled={isLoading}
                >
                    {isLoading ? "Creating Account..." : "Sign Up"}
                </button>
            </form>
            <p className="mt-4 text-gray-600 text-center">
                Already have an account? <Link href={`/login?callbackUrl=${encodeURIComponent(callbackUrl)}`} className="text-[#7BBFB3] font-semibold">Sign In</Link>
            </p>

        </div>
    );
}
