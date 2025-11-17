"use client";

import { Suspense } from "react";
import Image from "next/image";
import SignInForm from "../components/SignInForm";

export default function LoginPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-teal-50 via-white to-cyan-50">
      <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8 py-10">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
          <div className="hidden lg:block">
            <div className="relative w-full h-[460px] rounded-3xl bg-white/60 backdrop-blur shadow-sm border border-gray-100 overflow-hidden">
              <div className="absolute inset-0 bg-gradient-to-br from-[#44C3BB]/20 via-transparent to-[#44C3BB]/10" />
              <div className="absolute inset-0 p-8 flex flex-col justify-between">
                <div>
                  <h2 className="text-3xl font-extrabold text-black tracking-tight">
                    Welcome back to JobGen
                  </h2>
                  <p className="mt-3 text-gray-600 max-w-md">
                    Sign in to personalize your job search and get AI-powered CV
                    insights.
                  </p>
                </div>
                <div className="relative w-full h-64">
                  <Image
                    src="/professional-woman-dark-hair.png"
                    alt="Welcome illustration"
                    fill
                    className="object-contain object-bottom opacity-90"
                    priority
                  />
                </div>
              </div>
            </div>
          </div>

          <div>
            <Suspense
              fallback={
                <div className="max-w-md mx-auto mt-6 p-8 bg-white rounded-2xl shadow-sm border border-gray-100">
                  Loading…
                </div>
              }
            >
              <SignInForm />
            </Suspense>
          </div>
        </div>
      </div>
    </div>
  );
}
