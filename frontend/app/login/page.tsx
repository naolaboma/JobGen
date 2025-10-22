"use client";

import { Suspense } from "react";
import SignInForm from "../components/SignInForm";

export default function LoginPage() {
  return (
    <Suspense
      fallback={
        <div className="max-w-md mx-auto mt-12 p-8 bg-white rounded shadow">
          Loading…
        </div>
      }
    >
      <SignInForm />
    </Suspense>
  );
}
