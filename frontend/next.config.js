/** @type {import('next').NextConfig} */
const nextConfig = {
  devIndicators: {
    appIsrStatus: false,
  },
  reactStrictMode: false,
  output: "standalone",
};

module.exports = nextConfig;
