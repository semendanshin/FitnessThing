"use client";

import { use } from "react";

import { WorkoutResults } from "@/components/workout-results";

export default function RoutineDetailsPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);

  return <WorkoutResults enableBackButton id={id} />;
}
