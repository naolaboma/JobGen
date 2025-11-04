export default function PrivacyPage() {
  return (
    <main className="max-w-3xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-4">Privacy Policy</h1>
      <p className="text-gray-700 mb-6">
        Effective date: {new Date().getFullYear()}
      </p>

      <section className="space-y-2 mb-6">
        <h2 className="text-xl font-semibold">Information We Collect</h2>
        <p className="text-gray-700">
          Account details, profile/CV data you upload, usage analytics, and
          communications.
        </p>
      </section>

      <section className="space-y-2 mb-6">
        <h2 className="text-xl font-semibold">How We Use Information</h2>
        <p className="text-gray-700">
          To provide and improve JobGen services, personalize recommendations,
          support you, and ensure security.
        </p>
      </section>

      <section className="space-y-2 mb-6">
        <h2 className="text-xl font-semibold">Sharing</h2>
        <p className="text-gray-700">
          We do not sell your data. We may share with service providers under
          strict contracts and as required by law.
        </p>
      </section>

      <section className="space-y-2 mb-6">
        <h2 className="text-xl font-semibold">Security</h2>
        <p className="text-gray-700">
          We use industry-standard measures to protect your data. No method is
          100% secure.
        </p>
      </section>

      <section className="space-y-2 mb-6">
        <h2 className="text-xl font-semibold">Your Rights</h2>
        <p className="text-gray-700">
          You can access, correct, or delete your data, and control marketing
          communications, subject to applicable law.
        </p>
      </section>

      <section className="space-y-2">
        <h2 className="text-xl font-semibold">Contact</h2>
        <p className="text-gray-700">
          Questions? Visit our{" "}
          <a href="/contact" className="text-teal-600 hover:underline">
            Contact
          </a>{" "}
          page.
        </p>
      </section>
    </main>
  );
}
