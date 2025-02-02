"use client";

import { Button } from "@nextui-org/button";
import { Card, CardBody, CardHeader } from "@nextui-org/card";
import { DropdownItem } from "@nextui-org/dropdown";
import { Form } from "@nextui-org/form";
import { Input, Textarea } from "@nextui-org/input";
import {
  Modal,
  ModalContent,
  ModalHeader,
  useDisclosure,
} from "@nextui-org/modal";
import { useRouter } from "next/navigation";
import { use, useEffect, useState } from "react";
import { toast } from "react-toastify";
import Link from "next/link";

import { ChevronRightIcon, PlusIcon } from "@/config/icons";
import { ModalSelectExercise } from "@/components/pick-exercises-modal";
import { PageHeader } from "@/components/page-header";
import { Loading } from "@/components/loading";
import {
  WorkoutExerciseInstanceDetails,
  WorkoutRoutineDetailResponse,
} from "@/api/api.generated";
import { authApi, errUnauthorized } from "@/api/api";

export default function RoutineDetailsPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const [isLoading, setIsLoading] = useState(true);

  const [routineDetails, setRoutineDetails] =
    useState<WorkoutRoutineDetailResponse>({});

  const { id } = use(params);

  const { isOpen, onOpen, onClose } = useDisclosure();

  const {
    isOpen: isRenameOpen,
    onOpen: onRenameOpen,
    onOpenChange: onRenameOpenChange,
  } = useDisclosure();

  const router = useRouter();

  async function fetchRoutineDetails() {
    console.log(id);
    await authApi.v1
      .routineServiceGetRoutineDetail(id)
      .then((response) => {
        console.log(response.data);
        setRoutineDetails(response.data);
      })
      .catch((error) => {
        console.log(error);
        if (
          error === errUnauthorized ||
          (error as any).response?.status === 401
        ) {
          router.push("/auth/login");

          return;
        }
        toast.error("Ошибка при загрузке данных");
      });
  }

  async function addExerciseToRoutine(exerciseId: string) {
    await authApi.v1
      .routineServiceAddExerciseToRoutine(id, { exerciseId: exerciseId })
      .then((response) => {
        console.log(response.data);
      })
      .catch((error) => {
        console.log(error);
        if (error === errUnauthorized || error.response?.status === 401) {
          router.push("/auth/login");

          return;
        }
        toast.error("Ошибка при добавлении упражнения");
      });
  }

  async function submitPickExercise(exerciseIds: string[]) {
    console.log(exerciseIds);
    try {
      for (const exerciseId of exerciseIds) {
        await addExerciseToRoutine(exerciseId);
      }
      await fetchRoutineDetails();
    } catch (error) {
      console.log(error);
      if (
        error === errUnauthorized ||
        (error as any).response?.status === 401
      ) {
        router.push("/auth/login");

        return;
      }
      toast.error("Ошибка при добавлении упражнения");
    }
  }

  async function fetchData() {
    setIsLoading(true);
    await fetchRoutineDetails();
    setIsLoading(false);
  }

  useEffect(() => {
    fetchData();
  }, []);

  if (isLoading) {
    return <Loading />;
  }

  function ExercieInstanceCard({
    exerciseInstanceDetails,
  }: {
    exerciseInstanceDetails: WorkoutExerciseInstanceDetails;
  }) {
    const setsCount = exerciseInstanceDetails.sets?.length || 0;

    return (
      <Card
        key={exerciseInstanceDetails.exerciseInstance!.id}
        fullWidth
        as={Link}
        className="flex flex-row flex-grow p-2 gap-4 justify-between"
        href={`/routines/${id}/exerciseInstance/${exerciseInstanceDetails.exerciseInstance!.id}`}
        shadow="sm"
      >
        <div className="flex flex-col items-start justify-between p-2">
          <CardHeader className="p-0">
            <p className="text-m font-bold">
              {exerciseInstanceDetails.exercise!.name}
            </p>
          </CardHeader>
          <CardBody className="p-0">
            <div className="flex flex-row gap-1">
              <p className="text-xs text-gray-400 whitespace-nowrap">
                {setsCount}{" "}
                {"подход" +
                  (setsCount % 10 === 1
                    ? ""
                    : setsCount % 10 >= 2 && setsCount % 10 <= 4
                      ? "а"
                      : "ов")}
                {" •"}
              </p>
              <p className="text-xs text-gray-400 whitespace-nowrap">
                {exerciseInstanceDetails.exercise!.targetMuscleGroups!.join(
                  ", ",
                )}
              </p>
            </div>
          </CardBody>
        </div>
        <div className="flex flex-col items-center justify-center p-2">
          <ChevronRightIcon className="w-4 h-4" fill="currentColor" />
        </div>
      </Card>
    );
  }

  function ModalRenameRoutine() {
    const [routineName, setRoutineName] = useState(
      routineDetails.routine?.name,
    );
    const [routineDescription, setRoutineDescription] = useState(
      routineDetails.routine?.description,
    );
    const [isButtonLoading, setIsButtonLoading] = useState(false);

    async function updateRoutine() {
      setIsButtonLoading(true);
      try {
        await authApi.v1.routineServiceUpdateRoutine(id, {
          name: routineName,
          description: routineDescription,
        });
        await fetchRoutineDetails();
      } catch (error) {
        console.log(error);
        toast.error("Ошибка при обновлении рутины");
        if (
          error === errUnauthorized ||
          (error as any).response?.status === 401
        ) {
          router.push("/auth/login");

          return;
        }
      } finally {
        setIsButtonLoading(false);
      }
    }

    return (
      <Modal
        className="overflow-y-auto p-2 w-full"
        isOpen={isRenameOpen}
        placement="center"
        scrollBehavior="inside"
        size="xs"
        onClose={onRenameOpenChange}
      >
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader className="p-2">Переименовать рутину</ModalHeader>

              <Form
                className="inline-block text-center justify-center w-full max-w-lg p-0"
                onSubmit={async (e) => {
                  e.preventDefault();
                  await updateRoutine();
                  onClose();
                }}
              >
                <div className="grid grid-cols-1 gap-4 p-2">
                  <Input
                    autoFocus
                    placeholder="Название"
                    value={routineName}
                    onChange={(e) => setRoutineName(e.target.value)}
                  />
                  <Textarea
                    placeholder="Описание"
                    value={routineDescription ? routineDescription : ""}
                    onChange={(e) => setRoutineDescription(e.target.value)}
                  />
                  <Button
                    color="primary"
                    isLoading={isButtonLoading}
                    type="submit"
                  >
                    Сохранить
                  </Button>
                </div>
              </Form>
            </>
          )}
        </ModalContent>
      </Modal>
    );
  }

  async function onDelete() {
    await authApi.v1
      .routineServiceDeleteRoutine(id)
      .catch((error) => {
        console.log(error);
        if (
          error === errUnauthorized ||
          (error as any).response?.status === 401
        ) {
          router.push("/auth/login");

          return;
        }
        toast.error("Ошибка при удалении рутины");
      })
      .then(() => {
        router.back();
      });
  }

  return (
    <>
      <div className="flex flex-col py-4 gap-4">
        <PageHeader enableBackButton title={routineDetails.routine?.name!}>
          <DropdownItem key="rename" color="primary" onPress={onRenameOpen}>
            Переименовать
          </DropdownItem>
          <DropdownItem
            key="delete"
            className="text-danger"
            color="danger"
            onPress={onDelete}
          >
            Удалить
          </DropdownItem>
        </PageHeader>
        {routineDetails.routine?.description ? (
          <div className="px-4">
            <p className="text-sm text-gray-500">
              {routineDetails.routine?.description}
            </p>
          </div>
        ) : null}
        <div className="grid grid-cols-1 gap-4 px-4">
          {routineDetails.exerciseInstances?.map(
            (exerciseInstanceDetails: WorkoutExerciseInstanceDetails) => (
              <ExercieInstanceCard
                key={exerciseInstanceDetails.exerciseInstance!.id}
                exerciseInstanceDetails={exerciseInstanceDetails}
              />
            ),
          )}
          <Card className="p-2">
            <Button onPress={onOpen}>
              <PlusIcon className="w-4 h-4" />
              Добавить упражнение
            </Button>
          </Card>
        </div>
      </div>
      <ModalSelectExercise
        excludeExerciseIds={routineDetails.exerciseInstances!.map(
          (ei) => ei.exercise!.id!,
        )}
        isOpen={isOpen}
        onClose={onClose}
        onSubmit={submitPickExercise}
      />

      <ModalRenameRoutine />
    </>
  );
}
