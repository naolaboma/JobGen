import { ArrowRightIcon } from "lucide-react";
import React, { JSX } from "react";
import { Button } from "./ui/button";

export const AppsSection = (): JSX.Element => {
  return (
    <section className="flex w-full items-center justify-center gap-8 md:gap-12 lg:gap-[100px] px-4 md:px-8 lg:px-[220px] py-16 md:py-24 lg:py-[140px] relative bg-[#7bbfb3]">
      <img
        className="absolute w-full h-[400px] md:h-[550px] lg:h-[700px] top-2.5 left-0"
        alt="Element"
        src="/element-8.png"
      />

      <div className="inline-flex flex-col lg:flex-row items-center justify-center gap-8 md:gap-12 lg:gap-[100px] relative flex-[0_0_auto] w-full max-w-[1480px]">
        <img
          className="relative w-full max-w-[400px] md:max-w-[500px] lg:max-w-[582px] h-auto lg:h-[470.8px] object-contain"
          alt="Apps"
          src="/apps.png"
        />

        <div className="flex flex-col w-full lg:w-[798px] items-start gap-8 md:gap-12 lg:gap-[60px] relative">
          <div className="flex flex-col items-start gap-4 md:gap-6 relative self-stretch w-full flex-[0_0_auto]">
            <h1 className="relative self-stretch mt-[-1.00px] font-header-h1 font-[number:var(--header-h1-font-weight)] text-white text-2xl md:text-4xl lg:text-[length:var(--header-h1-font-size)] text-center lg:text-left tracking-[var(--header-h1-letter-spacing)] leading-tight md:leading-[var(--header-h1-line-height)] [font-style:var(--header-h1-font-style)]">
              Focus on Your Career, Not on Switching Tools
            </h1>

            <p className="relative self-stretch [font-family:'Inter',Helvetica] font-normal text-white text-base md:text-lg leading-6 md:leading-[18px] text-center lg:text-left">
              <span className="font-bold tracking-[-0.06px] leading-6 md:leading-[30px]">
                JobGen is designed to fit seamlessly into your existing
                workflow. We integrate with the platforms you already use to
                make your job search smoother and more effective.
              </span>
            </p>
          </div>

          <Button className="px-6 py-4 md:px-10 md:py-5 inline-flex items-center justify-center gap-2.5 relative flex-[0_0_auto] bg-[#38b2ac] rounded-lg h-auto hover:bg-[#319795] mx-auto lg:mx-0">
            <span className="relative w-fit mt-[-1.00px] [font-family:'Inter',Helvetica] font-normal text-white text-base md:text-lg leading-[18px]">
              <span className="font-bold tracking-[-0.06px] leading-[23px]">
                Read more
              </span>
            </span>
            <ArrowRightIcon className="relative flex-[0_0_auto] mr-[-0.50px] w-4 h-4" />
          </Button>
        </div>
      </div>
    </section>
  );
};
