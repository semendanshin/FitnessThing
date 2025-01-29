"use client";

import { use } from "react";

import { WorkoutResults } from "@/components/workout-results";
import { PageHeader } from "@/components/page-header";

export default function RoutineDetailsPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);

  return (
    <div className="py-4 flex flex-col h-full gap-4">
      <PageHeader enableBackButton={false} title="Так держать!" />
      <WorkoutResults className="px-4" id={id} />
    </div>
  );
}
