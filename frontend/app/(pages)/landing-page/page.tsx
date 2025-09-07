"use client";

import React, { JSX } from "react";
import { LandingHeader } from "../../components/landingPage/LandingHeader";
import { HeroSection } from "../../components/landingPage/HeroSection";
import { FeaturesSection } from "../../components/landingPage/FeaturesSection";
import { PricingSection } from "../../components/landingPage/PricingSection";
import { TestimonialsSection } from "../../components/landingPage/TestimonialsSection";
import { CtaSection } from "../../components/landingPage/CtaSection";
import { FooterSection } from "../../components/landingPage/FooterSection";

const LandingPage = (): JSX.Element => {
  const [mobileMenuOpen, setMobileMenuOpen] = React.useState(false);
  return (
    <div className="min-h-screen bg-white">
      <LandingHeader
        mobileMenuOpen={mobileMenuOpen}
        setMobileMenuOpen={setMobileMenuOpen}
      />
      <HeroSection />
      <FeaturesSection />
      <PricingSection />
      <TestimonialsSection />
      <CtaSection />
      <FooterSection />
    </div>
  );
};

export default LandingPage;
