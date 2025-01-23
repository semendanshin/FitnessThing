export default function LoginLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <section className="flex flex-col items-center justify-center gap-4 w-full">
      <div className="flex flex-col items-center justify-center gap-4 w-full max-w-lg p-4">
        {children}
      </div>
    </section>
  );
}
