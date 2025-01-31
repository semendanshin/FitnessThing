"use client";

import { Button } from "@nextui-org/button";
import { Card, CardBody } from "@nextui-org/card";
import { Link } from "@nextui-org/link";
import { useDisclosure } from "@nextui-org/modal";
import { DropdownItem, ScrollShadow } from "@nextui-org/react";
import { useRouter } from "next/navigation";
import { use, useEffect, useState } from "react";
import { toast } from "react-toastify";

import { ChevronRightIcon, PlusIcon } from "@/config/icons";
import { ModalSelectExercise } from "@/components/pick-exercises-modal";
import { PageHeader } from "@/components/page-header";
import { Loading } from "@/components/loading";
import {
  WorkoutExerciseLogDetails,
  WorkoutGetWorkoutResponse,
  WorkoutRoutineDetailResponse,
  WorkoutExerciseInstanceDetails,
} from "@/api/api.generated";
import { authApi, errUnauthorized } from "@/api/api";

function ExerciseLogCard({
  exerciseLogDetails,
  workoutId,
  exerciseInstanceDetails,
}: {
  exerciseLogDetails: WorkoutExerciseLogDetails;
  workoutId: string;
  exerciseInstanceDetails?: WorkoutExerciseInstanceDetails;
}) {
  return (
    <Card
      fullWidth
      as={Link}
      className="flex flex-row items-center justify-between p-4 cursor-pointer"
      href={`/workouts/${workoutId}/exerciseLogs/${exerciseLogDetails.exerciseLog?.id}`}
      shadow="sm"
    >
      <div className="flex flex-col items-start justify-between w-full gap-3">
        <div className="flex flex-col">
          <p className="text-lg font-bold">
            {exerciseLogDetails.exercise?.name}
          </p>
          {exerciseInstanceDetails && (
            <div className="text-xs font-light text-default-600">
              {exerciseInstanceDetails?.sets?.length} подходов x{" "}
              {(exerciseInstanceDetails?.sets?.reduce(
                (acc, set) => acc + set.reps!,
                0,
              )! /
                exerciseInstanceDetails.sets?.length!) |
                0}{" "}
              раз
            </div>
          )}
        </div>
        {exerciseLogDetails.setLogs!.length > 0 && (
          <CardBody className="flex flex-col w-full gap-1 p-0">
            {exerciseLogDetails.setLogs?.map((setLog, setNum) => (
              <div key={setLog.id} className="flex flex-row w-full gap-2">
                <div className="text-sm font-semibold min-w-fit w-3">
                  {setNum + 1}.
                </div>
                <div className="text-sm font-semibold w-fit">
                  {setLog?.weight} кг
                </div>
                <div className="text-sm font-semibold w-fit">x</div>
                <div className="text-sm font-semibold w-fit">
                  {setLog?.reps}
                </div>
              </div>
            ))}
          </CardBody>
        )}
      </div>
      <div className="flex flex-col items-center justify-between">
        <ChevronRightIcon className="w-4 h-4" fill="currentColor" />
      </div>
    </Card>
  );
}

