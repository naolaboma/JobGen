export default function TermsPage() {
  return (
    <main className="max-w-3xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-4">Terms of Service</h1>
      <p className="text-gray-700 mb-6">
        Welcome to JobGen. By accessing or using our services, you agree to
        these Terms of Service.
      </p>

      <section className="space-y-2 mb-6">
        <h2 className="text-xl font-semibold">1. Use of the Service</h2>
        <p className="text-gray-700">
          You agree to use JobGen only for lawful purposes and in accordance
          with these terms.
        </p>
      </section>

      <section className="space-y-2 mb-6">
        <h2 className="text-xl font-semibold">2. Accounts</h2>
        <p className="text-gray-700">
          You are responsible for maintaining the confidentiality of your
          account and credentials.
        </p>
      </section>

      <section className="space-y-2 mb-6">
        <h2 className="text-xl font-semibold">3. Content</h2>
        <p className="text-gray-700">
          You retain rights to your content. By submitting content, you grant us
          a limited license to operate and improve the service.
        </p>
      </section>

      <section className="space-y-2 mb-6">
        <h2 className="text-xl font-semibold">4. Intellectual Property</h2>
        <p className="text-gray-700">
          All JobGen trademarks, logos, and service marks are our property or
          our licensors’.
        </p>
      </section>

      <section className="space-y-2 mb-6">
        <h2 className="text-xl font-semibold">
          5. Disclaimers and Limitation of Liability
        </h2>
        <p className="text-gray-700">
          The service is provided “as is.” JobGen will not be liable for
          indirect or consequential damages to the extent permitted by law.
        </p>
      </section>

      <section className="space-y-2 mb-6">
        <h2 className="text-xl font-semibold">6. Changes</h2>
        <p className="text-gray-700">
          We may update these terms from time to time. We will post the updated
          terms with a new effective date.
        </p>
      </section>

      <section className="space-y-2">
        <h2 className="text-xl font-semibold">7. Contact</h2>
        <p className="text-gray-700">
          Questions about these terms? Visit our{" "}
          <a href="/contact" className="text-teal-600 hover:underline">
            Contact
          </a>{" "}
          page.
        </p>
      </section>
    </main>
  );
}
