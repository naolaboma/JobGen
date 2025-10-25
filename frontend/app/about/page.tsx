export default function AboutPage() {
  return (
    <main className="max-w-3xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-4">About JobGen</h1>
      <p className="text-gray-700 mb-6">
        JobGen helps candidates generate tailored resumes and cover letters,
        discover relevant jobs, and streamline applications with AI assistance.
      </p>

      <section className="space-y-3 mb-8">
        <h2 className="text-xl font-semibold">Our Mission</h2>
        <p className="text-gray-700">
          Make job search faster, smarter, and more inclusive by giving everyone
          access to AI-powered career tools.
        </p>
      </section>

      <section className="space-y-3 mb-8">
        <h2 className="text-xl font-semibold">What We Do</h2>
        <ul className="list-disc pl-5 text-gray-700 space-y-2">
          <li>Parse and enhance your CV to highlight your strengths.</li>
          <li>Match you with roles based on your skills and interests.</li>
          <li>Generate role-specific cover letters and application answers.</li>
          <li>Track applications and get timely reminders.</li>
        </ul>
      </section>

      <section className="space-y-3">
        <h2 className="text-xl font-semibold">Get in Touch</h2>
        <p className="text-gray-700">
          Have feedback or partnership ideas? Visit our{" "}
          <a href="/contact" className="text-teal-600 hover:underline">
            Contact
          </a>{" "}
          page—we’d love to hear from you.
        </p>
      </section>
    </main>
  );
}
