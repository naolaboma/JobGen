"use client";
import { useForm } from "react-hook-form";
import { useState, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { Epilogue, Inter, Poppins } from "next/font/google";
import { apiUrl } from "@/lib/api";

const epilogue = Epilogue({ subsets: ["latin"], weight: ["400", "700"] });
const inter = Inter({ subsets: ["latin"], weight: ["400", "700"] });
const poppins = Poppins({ subsets: ["latin"], weight: ["400", "700"] });

type FormValues = {
  email: string;
  otp: string;
};

function VerifyEmailInner() {
  const form = useForm<FormValues>();
  const { register, handleSubmit, formState } = form;
  const { errors } = formState;
  const [serverError, setServerError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);
  const router = useRouter();
  const searchParams = useSearchParams();
  const email = searchParams.get("email") || "";
  const callbackUrl = searchParams.get("callbackUrl") || "/chat";

  const onSubmit = async (data: FormValues) => {
    setServerError("");
    setIsLoading(true);

    try {
      const res = await fetch(apiUrl("/api/v1/auth/verify-email"), {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: data.email,
          otp: data.otp,
        }),
      });

      const result = await res.json();

      if (!res.ok) {
        setServerError(result.message || "Verification failed");
      } else {
        setIsSuccess(true);
        setTimeout(() => {
          router.push(`/login?callbackUrl=${encodeURIComponent(callbackUrl)}`);
        }, 2000);
      }
    } catch (err) {
      console.error("Verification error:", err);
      setServerError("Network error");
    } finally {
      setIsLoading(false);
    }
  };

  const handleResendOTP = async () => {
    if (!email) return;

    setServerError("");
    setIsLoading(true);

    try {
      const res = await fetch(apiUrl("/api/v1/auth/resend-otp"), {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: email,
          purpose: "EMAIL_VERIFICATION",
        }),
      });

      if (!res.ok) {
        const result = await res.json();
        setServerError(result.message || "Failed to resend OTP");
      } else {
        setServerError("OTP sent successfully!");
      }
    } catch (err) {
      console.error("Resend error:", err);
      setServerError("Network error");
    } finally {
      setIsLoading(false);
    }
  };

  if (isSuccess) {
    return (
      <div className="max-w-md mx-auto mt-12 p-8 bg-white rounded shadow text-center">
        <div className="text-green-500 text-6xl mb-4">✓</div>
        <h1 className={`${inter.className} text-2xl font-bold mb-4 text-black`}>
          Email Verified!
        </h1>
        <p className={`${inter.className} text-gray-600 mb-6`}>
          Your email has been successfully verified. You can now log in to your
          account.
        </p>
        <p className="text-sm text-gray-500">Redirecting to login page...</p>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto mt-12 p-8 bg-white rounded shadow">
      <h1
        className={`${inter.className} text-2xl font-bold mb-6 text-center text-black`}
      >
        Verify Your Email
      </h1>
      <div className="text-center mb-6">
        <p className={`${inter.className} text-gray-600`}>
          We've sent a verification code to your email address.
        </p>
        <p className={`${inter.className} text-gray-600 font-semibold mt-2`}>
          {email || "your email"}
        </p>
      </div>

      <form
        className={`${epilogue.className}`}
        onSubmit={handleSubmit(onSubmit)}
        noValidate
      >
        <div className="mb-4">
          <label
            className="block mb-1 text-black font-semibold"
            htmlFor="email"
          >
            Email Address
          </label>
          <input
            className="w-full border rounded px-3 py-2 placeholder-gray-400 text-gray-400"
            type="email"
            id="email"
            defaultValue={email}
            placeholder="Enter your email address"
            {...register("email", {
              required: "Email is required",
              pattern: {
                value:
                  /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/,
                message: "Invalid email format",
              },
            })}
          />
          <p className="text-red-500 text-sm">{errors.email?.message}</p>
        </div>

        <div className="mb-4">
          <label className="block mb-1 text-black font-semibold" htmlFor="otp">
            Verification Code
          </label>
          <input
            className="w-full border rounded px-3 py-2 placeholder-gray-400 text-gray-400 text-center text-2xl tracking-widest"
            type="text"
            id="otp"
            placeholder="000000"
            maxLength={6}
            {...register("otp", {
              required: "Verification code is required",
              pattern: {
                value: /^[0-9]{6}$/,
                message: "Please enter a valid 6-digit code",
              },
            })}
          />
          <p className="text-red-500 text-sm">{errors.otp?.message}</p>
        </div>

        {serverError && (
          <p
            className={`text-sm mb-4 ${
              serverError.includes("successfully")
                ? "text-green-500"
                : "text-red-500"
            }`}
          >
            {serverError}
          </p>
        )}

        <button
          className="w-full bg-[#7BBFB3] text-white py-2 rounded-xl font-semibold disabled:opacity-50 mb-4"
          type="submit"
          disabled={isLoading}
        >
          {isLoading ? "Verifying..." : "Verify Email"}
        </button>

        <div className="text-center">
          <button
            type="button"
            onClick={handleResendOTP}
            className="text-[#7BBFB3] font-semibold hover:underline disabled:opacity-50"
            disabled={isLoading || !email}
          >
            Didn't receive the code? Resend
          </button>
        </div>
      </form>

      <p className="mt-6 text-gray-600 text-center">
        <Link href="/login" className="text-[#7BBFB3] font-semibold">
          Back to Login
        </Link>
      </p>
    </div>
  );
}

export default function VerifyEmailPage() {
  return (
    <Suspense
      fallback={
        <div className="max-w-md mx-auto mt-12 p-8 bg-white rounded shadow">
          Loading…
        </div>
      }
    >
      <VerifyEmailInner />
    </Suspense>
  );
}
