"use client";

import { Suspense } from "react";
import SignUpForm from "../components/SignUpForm";

export default function RegisterPage() {
  return (
    <Suspense
      fallback={
        <div className="max-w-md mx-auto mt-12 p-8 bg-white rounded shadow">
          Loading…
        </div>
      }
    >
      <SignUpForm />
    </Suspense>
  );
}
