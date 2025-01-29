"use client";

import { use } from "react";
import { useRouter } from "next/navigation";
import { toast } from "react-toastify";
import { DropdownItem } from "@nextui-org/dropdown";

import { WorkoutResults } from "@/components/workout-results";
import { PageHeader } from "@/components/page-header";
import { authApi } from "@/api/api";

export default function RoutineDetailsPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);

  const router = useRouter();

  async function handleDelete() {
    await authApi.v1
      .workoutServiceDeleteWorkout(id)
      .then(() => {
        router.back();
      })
      .catch((error) => {
        console.log(error);
        toast.error("Failed to delete workout");
      });
  }

  return (
    <div className="py-4 flex flex-col gap-4">
      <PageHeader enableBackButton={true} title="Так держать!">
        <DropdownItem
          key="delete"
          className="text-danger"
          color="danger"
          onPress={handleDelete}
        >
          Удалить
        </DropdownItem>
      </PageHeader>
      <WorkoutResults className="px-4" id={id} />
    </div>
  );
}