export default function RoutineDetailsPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const [isLoading, setIsLoading] = useState(true);
  const [isError, setIsError] = useState(false);

  const [workoutDetails, setWorkoutDetails] =
    useState<WorkoutGetWorkoutResponse>({});

  const [routineDetails, setRoutineDetails] =
    useState<WorkoutRoutineDetailResponse>();

  const { id } = use(params);

  const router = useRouter();

  const { isOpen, onOpen, onClose } = useDisclosure();

  async function fetchWorkoutDetails() {
    await authApi.v1
      .workoutServiceGetWorkout(id)
      .then((response) => {
        console.log(response.data);
        setWorkoutDetails(response.data!);
      })
      .catch((error) => {
        console.log(error);
        if (error === errUnauthorized || error.response?.status === 401) {
          router.push("/auth/login");

          return;
        }
        throw error;
      });
  }

  async function fetchRoutineDetails(id: string) {
    await authApi.v1
      .routineServiceGetRoutineDetail(id)
      .then((response) => {
        console.log(response.data);
        setRoutineDetails(response.data);
      })
      .catch((error) => {
        console.log(error);
        if (error === errUnauthorized || error.response?.status === 401) {
          router.push("/auth/login");

          return;
        }
        throw error;
      });
  }

  async function fetchData() {
    setIsLoading(true);
    try {
      await fetchWorkoutDetails();
    } catch (error) {
      console.log(error);
      toast.error("Failed to fetch workout details");
      setIsError(true);
    } finally {
      setIsLoading(false);
    }
  }

  async function finishWorkout() {
    try {
      await authApi.v1
        .workoutServiceCompleteWorkout(id, {})
        .then((response) => {
          console.log(response.data);
          router.push(`/workouts/${id}/results`);
        })
        .catch((error) => {
          console.log(error);
          if (error === errUnauthorized || error.response?.status === 401) {
            router.push("/auth/login");

            return;
          }
          throw error;
        });
    } catch (error) {
      console.log(error);
      toast.error("Не удалось завершить тренировку");
    } finally {
      setIsLoading(false);
    }
  }

  async function addExercisesToWorkout(exerciseIds: string[]) {
    try {
      for (const exerciseId of exerciseIds) {
        await authApi.v1
          .workoutServiceLogExercise(id, {
            exerciseId,
          })
          .then((response) => {
            console.log(response.data);
          })
          .catch((error) => {
            console.log(error);
            if (error === errUnauthorized || error.response?.status === 401) {
              router.push("/auth/login");

              return;
            }
            throw error;
          });
      }
      onClose();
      await fetchWorkoutDetails();
    } catch (error) {
      console.log(error);
      toast.error("Failed to add exercises to workout");
    } finally {
      setIsLoading(false);
    }
  }

  async function onDelete() {
    try {
      await authApi.v1
        .workoutServiceDeleteWorkout(id)
        .then((response) => {
          console.log(response.data);
          router.push("/");
        })
        .catch((error) => {
          console.log(error);
          if (error === errUnauthorized || error.response?.status === 401) {
            router.push("/auth/login");

            return;
          }
          throw error;
        });
    } catch (error) {
      console.log(error);
      toast.error("Не удалось удалить тренировку");
    }
  }

  useEffect(() => {
    fetchData();
  }, []);

  useEffect(() => {
    if (workoutDetails.workout?.routineId) {
      fetchRoutineDetails(workoutDetails.workout?.routineId);
    }
  }, [workoutDetails]);

  if (isLoading) {
    return <Loading />;
  }

  if (isError) {
    return (
      <div className="p-4">
        <h2 className="text-lg text-red-500">Ошибка при загрузке данных</h2>
        <p>Проверьте соединение с сервером или обновите страницу.</p>
      </div>
    );
  }

  return (
    <>
      <div className="py-4 flex flex-col h-full flex-grow max-w-full basis-full">
        <PageHeader enableBackButton={true} title={"Тренировка"}>
          <DropdownItem
            key="delete"
            className="text-danger"
            color="danger"
            onPress={onDelete}
          >
            Удалить
          </DropdownItem>
        </PageHeader>
        <section className="flex flex-col gap-4 h-full overflow-y-scroll">
          <ScrollShadow size={50}>
            <div className="flex flex-col gap-2 p-4 ">
              {workoutDetails.exerciseLogs?.map((exerciseLogDetails, index) => (
                <ExerciseLogCard
                  key={index}
                  exerciseInstanceDetails={routineDetails?.exerciseInstances?.find(
                    (exerciseInstance) =>
                      exerciseInstance.exerciseInstance?.exerciseId ===
                      exerciseLogDetails.exerciseLog?.exerciseId,
                  )}
                  exerciseLogDetails={exerciseLogDetails}
                  workoutId={id}
                />
              ))}
              <Card className="p-2">
                <Button
                  className="w-full"
                  onPress={() => {
                    onOpen();
                  }}
                >
                  <PlusIcon className="w-4 h-4" />
                  <span>Добавить упражнение</span>
                </Button>
              </Card>
            </div>
          </ScrollShadow>
        </section>
        <section className="w-full px-4 py-2 ">
          <Button
            className="w-full"
            color="primary"
            onPress={async () => {
              await finishWorkout();
            }}
          >
            Завершить тренировку
          </Button>
        </section>
      </div>
      <ModalSelectExercise
        excludeExerciseIds={workoutDetails.exerciseLogs!.map(
          (exerciseLog) => exerciseLog.exerciseLog!.exerciseId!,
        )}
        isOpen={isOpen}
        onClose={onClose}
        onSubmit={addExercisesToWorkout}
      />
    </>
  );
}
