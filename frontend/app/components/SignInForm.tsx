"use client";
import { useForm } from "react-hook-form";
import Link from "next/link";
import { useState } from "react";
import { signIn } from "next-auth/react";
import { useSearchParams, useRouter } from "next/navigation";

type FormValues = {
  email: string;
  password: string;
};

export default function SignInForm() {
  const form = useForm<FormValues>();
  const { register, handleSubmit, formState } = form;
  const { errors } = formState;
  const [serverError, setServerError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const router = useRouter();
  const searchParams = useSearchParams();
  const callbackUrl = searchParams.get("callbackUrl") || "/";

  const onSubmit = async (data: FormValues) => {
    setServerError("");
    setIsLoading(true);

    try {
      const res = await signIn("credentials", {
        redirect: false,
        email: data.email,
        password: data.password,
        callbackUrl,
      });

      if (res?.error) {
        setServerError(res.error);
      } else if (res?.url) {
        router.push(res.url);
      } else {
        router.push(callbackUrl);
      }
    } catch (err) {
      console.error("Signin error:", err);
      setServerError("An unexpected error occurred.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleSocialSignIn = (provider: "google" | "github") => {
    signIn(provider, { callbackUrl });
  };

  return (
    <div className="max-w-md mx-auto mt-6 bg-white rounded-2xl shadow-sm border border-gray-100 p-6 sm:p-8">
      <div className="mb-6 text-center">
        <h1 className="text-2xl font-extrabold tracking-tight text-black">
          Sign in
        </h1>
        <p className="mt-2 text-sm text-gray-600">
          Welcome back. Let’s continue your search.
        </p>
      </div>

      <div className="grid gap-3 mb-5">
        <button
          type="button"
          onClick={() => handleSocialSignIn("google")}
          className="w-full inline-flex items-center justify-center gap-2 rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm font-semibold text-black hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-[#44C3BB]/40"
        >
          <span className="inline-block">
            {/* Google icon */}
            <svg
              width="18"
              height="18"
              viewBox="0 0 21 20"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M18.6712 8.36788H18V8.33329H10.5V11.6666H15.2096C14.5225 13.607 12.6762 15 10.5 15C7.73874 15 5.49999 12.7612 5.49999 9.99996C5.49999 7.23871 7.73874 4.99996 10.5 4.99996C11.7746 4.99996 12.9342 5.48079 13.8171 6.26621L16.1742 3.90913C14.6858 2.52204 12.695 1.66663 10.5 1.66663C5.89791 1.66663 2.16666 5.39788 2.16666 9.99996C2.16666 14.602 5.89791 18.3333 10.5 18.3333C15.1021 18.3333 18.8333 14.602 18.8333 9.99996C18.8333 9.44121 18.7758 8.89579 18.6712 8.36788Z"
                fill="#FFC107"
              />
              <path
                d="M3.12749 6.12121L5.8654 8.12913C6.60624 6.29496 8.4004 4.99996 10.5 4.99996C11.7746 4.99996 12.9342 5.48079 13.8171 6.26621L16.1742 3.90913C14.6858 2.52204 12.695 1.66663 10.5 1.66663C7.29915 1.66663 4.52332 3.47371 3.12749 6.12121Z"
                fill="#FF3D00"
              />
              <path
                d="M10.5 18.3333C12.6525 18.3333 14.6083 17.5095 16.0871 16.17L13.5079 13.9875C12.6432 14.6451 11.5865 15.0008 10.5 15C8.33251 15 6.49209 13.6179 5.79876 11.6891L3.08126 13.7829C4.46043 16.4816 7.26126 18.3333 10.5 18.3333Z"
                fill="#4CAF50"
              />
              <path
                d="M18.6713 8.36796H18V8.33337H10.5V11.6667H15.2096C14.8809 12.5902 14.2889 13.3972 13.5067 13.988L13.5079 13.9871L16.0871 16.1696C15.9046 16.3355 18.8333 14.1667 18.8333 10C18.8333 9.44129 18.7758 8.89587 18.6713 8.36796Z"
                fill="#1976D2"
              />
            </svg>
          </span>
          Continue with Google
        </button>
        <button
          type="button"
          onClick={() => handleSocialSignIn("github")}
          className="w-full inline-flex items-center justify-center gap-2 rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm font-semibold text-black hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-[#44C3BB]/40"
        >
          <span className="inline-block">
            {/* GitHub icon */}
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              fill="currentColor"
              viewBox="0 0 16 16"
            >
              <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27s1.36.09 2 .27c1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8" />
            </svg>
          </span>
          Continue with GitHub
        </button>
      </div>

      <div className="relative flex items-center gap-3 my-6">
        <div className="flex-1 h-px bg-gray-200" />
        <span className="text-xs text-gray-500">or with email</span>
        <div className="flex-1 h-px bg-gray-200" />
      </div>

      <form onSubmit={handleSubmit(onSubmit)} noValidate className="space-y-4">
        <div>
          <label
            className="block mb-1 text-sm text-black font-medium"
            htmlFor="email"
          >
            Email
          </label>
          <input
            id="email"
            type="email"
            autoComplete="email"
            placeholder="you@example.com"
            className="w-full rounded-xl border border-gray-300 bg-white px-3.5 py-2.5 text-sm text-black placeholder-gray-400 shadow-xs focus:outline-none focus:ring-2 focus:ring-[#44C3BB]/40 focus:border-[#44C3BB]"
            {...register("email", {
              required: "Email is required",
              pattern: {
                value:
                  /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/,
                message: "Invalid email format",
              },
            })}
          />
          {errors.email?.message && (
            <p className="mt-1 text-xs text-red-600">{errors.email.message}</p>
          )}
        </div>

        <div>
          <label
            className="block mb-1 text-sm text-black font-medium"
            htmlFor="password"
          >
            Password
          </label>
          <div className="relative">
            <input
              id="password"
              type={showPassword ? "text" : "password"}
              autoComplete="current-password"
              placeholder="••••••••"
              className="w-full rounded-xl border border-gray-300 bg-white px-3.5 py-2.5 pr-10 text-sm text-black placeholder-gray-400 shadow-xs focus:outline-none focus:ring-2 focus:ring-[#44C3BB]/40 focus:border-[#44C3BB]"
              {...register("password", {
                required: "Password is required",
                minLength: {
                  value: 8,
                  message: "Password must be at least 8 characters",
                },
              })}
            />
            <button
              type="button"
              onClick={() => setShowPassword((v) => !v)}
              aria-label={showPassword ? "Hide password" : "Show password"}
              className="absolute inset-y-0 right-0 px-3 text-gray-500 hover:text-gray-700 focus:outline-none"
            >
              {showPassword ? (
                // Eye-off
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-5 w-5"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                >
                  <path d="M3 3l18 18" />
                  <path d="M10.58 10.58A2 2 0 0 0 12 14a2 2 0 0 0 1.42-.58" />
                  <path d="M16.24 16.24A10.94 10.94 0 0 1 12 18c-5 0-9-4-9-6a10.94 10.94 0 0 1 4.53-5.94" />
                  <path d="M14.12 9.88A2 2 0 0 0 12 8a2 2 0 0 0-2 2" />
                  <path d="M17.94 13.06A10.94 10.94 0 0 0 21 12c0-2-4-6-9-6a10.94 10.94 0 0 0-1.18.06" />
                </svg>
              ) : (
                // Eye
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-5 w-5"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                >
                  <path d="M1 12s4-7 11-7 11 7 11 7-4 7-11 7-11-7-11-7Z" />
                  <circle cx="12" cy="12" r="3" />
                </svg>
              )}
            </button>
          </div>
          {errors.password?.message && (
            <p className="mt-1 text-xs text-red-600">
              {errors.password.message}
            </p>
          )}
        </div>

        {serverError && (
          <div className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700">
            {serverError}
          </div>
        )}

        <button
          className="w-full inline-flex items-center justify-center gap-2 rounded-xl bg-[#44C3BB] text-white py-2.5 font-semibold shadow-sm hover:bg-[#3bb3ac] focus:outline-none focus:ring-2 focus:ring-[#44C3BB]/40 disabled:opacity-60"
          type="submit"
          disabled={isLoading}
        >
          {isLoading && (
            <svg className="h-4 w-4 animate-spin" viewBox="0 0 24 24">
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              />
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
              />
            </svg>
          )}
          {isLoading ? "Signing in..." : "Sign in"}
        </button>
      </form>

      <div className="mt-5 flex items-center justify-between text-sm">
        <Link href="/support" className="text-gray-600 hover:text-black">
          Forgot password?
        </Link>
        <p className="text-gray-600">
          No account?{" "}
          <Link
            href={`/register?callbackUrl=${encodeURIComponent(callbackUrl)}`}
            className="text-[#44C3BB] font-semibold hover:underline"
          >
            Sign up
          </Link>
        </p>
      </div>
    </div>
  );
}
