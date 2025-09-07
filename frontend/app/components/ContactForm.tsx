"use client";

import React from "react";
import { useForm, SubmitHandler } from "react-hook-form";
import { useDispatch, useSelector } from "react-redux";
import { AppDispatch, RootState } from "@/store/store";
import {
  submitContactFormStart,
  submitContactFormSuccess,
  submitContactFormFailure,
  clearContactFormStatus,
} from "@/lib/redux/slices/contactSlice";
import { useSubmitContactFormMutation } from "@/lib/redux/api/ContactApi";

const ContactForm: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const { isLoading, error, isSubmitted } = useSelector(
    (state: RootState) => state.contact
  );

  const [
    submitContactForm,
    { isLoading: mutationLoading, error: mutationError },
  ] = useSubmitContactFormMutation();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<IFormInput>();

  const onSubmit: SubmitHandler<IFormInput> = async (data) => {
    dispatch(submitContactFormStart());
    try {
      await submitContactForm(data).unwrap();
      dispatch(submitContactFormSuccess(data));
      reset();
      alert("Form submitted successfully! The admin has been notified.");
    } catch (err: any) {
      const errorMessage = err?.data?.message || "Submission failed";
      dispatch(submitContactFormFailure(errorMessage));
      alert(`Submission failed: ${errorMessage}`);
    }
  };

  React.useEffect(() => {
    dispatch(clearContactFormStatus());
    return () => {
      dispatch(clearContactFormStatus());
    };
  }, [dispatch]);

  return (
    <div className="max-w-2xl mx-auto mt-12 p-6 bg-white dark:bg-gray-800 rounded-lg shadow-2xl border border-gray-200 dark:border-gray-700">
      <h1 className="text-2xl font-bold mb-6 text-center text-gray-900 dark:text-gray-100">
        Contact Us
      </h1>
      <h3 className="font-bold text-center mb-12 text-gray-700 dark:text-gray-300">
        Have a question or feedback? We’re here to help!
      </h3>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {/* Name Field */}
        <div className="mb-7">
          <label
            htmlFor="name"
            className="block text-sm font-bold text-gray-700 dark:text-gray-300"
          >
            Name
          </label>
          <input
            id="name"
            type="text"
            placeholder="Enter your name"
            className={`mt-1 block w-full px-3 py-2 border font-semibold bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 ${
              errors.name
                ? "border-red-500"
                : "border-gray-300 dark:border-gray-700"
            } rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-[#44C3BB] focus:border-transparent`}
            {...register("name", { required: "Name is required" })}
            disabled={isLoading || mutationLoading}
          />
          {errors.name && (
            <p className="mt-1 text-sm text-red-500">{errors.name.message}</p>
          )}
        </div>

        {/* Email Field */}
        <div className="mb-7">
          <label
            htmlFor="email"
            className="block text-sm font-bold text-gray-700 dark:text-gray-300"
          >
            Email
          </label>
          <input
            id="email"
            type="email"
            placeholder="Enter your Email"
            className={`mt-1 block w-full px-3 py-2 border font-semibold bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 ${
              errors.email
                ? "border-red-500"
                : "border-gray-300 dark:border-gray-700"
            } rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-[#44C3BB] focus:border-transparent`}
            {...register("email", {
              required: "Email is required",
              pattern: {
                value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                message: "Invalid email address",
              },
            })}
            disabled={isLoading || mutationLoading}
          />
          {errors.email && (
            <p className="mt-1 text-sm text-red-500">{errors.email.message}</p>
          )}
        </div>

        {/* Subject Field */}
        <div className="mb-7">
          <label
            htmlFor="subject"
            className="block text-sm font-bold text-gray-700 dark:text-gray-300"
          >
            Subject
          </label>
          <input
            id="subject"
            type="text"
            placeholder="Enter the subject of your message"
            className={`mt-1 block w-full px-3 py-2 border font-semibold bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 ${
              errors.subject
                ? "border-red-500"
                : "border-gray-300 dark:border-gray-700"
            } rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-[#44C3BB] focus:border-transparent`}
            {...register("subject", { required: "Subject is required" })}
            disabled={isLoading || mutationLoading}
          />
          {errors.subject && (
            <p className="mt-1 text-sm text-red-500">
              {errors.subject.message}
            </p>
          )}
        </div>

        {/* Message Field */}
        <div className="mb-7">
          <label
            htmlFor="message"
            className="block text-sm font-bold text-gray-700 dark:text-gray-300"
          >
            Message
          </label>
          <textarea
            id="message"
            rows={6}
            placeholder="Enter your comment"
            className={`mt-1 block w-full px-3 py-2 border font-semibold bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 ${
              errors.message
                ? "border-red-500"
                : "border-gray-300 dark:border-gray-700"
            } rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-[#44C3BB] focus:border-transparent`}
            {...register("message", { required: "Message is required" })}
            disabled={isLoading || mutationLoading}
          ></textarea>
          {errors.message && (
            <p className="mt-1 text-sm text-red-500">
              {errors.message.message}
            </p>
          )}
        </div>

        {/* Submit Button */}
        <button
          type="submit"
          className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] hover:from-[#3AB5AD] hover:to-[#44C3BB] focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-[#44C3BB] disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={isLoading || mutationLoading}
        >
          {isLoading || mutationLoading ? "Submitting..." : "Submit"}
        </button>

        {/* Display general submission status */}
        {isSubmitted && !error && (
          <p className="mt-2 text-sm text-green-600 dark:text-green-400">
            Form submitted successfully! The admin has been notified.
          </p>
        )}
        {error && (
          <p className="mt-2 text-sm text-red-600 dark:text-red-400">
            Error: {error}
          </p>
        )}
      </form>
    </div>
  );
};

export default ContactForm;
