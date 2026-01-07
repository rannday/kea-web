document.addEventListener("DOMContentLoaded", function () {
  if ("speculationrules" in document) {
      console.log("[Prefetch] Speculation Rules API supported.");
      return; // The browser supports Speculation Rules API; no need for fallback.
  }

  console.log("[Prefetch] Speculation Rules API NOT supported. Using link prefetch fallback.");
  
  const PREFETCH_URLS = [
      "/search",
      "/about",
      "/stats",
      "/debug"
  ];

  PREFETCH_URLS.forEach(url => {
      if (!document.querySelector(`link[rel="prefetch"][href="${url}"]`)) {
          const link = document.createElement("link");
          link.rel = "prefetch";
          link.href = url;
          document.head.appendChild(link);
          console.log(`[Prefetch] Added fallback prefetch: ${url}`);
      }
  });
});
